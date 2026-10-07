package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// mergeStaleTimeout is how long a processing task may stay untouched before a
// retry is allowed to claim it (covers executor crashes).
const mergeStaleTimeout = 5 * time.Minute

// mergeExecMu serializes merge execution inside one backend process so two
// different batches can never re-point the same user's rows concurrently.
var mergeExecMu sync.Mutex

// PlantMergeService implements plant species merge with checkpoint recovery.
type PlantMergeService struct {
	db        *gorm.DB
	repo      *repository.PlantMergeRepository
	plantRepo *repository.PlantSpeciesRepository
	logger    *slog.Logger
}

// NewPlantMergeService creates a PlantMergeService.
func NewPlantMergeService(db *gorm.DB, repo *repository.PlantMergeRepository, plantRepo *repository.PlantSpeciesRepository, logger *slog.Logger) *PlantMergeService {
	return &PlantMergeService{db: db, repo: repo, plantRepo: plantRepo, logger: logger}
}

// normalizeSourceIDs dedupes and sorts source ids and collects validation issues.
func normalizeSourceIDs(keepID uint, sourceIDs []uint) ([]uint, []string) {
	seen := make(map[uint]bool, len(sourceIDs))
	ids := make([]uint, 0, len(sourceIDs))
	issues := []string{}
	for _, id := range sourceIDs {
		if id == 0 {
			issues = append(issues, "PlantMerge source_ids contains invalid id=0")
			continue
		}
		if id == keepID {
			issues = append(issues, fmt.Sprintf("PlantSpecies[id=%d] merge failed: keep card cannot be a source card", id))
			continue
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	if len(ids) == 0 {
		issues = append(issues, "PlantMerge source_ids empty: at least one source card is required")
	}
	return ids, issues
}

// Precheck validates a merge batch and estimates how each association type
// will be affected. Issues are reported in the response instead of an error
// so the page can show them next to the stats.
func (s *PlantMergeService) Precheck(keepID uint, sourceIDs []uint) (*model.MergePrecheck, error) {
	ids, issues := normalizeSourceIDs(keepID, sourceIDs)
	pre := &model.MergePrecheck{
		KeepID:      keepID,
		SourceIDs:   ids,
		SourceNames: []string{},
		Issues:      issues,
	}
	keep, err := s.plantRepo.FindByID(keepID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			pre.Issues = append(pre.Issues, fmt.Sprintf("PlantSpecies[id=%d] keep card not found", keepID))
		} else {
			return nil, fmt.Errorf("plant merge precheck keep find: %w", err)
		}
	} else {
		pre.KeepName = keep.Name
	}
	validSources := make([]uint, 0, len(ids))
	for _, id := range ids {
		p, err := s.plantRepo.FindByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				pre.Issues = append(pre.Issues, fmt.Sprintf("PlantSpecies[id=%d] source card not found or already retired", id))
				continue
			}
			return nil, fmt.Errorf("plant merge precheck source find: %w", err)
		}
		pre.SourceNames = append(pre.SourceNames, p.Name)
		validSources = append(validSources, id)
	}
	if len(validSources) > 0 {
		if err := s.fillPrecheckStats(pre, keepID, validSources); err != nil {
			return nil, err
		}
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantMergePrecheck, keepID, ids))
	return pre, nil
}

// fillPrecheckStats queries the per-association estimates.
func (s *PlantMergeService) fillPrecheckStats(pre *model.MergePrecheck, keepID uint, sourceIDs []uint) error {
	var err error
	if pre.Stats.FavoritesTotal, err = s.repo.CountFavoritesOnPlants(sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck favorites count: %w", err)
	}
	if pre.Stats.FavoritesDedupe, err = s.repo.CountFavoriteConflicts(keepID, sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck favorites conflict: %w", err)
	}
	if pre.Stats.GardensTotal, err = s.repo.CountGardensOnPlants(sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck gardens count: %w", err)
	}
	if pre.Stats.GardensDedupe, err = s.repo.CountGardenConflicts(keepID, sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck gardens conflict: %w", err)
	}
	if pre.Stats.RemindersTotal, err = s.repo.CountRemindersOnPlants(sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck reminders count: %w", err)
	}
	if pre.Stats.RemindersMerge, err = s.repo.CountReminderDuplicates(keepID, sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck reminders duplicates: %w", err)
	}
	if pre.Stats.PestsTotal, err = s.repo.CountPestsOnPlants(sourceIDs); err != nil {
		return fmt.Errorf("plant merge precheck pests count: %w", err)
	}
	return nil
}

