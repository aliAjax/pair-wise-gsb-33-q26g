package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// PlantMergeRepository handles persistence of plant merge tasks and the
// transactional migration units that re-point associations to the kept card.
type PlantMergeRepository struct {
	db *gorm.DB
}

// NewPlantMergeRepository creates a PlantMergeRepository.
func NewPlantMergeRepository(db *gorm.DB) *PlantMergeRepository {
	return &PlantMergeRepository{db: db}
}

// Create inserts a merge task; the unique task key turns duplicates into ErrDuplicate.
func (r *PlantMergeRepository) Create(t *model.PlantMergeTask) error {
	if err := r.db.Create(t).Error; err != nil {
		if isDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

// FindByKey locates a merge task by its idempotency key.
func (r *PlantMergeRepository) FindByKey(key string) (*model.PlantMergeTask, error) {
	var t model.PlantMergeTask
	if err := r.db.Where("task_key = ?", key).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// FindByID locates a merge task by primary key.
func (r *PlantMergeRepository) FindByID(id uint) (*model.PlantMergeTask, error) {
	var t model.PlantMergeTask
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// Save persists all fields of a merge task.
func (r *PlantMergeRepository) Save(t *model.PlantMergeTask) error {
	return r.db.Save(t).Error
}

// List returns merge tasks ordered by recency with pagination.
func (r *PlantMergeRepository) List(page, pageSize int) ([]model.PlantMergeTask, int64, error) {
	var items []model.PlantMergeTask
	var total int64
	if err := r.db.Model(&model.PlantMergeTask{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Claim atomically moves a task into processing state. Pending and failed
// tasks can always be claimed; a processing task can only be claimed once it
// is stale (its executor crashed), so a retry resumes instead of hanging.
func (r *PlantMergeRepository) Claim(id uint, staleBefore time.Time) (bool, error) {
	res := r.db.Model(&model.PlantMergeTask{}).
		Where("id = ? AND (status IN ? OR (status = ? AND updated_at < ?))",
			id,
			[]string{model.MergeTaskPending, model.MergeTaskFailed},
			model.MergeTaskProcessing, staleBefore).
		Update("status", model.MergeTaskProcessing)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// CountFavoritesOnPlants counts plant favorites on the given cards.
func (r *PlantMergeRepository) CountFavoritesOnPlants(ids []uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.Favorite{}).
		Where("target_type = ? AND target_id IN ?", constants.FavoriteTargetPlant, ids).
		Count(&n).Error
	return n, err
}

// CountFavoriteConflicts counts source-card favorites whose user already
// favorites the keep card; those rows will be deduped away.
func (r *PlantMergeRepository) CountFavoriteConflicts(keepID uint, sourceIDs []uint) (int64, error) {
	var n int64
	sub := r.db.Model(&model.Favorite{}).
		Select("user_id").
		Where("target_type = ? AND target_id = ?", constants.FavoriteTargetPlant, keepID)
	err := r.db.Model(&model.Favorite{}).
		Where("target_type = ? AND target_id IN ? AND user_id IN (?)", constants.FavoriteTargetPlant, sourceIDs, sub).
		Count(&n).Error
	return n, err
}

// CountGardensOnPlants counts garden items on the given cards.
func (r *PlantMergeRepository) CountGardensOnPlants(ids []uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.UserGarden{}).Where("plant_species_id IN ?", ids).Count(&n).Error
	return n, err
}

// CountGardenConflicts counts source-card garden items whose user already
// grows the keep card; those rows will be deduped away.
func (r *PlantMergeRepository) CountGardenConflicts(keepID uint, sourceIDs []uint) (int64, error) {
	var n int64
	sub := r.db.Model(&model.UserGarden{}).
		Select("user_id").
		Where("plant_species_id = ?", keepID)
	err := r.db.Model(&model.UserGarden{}).
		Where("plant_species_id IN ? AND user_id IN (?)", sourceIDs, sub).
		Count(&n).Error
	return n, err
}

// CountRemindersOnPlants counts care reminders on the given cards.
func (r *PlantMergeRepository) CountRemindersOnPlants(ids []uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.CareReminder{}).Where("plant_species_id IN ?", ids).Count(&n).Error
	return n, err
}

// CountReminderDuplicates counts reminders that will be merged away: within
// each (user, task title) group spanning the keep and source cards, all but
// the nearest-dated reminder are redundant.
func (r *PlantMergeRepository) CountReminderDuplicates(keepID uint, sourceIDs []uint) (int64, error) {
	var n int64
	ids := append([]uint{keepID}, sourceIDs...)
	err := r.db.Raw(`SELECT COALESCE(SUM(cnt - 1), 0) FROM (
		SELECT COUNT(*) AS cnt FROM care_reminders
		WHERE plant_species_id IN ?
		GROUP BY user_id, task_title HAVING cnt > 1) t`, ids).
		Scan(&n).Error
	return n, err
}

// CountPestsOnPlants counts disease/pest entries on the given cards.
func (r *PlantMergeRepository) CountPestsOnPlants(ids []uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.DiseasePest{}).Where("plant_species_id IN ?", ids).Count(&n).Error
	return n, err
}

// DedupeFavoritesTx removes source-card favorites of users who already
// favorite the keep card, so the re-point never violates the unique index.
func (r *PlantMergeRepository) DedupeFavoritesTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	var userIDs []uint
	if err := tx.Model(&model.Favorite{}).
		Where("target_type = ? AND target_id = ?", constants.FavoriteTargetPlant, keepID).
		Pluck("user_id", &userIDs).Error; err != nil {
		return 0, err
	}
	if len(userIDs) == 0 {
		return 0, nil
	}
	res := tx.Where("target_type = ? AND target_id = ? AND user_id IN ?",
		constants.FavoriteTargetPlant, sourceID, userIDs).
		Delete(&model.Favorite{})
	return res.RowsAffected, res.Error
}

// MoveFavoritesTx re-points the remaining source-card favorites to the keep card.
func (r *PlantMergeRepository) MoveFavoritesTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	res := tx.Model(&model.Favorite{}).
		Where("target_type = ? AND target_id = ?", constants.FavoriteTargetPlant, sourceID).
		Update("target_id", keepID)
	return res.RowsAffected, res.Error
}

// DedupeGardensTx removes source-card garden items of users who already grow
// the keep card, so the re-point never violates the unique index.
func (r *PlantMergeRepository) DedupeGardensTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	var userIDs []uint
	if err := tx.Model(&model.UserGarden{}).
		Where("plant_species_id = ?", keepID).
		Pluck("user_id", &userIDs).Error; err != nil {
		return 0, err
	}
	if len(userIDs) == 0 {
		return 0, nil
	}
	res := tx.Where("plant_species_id = ? AND user_id IN ?", sourceID, userIDs).
		Delete(&model.UserGarden{})
	return res.RowsAffected, res.Error
}

// MoveGardensTx re-points the remaining source-card garden items to the keep card.
func (r *PlantMergeRepository) MoveGardensTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	res := tx.Model(&model.UserGarden{}).
		Where("plant_species_id = ?", sourceID).
		Update("plant_species_id", keepID)
	return res.RowsAffected, res.Error
}

