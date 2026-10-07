package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// PlantMergeRepository persists merge tasks and executes the association
// migration SQL. Every stage statement is idempotent: rows already moved to
// the keep card no longer match the source filter, so a retry after failure
// only picks up the unfinished objects.
type PlantMergeRepository struct {
	db *gorm.DB
}

// NewPlantMergeRepository creates a PlantMergeRepository.
func NewPlantMergeRepository(db *gorm.DB) *PlantMergeRepository {
	return &PlantMergeRepository{db: db}
}

// Create inserts a merge task. The unique (keep_id, source_id) pair turns
// concurrent submissions of the same batch into ErrDuplicate for the loser.
func (r *PlantMergeRepository) Create(t *model.PlantMergeTask) error {
	if err := r.db.Create(t).Error; err != nil {
		if isDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

// FindByID locates a merge task by id.
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

// FindByPair locates a merge task by its unique keep/source pair.
func (r *PlantMergeRepository) FindByPair(keepID, sourceID uint) (*model.PlantMergeTask, error) {
	var t model.PlantMergeTask
	if err := r.db.Where("keep_id = ? AND source_id = ?", keepID, sourceID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// List returns merge tasks ordered by recency with pagination.
func (r *PlantMergeRepository) List(page, pageSize int) ([]model.PlantMergeTask, int64, error) {
	var items []model.PlantMergeTask
	var total int64
	q := r.db.Model(&model.PlantMergeTask{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Claim atomically marks a task running when it is pending/failed, or when a
// previous runner crashed and left it running stale. RowsAffected == 0 means
// another executor owns the task. GORM refreshes updated_at as the heartbeat.
func (r *PlantMergeRepository) Claim(id uint, staleBefore time.Time) (bool, error) {
	res := r.db.Model(&model.PlantMergeTask{}).
		Where("id = ? AND (status IN ? OR (status = ? AND updated_at < ?))",
			id, []string{model.MergePending, model.MergeFailed}, model.MergeRunning, staleBefore).
		Updates(map[string]interface{}{"status": model.MergeRunning})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// SaveCheckpointTx records the last finished stage and accumulated stats
// inside the same transaction as the stage migration.
func (r *PlantMergeRepository) SaveCheckpointTx(tx *gorm.DB, id uint, stage, statsJSON string) error {
	return tx.Model(&model.PlantMergeTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"stage": stage, "stats": statsJSON}).Error
}

// MarkFailed stores the failure reason so a later retry resumes from Stage.
func (r *PlantMergeRepository) MarkFailed(id uint, lastErr string) error {
	return r.db.Model(&model.PlantMergeTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": model.MergeFailed, "last_error": lastErr}).Error
}

// MarkDone finishes a task.
func (r *PlantMergeRepository) MarkDone(id uint, statsJSON string) error {
	now := time.Now()
	return r.db.Model(&model.PlantMergeTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": model.MergeDone, "stage": model.MergeStageRetire,
			"stats": statsJSON, "last_error": "", "finished_at": &now,
		}).Error
}

// MigrateFavoritesTx re-attaches plant favorites from source to keep. Users
// who favorited both cards keep exactly one row (unique user+target).
func (r *PlantMergeRepository) MigrateFavoritesTx(tx *gorm.DB, keepID, sourceID uint) (moved, deduped int64, err error) {
	res := tx.Exec(`DELETE f FROM favorites f
		JOIN favorites k ON k.user_id = f.user_id AND k.target_type = 'plant' AND k.target_id = ?
		WHERE f.target_type = 'plant' AND f.target_id = ?`, keepID, sourceID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	deduped = res.RowsAffected
	res = tx.Exec(`UPDATE favorites SET target_id = ? WHERE target_type = 'plant' AND target_id = ?`, keepID, sourceID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	return res.RowsAffected, deduped, nil
}

// MigrateGardensTx re-attaches garden entries from source to keep. Users
// growing both cards keep exactly one entry; a reminder linked only on the
// dropped entry is carried over to the surviving one first.
func (r *PlantMergeRepository) MigrateGardensTx(tx *gorm.DB, keepID, sourceID uint) (moved, deduped int64, err error) {
	if res := tx.Exec(`UPDATE user_gardens k
		JOIN user_gardens g ON g.user_id = k.user_id AND g.plant_species_id = ?
		SET k.care_reminder_id = g.care_reminder_id
		WHERE k.plant_species_id = ? AND k.care_reminder_id = 0 AND g.care_reminder_id <> 0`, sourceID, keepID); res.Error != nil {
		return 0, 0, res.Error
	}
	res := tx.Exec(`DELETE g FROM user_gardens g
		JOIN user_gardens k ON k.user_id = g.user_id AND k.plant_species_id = ?
		WHERE g.plant_species_id = ?`, keepID, sourceID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	deduped = res.RowsAffected
	res = tx.Exec(`UPDATE user_gardens SET plant_species_id = ? WHERE plant_species_id = ?`, keepID, sourceID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	return res.RowsAffected, deduped, nil
}

// MigrateRemindersTx re-attaches care reminders from source to keep. When the
// same user owns the same task on both cards the copies are folded into one
// row carrying the most recent remind date, and the task stays open unless
// both copies were done. Garden entries linked to a folded reminder are
// re-pointed to the surviving one before it is deleted.
func (r *PlantMergeRepository) MigrateRemindersTx(tx *gorm.DB, keepID, sourceID uint) (moved, merged int64, err error) {
	if res := tx.Exec(`UPDATE care_reminders k
		JOIN care_reminders r ON r.user_id = k.user_id AND r.task_title = k.task_title
		SET k.remind_date = GREATEST(k.remind_date, r.remind_date),
			k.status = IF(k.status = 'done' AND r.status = 'done', 'done', 'pending')
		WHERE k.plant_species_id = ? AND r.plant_species_id = ?`, keepID, sourceID); res.Error != nil {
		return 0, 0, res.Error
	}
	if res := tx.Exec(`UPDATE user_gardens g
		JOIN care_reminders r ON g.care_reminder_id = r.id
		JOIN care_reminders k ON k.user_id = r.user_id AND k.task_title = r.task_title
		SET g.care_reminder_id = k.id
		WHERE r.plant_species_id = ? AND k.plant_species_id = ?`, sourceID, keepID); res.Error != nil {
		return 0, 0, res.Error
	}
	res := tx.Exec(`DELETE r FROM care_reminders r
		JOIN care_reminders k ON k.user_id = r.user_id AND k.task_title = r.task_title
		WHERE r.plant_species_id = ? AND k.plant_species_id = ?`, sourceID, keepID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	merged = res.RowsAffected
	res = tx.Exec(`UPDATE care_reminders SET plant_species_id = ? WHERE plant_species_id = ?`, keepID, sourceID)
	if res.Error != nil {
		return 0, 0, res.Error
	}
	return res.RowsAffected, merged, nil
}

// MigratePestsTx re-attaches disease/pest manual entries from source to keep.
func (r *PlantMergeRepository) MigratePestsTx(tx *gorm.DB, keepID, sourceID uint) (int64, error) {
	res := tx.Exec(`UPDATE disease_pests SET plant_species_id = ? WHERE plant_species_id = ?`, keepID, sourceID)
	return res.RowsAffected, res.Error
}

// RetireSourceTx removes the merged source card from the library. All four
// association types have been re-attached by earlier stages at this point. A
// card that is already gone (e.g. removed manually after a failure) counts
// as retired so a retry can still finish the task.
func (r *PlantMergeRepository) RetireSourceTx(tx *gorm.DB, sourceID uint) error {
	return tx.Exec(`DELETE FROM plant_species WHERE id = ?`, sourceID).Error
}

// CountFavorites counts plant favorites pointing at a card.
func (r *PlantMergeRepository) CountFavorites(plantID uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.Favorite{}).
		Where("target_type = 'plant' AND target_id = ?", plantID).Count(&n).Error
	return n, err
}

// CountFavoriteOverlaps counts users who favorited both cards.
func (r *PlantMergeRepository) CountFavoriteOverlaps(keepID, sourceID uint) (int64, error) {
	var n int64
	err := r.db.Table("favorites f").
		Joins("JOIN favorites k ON k.user_id = f.user_id AND k.target_type = 'plant' AND k.target_id = ?", keepID).
		Where("f.target_type = 'plant' AND f.target_id = ?", sourceID).
		Distinct("f.user_id").Count(&n).Error
	return n, err
}

// CountGardens counts garden entries pointing at a card.
func (r *PlantMergeRepository) CountGardens(plantID uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.UserGarden{}).Where("plant_species_id = ?", plantID).Count(&n).Error
	return n, err
}

// CountGardenOverlaps counts users growing both cards.
func (r *PlantMergeRepository) CountGardenOverlaps(keepID, sourceID uint) (int64, error) {
	var n int64
	err := r.db.Table("user_gardens g").
		Joins("JOIN user_gardens k ON k.user_id = g.user_id AND k.plant_species_id = ?", keepID).
		Where("g.plant_species_id = ?", sourceID).
		Distinct("g.user_id").Count(&n).Error
	return n, err
}

// CountReminders counts reminders pointing at a card.
func (r *PlantMergeRepository) CountReminders(plantID uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.CareReminder{}).Where("plant_species_id = ?", plantID).Count(&n).Error
	return n, err
}

// CountReminderOverlaps counts source reminders that have a same-user,
// same-title counterpart on the keep card and will be folded by latest date.
func (r *PlantMergeRepository) CountReminderOverlaps(keepID, sourceID uint) (int64, error) {
	var n int64
	err := r.db.Table("care_reminders r").
		Joins("JOIN care_reminders k ON k.user_id = r.user_id AND k.task_title = r.task_title AND k.plant_species_id = ?", keepID).
		Where("r.plant_species_id = ?", sourceID).
		Count(&n).Error
	return n, err
}

// CountPests counts disease/pest entries pointing at a card.
func (r *PlantMergeRepository) CountPests(plantID uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.DiseasePest{}).Where("plant_species_id = ?", plantID).Count(&n).Error
	return n, err
}