// Merge submits a plant merge batch. Submission is idempotent: the same keep
// card plus source set maps to one task, so a concurrent duplicate window
// receives the existing task and sees the completed result.
func (s *PlantMergeService) Merge(operatorID, keepID uint, sourceIDs []uint) (*model.PlantMergeTask, error) {
	ids, issues := normalizeSourceIDs(keepID, sourceIDs)
	if len(issues) > 0 {
		return nil, util.NewAppError(422, constants.CodeValidationError, strings.Join(issues, "; "))
	}

	// The idempotency lookup comes first: once a batch is merged, its source
	// cards are retired, so a replayed submission can only be recognized by
	// its task key, never by re-validating the cards.
	key := model.BuildMergeTaskKey(keepID, ids)
	task, err := s.repo.FindByKey(key)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("plant merge task find: %w", err)
	}
	if errors.Is(err, repository.ErrNotFound) {
		if _, err := s.plantRepo.FindByID(keepID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.NewAppError(404, constants.CodeNotFound,
					fmt.Sprintf("PlantSpecies[id=%d] keep card not found", keepID))
			}
			return nil, fmt.Errorf("plant merge keep find: %w", err)
		}
		for _, id := range ids {
			if _, err := s.plantRepo.FindByID(id); err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return nil, util.NewAppError(404, constants.CodeNotFound,
						fmt.Sprintf("PlantSpecies[id=%d] source card not found or already retired", id))
				}
				return nil, fmt.Errorf("plant merge source find: %w", err)
			}
		}
		task = &model.PlantMergeTask{
			TaskKey:    key,
			KeepID:     keepID,
			SourceIDs:  model.EncodeSourceIDs(ids),
			Status:     model.MergeTaskPending,
			Checkpoint: "{}",
			Precheck:   "{}",
			Result:     "{}",
			OperatorID: operatorID,
		}
		if pre, perr := s.Precheck(keepID, ids); perr == nil {
			if b, merr := json.Marshal(pre); merr == nil {
				task.Precheck = string(b)
			}
		}
		if cerr := s.repo.Create(task); cerr != nil {
			if !errors.Is(cerr, repository.ErrDuplicate) {
				return nil, fmt.Errorf("plant merge task create: %w", cerr)
			}
			// Another window created the same batch first; fall through to it.
			s.logger.Info(fmt.Sprintf(constants.LogPlantMergeDuplicate, key))
			task, err = s.repo.FindByKey(key)
			if err != nil {
				return nil, fmt.Errorf("plant merge task refind: %w", err)
			}
		} else {
			s.logger.Info(fmt.Sprintf(constants.LogPlantMergeSubmit, task.ID, keepID, task.SourceIDs))
		}
	}

	switch task.Status {
	case model.MergeTaskSucceeded:
		// First submitter already finished: the late window sees the result.
		return task, nil
	case model.MergeTaskProcessing:
		if time.Since(task.UpdatedAt) < mergeStaleTimeout {
			return task, nil
		}
	}

	claimed, err := s.repo.Claim(task.ID, time.Now().Add(-mergeStaleTimeout))
	if err != nil {
		return nil, fmt.Errorf("plant merge task claim: %w", err)
	}
	if !claimed {
		// Lost the claim race; report the winner's current state.
		return s.repo.FindByID(task.ID)
	}
	if task.Status == model.MergeTaskFailed || task.Status == model.MergeTaskProcessing {
		s.logger.Info(fmt.Sprintf(constants.LogPlantMergeResume, task.ID))
	}
	task.Status = model.MergeTaskProcessing

	mergeExecMu.Lock()
	defer mergeExecMu.Unlock()
	if err := s.execute(task); err != nil {
		return task, err
	}
	return s.repo.FindByID(task.ID)
}