// MoveRemindersTx re-points source-card reminders to the keep card.
func (r *PlantMergeRepository) MoveRemindersTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	res := tx.Model(&model.CareReminder{}).
		Where("plant_species_id = ?", sourceID).
		Update("plant_species_id", keepID)
	return res.RowsAffected, res.Error
}

// DedupeRemindersTx merges duplicate reminders on the keep card: within each
// (user, task title) group only the nearest-dated reminder survives (ties keep
// the smallest id). Running it repeatedly is a no-op, which keeps retries safe.
func (r *PlantMergeRepository) DedupeRemindersTx(tx *gorm.DB, keepID uint) (int64, error) {
	var rows []model.CareReminder
	if err := tx.Where("plant_species_id = ?", keepID).
		Order("user_id, task_title, remind_date ASC, id ASC").
		Find(&rows).Error; err != nil {
		return 0, err
	}
	type groupKey struct {
		userID uint
		title  string
	}
	seen := make(map[groupKey]bool, len(rows))
	dropIDs := make([]uint, 0)
	for _, row := range rows {
		key := groupKey{userID: row.UserID, title: row.TaskTitle}
		if seen[key] {
			dropIDs = append(dropIDs, row.ID)
			continue
		}
		seen[key] = true
	}
	if len(dropIDs) == 0 {
		return 0, nil
	}
	res := tx.Where("id IN ?", dropIDs).Delete(&model.CareReminder{})
	return res.RowsAffected, res.Error
}

// MovePestsTx re-points source-card disease/pest entries to the keep card.
func (r *PlantMergeRepository) MovePestsTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	res := tx.Model(&model.DiseasePest{}).
		Where("plant_species_id = ?", sourceID).
		Update("plant_species_id", keepID)
	return res.RowsAffected, res.Error
}

// RetirePlantTx removes the merged source card from the library.
func (r *PlantMergeRepository) RetirePlantTx(tx *gorm.DB, sourceID uint) (int64, error) {
	res := tx.Delete(&model.PlantSpecies{}, sourceID)
	return res.RowsAffected, res.Error
}

// UpdateCheckpointTx persists the checkpoint inside the unit transaction, so
// progress and data changes commit or roll back together.
func (r *PlantMergeRepository) UpdateCheckpointTx(tx *gorm.DB, taskID uint, checkpoint string) error {
	return tx.Model(&model.PlantMergeTask{}).
		Where("id = ?", taskID).
		Update("checkpoint", checkpoint).Error
}
