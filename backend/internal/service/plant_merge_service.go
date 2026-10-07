package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/dto"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

const (
	// mergeStaleThreshold marks a running task as crashed so it can be reclaimed.
	mergeStaleThreshold = 2 * time.Minute
	// mergeWaitTimeout bounds how long a late submission waits for the
	// concurrent executor before returning the still-running task.
	mergeWaitTimeout = 10 * time.Second
	mergeWaitPoll    = 200 * time.Millisecond
)

// PlantMergeService merges duplicate variety cards: the four association
// types move to the keep card, the source card retires. Progress is
// checkpointed per stage so a failed migration resumes where it stopped.
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

// Precheck reports how many favorites, garden entries, reminders and pest
// entries each source card holds, and how many of them will be deduplicated
// or folded because the user already relates to the keep card.
func (s *PlantMergeService) Precheck(keepID uint, sourceIDs []uint) (*dto.MergePrecheckResult, error) {
	if err := validateSourceIDs(keepID, sourceIDs); err != nil {
		return nil, err
	}
	keep, err := s.plantRepo.FindByID(keepID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("PlantSpecies[id=%d] keep card not found", keepID))
		}
		return nil, fmt.Errorf("plant merge precheck find keep: %w", err)
	}
	result := &dto.MergePrecheckResult{KeepID: keep.ID, KeepName: keep.Name, Items: []dto.MergePrecheckItem{}}
	for _, sourceID := range sourceIDs {
		item, err := s.precheckOne(keep.ID, sourceID)
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, *item)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantMergePrecheck, keep.ID, len(sourceIDs)))
	return result, nil
}

// precheckOne builds the report for one source card. A source whose pair task
// already completed is reported as already merged instead of an error, so a
// late window re-running precheck on the same batch sees the done state.
func (s *PlantMergeService) precheckOne(keepID, sourceID uint) (*dto.MergePrecheckItem, error) {
	if task, err := s.repo.FindByPair(keepID, sourceID); err == nil && task.Status == model.MergeDone {
		return &dto.MergePrecheckItem{SourceID: sourceID, SourceName: task.SourceName, AlreadyMerged: true}, nil
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("plant merge precheck find task: %w", err)
	}
	source, err := s.plantRepo.FindByID(sourceID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("PlantSpecies[id=%d] source card not found", sourceID))
		}
		return nil, fmt.Errorf("plant merge precheck find source: %w", err)
	}
	item := &dto.MergePrecheckItem{SourceID: sourceID, SourceName: source.Name}
	counts := []struct {
		dst *int64
		fn  func() (int64, error)
	}{
		{&item.FavoritesTotal, func() (int64, error) { return s.repo.CountFavorites(sourceID) }},
		{&item.FavoritesDedup, func() (int64, error) { return s.repo.CountFavoriteOverlaps(keepID, sourceID) }},
		{&item.GardensTotal, func() (int64, error) { return s.repo.CountGardens(sourceID) }},
		{&item.GardensDedup, func() (int64, error) { return s.repo.CountGardenOverlaps(keepID, sourceID) }},
		{&item.RemindersTotal, func() (int64, error) { return s.repo.CountReminders(sourceID) }},
		{&item.RemindersMerge, func() (int64, error) { return s.repo.CountReminderOverlaps(keepID, sourceID) }},
		{&item.PestsTotal, func() (int64, error) { return s.repo.CountPests(sourceID) }},
	}
	for _, c := range counts {
		n, err := c.fn()
		if err != nil {
			return nil, fmt.Errorf("plant merge precheck count: %w", err)
		}
		*c.dst = n
	}
	return item, nil
}

// Merge executes (or resumes, or replays) one merge task per source card.
// The first submission of a (keep, source) pair runs the migration; a
// concurrent duplicate submission follows the existing task and observes its
// finished result, so the same batch is never applied twice.
func (s *PlantMergeService) Merge(operatorID, keepID uint, sourceIDs []uint) ([]dto.PlantMergeResult, error) {
	if err := validateSourceIDs(keepID, sourceIDs); err != nil {
		return nil, err
	}
	keep, err := s.plantRepo.FindByID(keepID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("PlantSpecies[id=%d] keep card not found", keepID))
		}
		return nil, fmt.Errorf("plant merge find keep: %w", err)
	}
	// Validate every source before executing anything, so a bad id fails the
	// whole batch instead of leaving a half-submitted batch behind.
	for _, sourceID := range sourceIDs {
		if _, err := s.repo.FindByPair(keepID, sourceID); err == nil {
			continue
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("plant merge find task: %w", err)
		}
		if _, err := s.plantRepo.FindByID(sourceID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.NewAppError(404, constants.CodeNotFound,
					fmt.Sprintf("PlantSpecies[id=%d] source card not found", sourceID))
			}
			return nil, fmt.Errorf("plant merge find source: %w", err)
		}
	}
	results := make([]dto.PlantMergeResult, 0, len(sourceIDs))
	for _, sourceID := range sourceIDs {
		task, err := s.mergeOne(operatorID, keep, sourceID)
		if err != nil {
			return nil, err
		}
		results = append(results, toMergeResult(task))
	}
	return results, nil
}

