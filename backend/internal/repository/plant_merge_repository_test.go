package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

func TestPlantMergeRepositoryCreateDuplicate(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPlantMergeRepository(db)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `plant_merge_tasks`")).
		WillReturnError(errors.New("Duplicate entry 'K1-S2' for key 'plant_merge_tasks.task_key'"))
	mock.ExpectRollback()
	task := &model.PlantMergeTask{TaskKey: "K1-S2", KeepID: 1, SourceIDs: "[2]", Status: model.MergeTaskPending}
	if err := repo.Create(task); !errors.Is(err, ErrDuplicate) {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}

func TestPlantMergeRepositoryFindByKeyNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPlantMergeRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `plant_merge_tasks` WHERE task_key = ?")).
		WithArgs("K1-S2", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	if _, err := repo.FindByKey("K1-S2"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantMergeRepositoryClaim(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPlantMergeRepository(db)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `plant_merge_tasks` SET `status`=?")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	ok, err := repo.Claim(9, time.Now().Add(-5*time.Minute))
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}
	if !ok {
		t.Error("expected claim to succeed when one row is updated")
	}
}

func TestPlantMergeRepositoryClaimLost(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewPlantMergeRepository(db)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `plant_merge_tasks` SET `status`=?")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	ok, err := repo.Claim(9, time.Now().Add(-5*time.Minute))
	if err != nil {
		t.Fatalf("Claim: %v", err)
	}
	if ok {
		t.Error("expected claim to fail when no row matches")
	}
}
