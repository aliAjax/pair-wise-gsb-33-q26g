package service

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

func newMergeService(t *testing.T) (*PlantMergeService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock := newServiceDB(t)
	return NewPlantMergeService(db, repository.NewPlantMergeRepository(db),
		repository.NewPlantSpeciesRepository(db), newTestLogger()), mock
}

var mergeTaskColumns = []string{
	"id", "keep_id", "source_id", "keep_name", "source_name", "status", "stage",
	"stats", "last_error", "operator_id", "created_at", "updated_at", "finished_at",
}

// mergeTaskRow mocks the FindByID query at the start/end of execute.
func mergeTaskRow(mock sqlmock.Sqlmock, id, keepID, sourceID uint, status, stage, stats string) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `plant_merge_tasks` WHERE `plant_merge_tasks`.`id` = ?")).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(mergeTaskColumns).
			AddRow(id, keepID, sourceID, "保留卡", "待合并卡", status, stage, stats, "", uint(1),
				time.Now(), time.Now(), nil))
}

// stageStmt describes one migration statement and its affected row count.
type stageStmt struct {
	sql  string
	args []driver.Value
	rows int64
}

// expectStage mocks one stage transaction: the migration statements plus the
// checkpoint update carrying the accumulated stats JSON.
func expectStage(mock sqlmock.Sqlmock, taskID uint, stage, wantStats string, stmts ...stageStmt) {
	mock.ExpectBegin()
	for _, s := range stmts {
		e := mock.ExpectExec(regexp.QuoteMeta(s.sql))
		if len(s.args) > 0 {
			e.WithArgs(s.args...)
		}
		e.WillReturnResult(sqlmock.NewResult(0, s.rows))
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `plant_merge_tasks` SET")).
		WithArgs(stage, wantStats, sqlmock.AnyArg(), taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

// expectTaskUpdate mocks a standalone task UPDATE (GORM wraps it in an
// implicit transaction).
func expectTaskUpdate(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `plant_merge_tasks` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

func TestValidateSourceIDs(t *testing.T) {
	cases := []struct {
		name    string
		keep    uint
		sources []uint
		wantErr bool
	}{
		{"ok", 1, []uint{2, 3}, false},
		{"empty", 1, nil, true},
		{"keep equals source", 1, []uint{1}, true},
		{"duplicate source", 1, []uint{2, 2}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateSourceIDs(c.keep, c.sources)
			if c.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.wantErr {
				var appErr *util.AppError
				if !errors.As(err, &appErr) || appErr.HTTPStatus != 422 {
					t.Fatalf("expected 422 AppError, got %v", err)
				}
			}
		})
	}
}

func TestParseMergeStats(t *testing.T) {
	if got := parseMergeStats(""); got != (dto.MergeStats{}) {
		t.Fatalf("empty stats should be zero, got %+v", got)
	}
	got := parseMergeStats(`{"favorites_moved":2,"pests_moved":3}`)
	if got.FavoritesMoved != 2 || got.PestsMoved != 3 {
		t.Fatalf("unexpected stats: %+v", got)
	}
	if round := parseMergeStats(marshalMergeStats(&got)); round != got {
		t.Fatalf("round trip mismatch: %+v != %+v", round, got)
	}
}

// TestMergeExecuteRunsAllStages verifies the five stages run in order, each in
// its own transaction with a stats checkpoint, then the task is marked done.
func TestMergeExecuteRunsAllStages(t *testing.T) {
	svc, mock := newMergeService(t)
	mergeTaskRow(mock, 1, 10, 20, model.MergeRunning, "", "{}")

	s1 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":0,"gardens_deduped":0,"reminders_moved":0,"reminders_merged":0,"pests_moved":0}`
	expectStage(mock, 1, model.MergeStageFavorites, s1,
		stageStmt{`DELETE f FROM favorites f
		JOIN favorites k ON k.user_id = f.user_id AND k.target_type = 'plant' AND k.target_id = ?
		WHERE f.target_type = 'plant' AND f.target_id = ?`, []driver.Value{uint(10), uint(20)}, 1},
		stageStmt{`UPDATE favorites SET target_id = ? WHERE target_type = 'plant' AND target_id = ?`, []driver.Value{uint(10), uint(20)}, 2},
	)
	s2 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":1,"reminders_moved":0,"reminders_merged":0,"pests_moved":0}`
	expectStage(mock, 1, model.MergeStageGardens, s2,
		stageStmt{`UPDATE user_gardens k
		JOIN user_gardens g ON g.user_id = k.user_id AND g.plant_species_id = ?
		SET k.care_reminder_id = g.care_reminder_id
		WHERE k.plant_species_id = ? AND k.care_reminder_id = 0 AND g.care_reminder_id <> 0`, []driver.Value{uint(20), uint(10)}, 0},
		stageStmt{`DELETE g FROM user_gardens g
		JOIN user_gardens k ON k.user_id = g.user_id AND k.plant_species_id = ?
		WHERE g.plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 1},
		stageStmt{`UPDATE user_gardens SET plant_species_id = ? WHERE plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 1},
	)
	s3 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":1,"reminders_moved":3,"reminders_merged":1,"pests_moved":0}`
	expectStage(mock, 1, model.MergeStageReminders, s3,
		stageStmt{`UPDATE care_reminders k
		JOIN care_reminders r ON r.user_id = k.user_id AND r.task_title = k.task_title
		SET k.remind_date = GREATEST(k.remind_date, r.remind_date),
			k.status = IF(k.status = 'done' AND r.status = 'done', 'done', 'pending')
		WHERE k.plant_species_id = ? AND r.plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 1},
		stageStmt{`UPDATE user_gardens g
		JOIN care_reminders r ON g.care_reminder_id = r.id
		JOIN care_reminders k ON k.user_id = r.user_id AND k.task_title = r.task_title
		SET g.care_reminder_id = k.id
		WHERE r.plant_species_id = ? AND k.plant_species_id = ?`, []driver.Value{uint(20), uint(10)}, 0},
		stageStmt{`DELETE r FROM care_reminders r
		JOIN care_reminders k ON k.user_id = r.user_id AND k.task_title = r.task_title
		WHERE r.plant_species_id = ? AND k.plant_species_id = ?`, []driver.Value{uint(20), uint(10)}, 1},
		stageStmt{`UPDATE care_reminders SET plant_species_id = ? WHERE plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 3},
	)
	s4 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":1,"reminders_moved":3,"reminders_merged":1,"pests_moved":2}`
	expectStage(mock, 1, model.MergeStagePests, s4,
		stageStmt{`UPDATE disease_pests SET plant_species_id = ? WHERE plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 2},
	)
	expectStage(mock, 1, model.MergeStageRetire, s4,
		stageStmt{`DELETE FROM plant_species WHERE id = ?`, []driver.Value{uint(20)}, 1},
	)
	expectTaskUpdate(mock) // MarkDone
	mergeTaskRow(mock, 1, 10, 20, model.MergeDone, model.MergeStageRetire, s4)

	task, err := svc.execute(1)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if task.Status != model.MergeDone {
		t.Fatalf("expected done, got %s", task.Status)
	}
	stats := parseMergeStats(task.Stats)
	want := dto.MergeStats{FavoritesMoved: 2, FavoritesDeduped: 1, GardensMoved: 1, GardensDeduped: 1, RemindersMoved: 3, RemindersMerged: 1, PestsMoved: 2}
	if stats != want {
		t.Fatalf("stats mismatch: got %+v want %+v", stats, want)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestMergeExecuteResumesFromCheckpoint verifies a failed task reruns only the
// stages after its checkpoint instead of repeating finished ones.
func TestMergeExecuteResumesFromCheckpoint(t *testing.T) {
	svc, mock := newMergeService(t)
	stats := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":0,"reminders_moved":0,"reminders_merged":0,"pests_moved":0}`
	mergeTaskRow(mock, 1, 10, 20, model.MergeFailed, model.MergeStageGardens, stats)

	// favorites and gardens stages must NOT run again.
	s3 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":0,"reminders_moved":1,"reminders_merged":0,"pests_moved":0}`
	expectStage(mock, 1, model.MergeStageReminders, s3,
		stageStmt{`UPDATE care_reminders k`, nil, 0},
		stageStmt{`UPDATE user_gardens g`, nil, 0},
		stageStmt{`DELETE r FROM care_reminders r`, nil, 0},
		stageStmt{`UPDATE care_reminders SET plant_species_id = ? WHERE plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 1},
	)
	s4 := `{"favorites_moved":2,"favorites_deduped":1,"gardens_moved":1,"gardens_deduped":0,"reminders_moved":1,"reminders_merged":0,"pests_moved":4}`
	expectStage(mock, 1, model.MergeStagePests, s4,
		stageStmt{`UPDATE disease_pests SET plant_species_id = ? WHERE plant_species_id = ?`, []driver.Value{uint(10), uint(20)}, 4},
	)
	expectStage(mock, 1, model.MergeStageRetire, s4,
		stageStmt{`DELETE FROM plant_species WHERE id = ?`, []driver.Value{uint(20)}, 1},
	)
	expectTaskUpdate(mock) // MarkDone
	mergeTaskRow(mock, 1, 10, 20, model.MergeDone, model.MergeStageRetire, s4)

	task, err := svc.execute(1)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if task.Status != model.MergeDone {
		t.Fatalf("expected done, got %s", task.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestMergeExecuteFailureMarksFailed verifies a stage error rolls back its
// transaction and persists status=failed so a later retry resumes there.
func TestMergeExecuteFailureMarksFailed(t *testing.T) {
	svc, mock := newMergeService(t)
	mergeTaskRow(mock, 1, 10, 20, model.MergeRunning, "", "{}")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE f FROM favorites")).
		WithArgs(uint(10), uint(20)).
		WillReturnError(errors.New("deadlock found"))
	mock.ExpectRollback()
	expectTaskUpdate(mock) // MarkFailed
	mergeTaskRow(mock, 1, 10, 20, model.MergeFailed, "", "{}")

	task, err := svc.execute(1)
	if err != nil {
		t.Fatalf("execute should return the failed task, got error: %v", err)
	}
	if task.Status != model.MergeFailed {
		t.Fatalf("expected failed, got %s", task.Status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestMergeOneReplaysDoneTask verifies a repeated submission of a finished
// pair returns the completed result without claiming or rerunning anything.
func TestMergeOneReplaysDoneTask(t *testing.T) {
	svc, mock := newMergeService(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `plant_merge_tasks` WHERE keep_id = ? AND source_id = ?")).
		WithArgs(uint(10), uint(20), 1).
		WillReturnRows(sqlmock.NewRows(mergeTaskColumns).
			AddRow(7, 10, 20, "保留卡", "待合并卡", model.MergeDone, model.MergeStageRetire,
				`{"pests_moved":2}`, "", uint(1), time.Now(), time.Now(), time.Now()))

	keep := &model.PlantSpecies{ID: 10, Name: "保留卡"}
	task, err := svc.mergeOne(1, keep, 20)
	if err != nil {
		t.Fatalf("mergeOne: %v", err)
	}
	if task.ID != 7 || task.Status != model.MergeDone {
		t.Fatalf("expected replay of done task 7, got %+v", task)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

// TestMergeRejectsKeepEqualsSource ensures the batch validation fails fast.
func TestMergeRejectsKeepEqualsSource(t *testing.T) {
	svc, _ := newMergeService(t)
	if _, err := svc.Merge(1, 5, []uint{5}); err == nil {
		t.Fatal("expected validation error")
	}
}