// mergeOne gets or creates the pair task and drives it to a terminal state.
func (s *PlantMergeService) mergeOne(operatorID uint, keep *model.PlantSpecies, sourceID uint) (*model.PlantMergeTask, error) {
	task, err := s.repo.FindByPair(keep.ID, sourceID)
	if errors.Is(err, repository.ErrNotFound) {
		source, findErr := s.plantRepo.FindByID(sourceID)
		if findErr != nil {
			return nil, fmt.Errorf("plant merge find source: %w", findErr)
		}
		task = &model.PlantMergeTask{
			KeepID: keep.ID, SourceID: sourceID, KeepName: keep.Name, SourceName: source.Name,
			Status: model.MergePending, Stats: "{}", OperatorID: operatorID,
		}
		if createErr := s.repo.Create(task); createErr != nil {
			if !errors.Is(createErr, repository.ErrDuplicate) {
				return nil, fmt.Errorf("plant merge task create: %w", createErr)
			}
			// A concurrent window submitted the same pair first; follow its task.
			task, err = s.repo.FindByPair(keep.ID, sourceID)
			if err != nil {
				return nil, fmt.Errorf("plant merge task refetch: %w", err)
			}
		}
	} else if err != nil {
		return nil, fmt.Errorf("plant merge find task: %w", err)
	}
	if task.Status == model.MergeDone {
		s.logger.Info(fmt.Sprintf(constants.LogPlantMergeReplay, task.ID))
		return task, nil
	}
	claimed, err := s.repo.Claim(task.ID, time.Now().Add(-mergeStaleThreshold))
	if err != nil {
		return nil, fmt.Errorf("plant merge claim: %w", err)
	}
	if !claimed {
		s.logger.Info(fmt.Sprintf(constants.LogPlantMergeWait, task.ID))
		return s.waitForTask(task.ID)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantMergeStart, task.ID, task.KeepID, task.SourceID))
	return s.execute(task.ID)
}