// execute runs the unfinished (phase, source) units of a claimed task and
// checkpoints each one inside its own transaction.
func (s *PlantMergeService) execute(task *model.PlantMergeTask) error {
	if _, err := s.plantRepo.FindByID(task.KeepID); err != nil {
		return s.failTask(task, fmt.Sprintf("PlantSpecies[id=%d] keep card not found, merge aborted", task.KeepID))
	}
	cp := task.CheckpointData()
	for _, phase := range model.MergePhases() {
		for _, sourceID := range task.SourceIDList() {
			if cp.Done(phase, sourceID) {
				continue
			}
			if err := s.runUnit(task, cp, phase, sourceID); err != nil {
				return s.failTask(task, err.Error())
			}
			// Keep the in-memory checkpoint in sync so a later Save (failure
			// mark or finalize) never clobbers the committed progress.
			task.Checkpoint = cp.JSON()
			s.logger.Info(fmt.Sprintf(constants.LogPlantMergeUnitDone, task.ID, phase, sourceID))
		}
	}
	resultJSON, err := json.Marshal(cp.Result)
	if err != nil {
		return s.failTask(task, fmt.Sprintf("result encode: %v", err))
	}
	task.Status = model.MergeTaskSucceeded
	task.Result = string(resultJSON)
	task.LastError = ""
	if err := s.repo.Save(task); err != nil {
		return fmt.Errorf("plant merge task finalize: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantMergeSuccess, task.ID, task.KeepID))
	return nil
}

// runUnit executes one (phase, source) unit and checkpoints it atomically.
func (s *PlantMergeService) runUnit(task *model.PlantMergeTask, cp *model.MergeCheckpoint, phase string, sourceID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		switch phase {
		case model.MergePhaseFavorites:
			var deduped, moved int64
			if deduped, err = s.repo.DedupeFavoritesTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("favorites dedupe source=%d: %w", sourceID, err)
			}
			if moved, err = s.repo.MoveFavoritesTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("favorites move source=%d: %w", sourceID, err)
			}
			cp.Result.FavoritesDeduped += deduped
			cp.Result.FavoritesMoved += moved
		case model.MergePhaseGardens:
			var deduped, moved int64
			if deduped, err = s.repo.DedupeGardensTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("gardens dedupe source=%d: %w", sourceID, err)
			}
			if moved, err = s.repo.MoveGardensTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("gardens move source=%d: %w", sourceID, err)
			}
			cp.Result.GardensDeduped += deduped
			cp.Result.GardensMoved += moved
		case model.MergePhaseReminders:
			var moved, merged int64
			if moved, err = s.repo.MoveRemindersTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("reminders move source=%d: %w", sourceID, err)
			}
			if merged, err = s.repo.DedupeRemindersTx(tx, task.KeepID); err != nil {
				return fmt.Errorf("reminders merge source=%d: %w", sourceID, err)
			}
			cp.Result.RemindersMoved += moved
			cp.Result.RemindersMerged += merged
		case model.MergePhasePests:
			var moved int64
			if moved, err = s.repo.MovePestsTx(tx, task.KeepID, sourceID); err != nil {
				return fmt.Errorf("pests move source=%d: %w", sourceID, err)
			}
			cp.Result.PestsMoved += moved
		case model.MergePhaseRetire:
			var retired int64
			if retired, err = s.repo.RetirePlantTx(tx, sourceID); err != nil {
				return fmt.Errorf("retire source=%d: %w", sourceID, err)
			}
			cp.Result.RetiredPlants += retired
		default:
			return fmt.Errorf("unknown merge phase %q", phase)
		}
		cp.Mark(phase, sourceID)
		if err := s.repo.UpdateCheckpointTx(tx, task.ID, cp.JSON()); err != nil {
			return fmt.Errorf("checkpoint update phase=%s source=%d: %w", phase, sourceID, err)
		}
		return nil
	})
}

// failTask marks the task failed so a later retry resumes from the checkpoint.
func (s *PlantMergeService) failTask(task *model.PlantMergeTask, msg string) error {
	task.Status = model.MergeTaskFailed
	task.LastError = msg
	if err := s.repo.Save(task); err != nil {
		return fmt.Errorf("plant merge task fail mark: %w", err)
	}
	s.logger.Error(fmt.Sprintf(constants.LogPlantMergeFailed, task.ID, msg))
	return util.NewAppError(500, constants.CodeInternalError,
		fmt.Sprintf("PlantMergeTask[id=%d] %s: %s", task.ID, constants.MsgPlantMergeFailed, msg))
}

// GetTask returns a merge task by id for status polling.
func (s *PlantMergeService) GetTask(id uint) (*model.PlantMergeTask, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("PlantMergeTask[id=%d] not found", id))
		}
		return nil, fmt.Errorf("plant merge task get: %w", err)
	}
	return task, nil
}

// ListTasks returns merge tasks ordered by recency.
func (s *PlantMergeService) ListTasks(page, pageSize int) ([]model.PlantMergeTask, int64, error) {
	items, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("plant merge task list: %w", err)
	}
	return items, total, nil
}
