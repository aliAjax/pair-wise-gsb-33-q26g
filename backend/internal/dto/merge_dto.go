package dto

import "time"

// PlantMergeRequest selects the keep card and the source cards merged into it.
type PlantMergeRequest struct {
	KeepID    uint   `json:"keep_id" binding:"required"`
	SourceIDs []uint `json:"source_ids" binding:"required,min=1,max=20"`
}

// MergeStats counts migrated and deduplicated associations of one merge task.
type MergeStats struct {
	FavoritesMoved   int `json:"favorites_moved"`
	FavoritesDeduped int `json:"favorites_deduped"`
	GardensMoved     int `json:"gardens_moved"`
	GardensDeduped   int `json:"gardens_deduped"`
	RemindersMoved   int `json:"reminders_moved"`
	RemindersMerged  int `json:"reminders_merged"`
	PestsMoved       int `json:"pests_moved"`
}

// MergePrecheckItem reports the migration impact for one source card.
type MergePrecheckItem struct {
	SourceID       uint   `json:"source_id"`
	SourceName     string `json:"source_name"`
	FavoritesTotal int64  `json:"favorites_total"`
	FavoritesDedup int64  `json:"favorites_dedup"`
	GardensTotal   int64  `json:"gardens_total"`
	GardensDedup   int64  `json:"gardens_dedup"`
	RemindersTotal int64  `json:"reminders_total"`
	RemindersMerge int64  `json:"reminders_merge"`
	PestsTotal     int64  `json:"pests_total"`
	AlreadyMerged  bool   `json:"already_merged"`
}

// MergePrecheckResult is the precheck report for a merge batch.
type MergePrecheckResult struct {
	KeepID   uint                `json:"keep_id"`
	KeepName string              `json:"keep_name"`
	Items    []MergePrecheckItem `json:"items"`
}

// PlantMergeResult is the API view of one merge task.
type PlantMergeResult struct {
	ID         uint        `json:"id"`
	KeepID     uint        `json:"keep_id"`
	SourceID   uint        `json:"source_id"`
	KeepName   string      `json:"keep_name"`
	SourceName string      `json:"source_name"`
	Status     string      `json:"status"`
	Stage      string      `json:"stage"`
	Stats      *MergeStats `json:"stats"`
	LastError  string      `json:"last_error"`
	OperatorID uint        `json:"operator_id"`
	CreatedAt  time.Time   `json:"created_at"`
	FinishedAt *time.Time  `json:"finished_at"`
}