// execute runs the remaining stages after the checkpoint, each in its own
// transaction together with the checkpoint update. A stage failure persists
// status=failed; the next submission of the same pair resumes from Stage.
func (s *PlantMergeService) execute(taskID uint) (*model.PlantMergeTask, error) {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("plant merge execute find task: %w", err)
	}
	stats := parseMergeStats(task.Stats)
	for i := stageIndex(task.Stage) + 1; i < len(model.MergeStageOrder); i++ {
		stage := model.MergeStageOrder[i]
		err := s.db.Transaction(func(tx *gorm.DB) error {
			if err := s.runStage(tx, stage, task.KeepID, task.SourceID, &stats); err != nil {
				return err
			}
			return s.repo.SaveCheckpointTx(tx, task.ID, stage, marshalMergeStats(&stats))
		})
		if err != nil {
			lastErr := fmt.Sprintf("stage %s: %v", stage, err)
			if markErr := s.repo.MarkFailed(task.ID, lastErr); markErr != nil {
				return nil, fmt.Errorf("plant merge mark failed: %w", markErr)
			}
			s.logger.Error(fmt.Sprintf(constants.LogPlantMergeFailed, task.ID, stage), "error", err)
			return s.repo.FindByID(task.ID)
		}
		s.logger.Info(fmt.Sprintf(constants.LogPlantMergeStageDone, task.ID, stage))
	}
	if err := s.repo.MarkDone(task.ID, marshalMergeStats(&stats)); err != nil {
		return nil, fmt.Errorf("plant merge mark done: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPlantMergeSuccess, task.ID, task.KeepID, task.SourceID))
	return s.repo.FindByID(task.ID)
}

// runStage applies one migration stage inside tx and folds its counts into stats.
func (s *PlantMergeService) runStage(tx *gorm.DB, stage string, keepID, sourceID uint, stats *dto.MergeStats) error {
	switch stage {
	case model.MergeStageFavorites:
		moved, deduped, err := s.repo.MigrateFavoritesTx(tx, keepID, sourceID)
		stats.FavoritesMoved += int(moved)
		stats.FavoritesDeduped += int(deduped)
		return err
	case model.MergeStageGardens:
		moved, deduped, err := s.repo.MigrateGardensTx(tx, keepID, sourceID)
		stats.GardensMoved += int(moved)
		stats.GardensDeduped += int(deduped)
		return err
	case model.MergeStageReminders:
		moved, merged, err := s.repo.MigrateRemindersTx(tx, keepID, sourceID)
		stats.RemindersMoved += int(moved)
		stats.RemindersMerged += int(merged)
		return err
	case model.MergeStagePests:
		moved, err := s.repo.MigratePestsTx(tx, keepID, sourceID)
		stats.PestsMoved += int(moved)
		return err
	case model.MergeStageRetire:
		return s.repo.RetireSourceTx(tx, sourceID)
	default:
		return fmt.Errorf("unknown merge stage %q", stage)
	}
}

// waitForTask polls a task owned by a concurrent executor until it reaches a
// terminal state or the wait budget runs out.
func (s *PlantMergeService) waitForTask(taskID uint) (*model.PlantMergeTask, error) {
	deadline := time.Now().Add(mergeWaitTimeout)
	for {
		task, err := s.repo.FindByID(taskID)
		if err != nil {
			return nil, fmt.Errorf("plant merge wait find task: %w", err)
		}
		if task.Status == model.MergeDone || task.Status == model.MergeFailed || time.Now().After(deadline) {
			return task, nil
		}
		time.Sleep(mergeWaitPoll)
	}
}

// GetTask returns one merge task for status polling.
func (s *PlantMergeService) GetTask(id uint) (*dto.PlantMergeResult, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("PlantMergeTask[id=%d] not found", id))
		}
		return nil, fmt.Errorf("plant merge get task: %w", err)
	}
	result := toMergeResult(task)
	return &result, nil
}

// ListTasks returns recent merge tasks.
func (s *PlantMergeService) ListTasks(page, pageSize int) ([]dto.PlantMergeResult, int64, error) {
	tasks, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("plant merge list tasks: %w", err)
	}
	results := make([]dto.PlantMergeResult, 0, len(tasks))
	for i := range tasks {
		results = append(results, toMergeResult(&tasks[i]))
	}
	return results, total, nil
}

// validateSourceIDs rejects empty batches, duplicates and keep==source.
func validateSourceIDs(keepID uint, sourceIDs []uint) error {
	if len(sourceIDs) == 0 {
		return util.NewAppError(422, constants.CodeValidationError, "plant merge failed: source_ids is empty")
	}
	seen := make(map[uint]bool, len(sourceIDs))
	for _, id := range sourceIDs {
		if id == keepID {
			return util.NewAppError(422, constants.CodeValidationError,
				fmt.Sprintf("plant merge failed: source id=%d equals keep card", id))
		}
		if seen[id] {
			return util.NewAppError(422, constants.CodeValidationError,
				fmt.Sprintf("plant merge failed: duplicate source id=%d", id))
		}
		seen[id] = true
	}
	return nil
}

// stageIndex returns the position of a completed stage, or -1 when none ran.
func stageIndex(stage string) int {
	for i, s := range model.MergeStageOrder {
		if s == stage {
			return i
		}
	}
	return -1
}

// parseMergeStats decodes the accumulated stats JSON, tolerating empty input.
func parseMergeStats(raw string) dto.MergeStats {
	var stats dto.MergeStats
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &stats)
	}
	return stats
}

// marshalMergeStats encodes stats for the checkpoint column.
func marshalMergeStats(stats *dto.MergeStats) string {
	b, _ := json.Marshal(stats)
	return string(b)
}

// toMergeResult maps a task model to its API view.
func toMergeResult(t *model.PlantMergeTask) dto.PlantMergeResult {
	stats := parseMergeStats(t.Stats)
	return dto.PlantMergeResult{
		ID: t.ID, KeepID: t.KeepID, SourceID: t.SourceID,
		KeepName: t.KeepName, SourceName: t.SourceName,
		Status: t.Status, Stage: t.Stage, Stats: &stats,
		LastError: t.LastError, OperatorID: t.OperatorID,
		CreatedAt: t.CreatedAt, FinishedAt: t.FinishedAt,
	}
}
