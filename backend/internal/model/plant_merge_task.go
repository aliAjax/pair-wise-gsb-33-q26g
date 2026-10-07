package model

import "time"

// Plant merge task statuses.
const (
	MergePending = "pending"
	MergeRunning = "running"
	MergeDone    = "done"
	MergeFailed  = "failed"
)

// Merge stages execute in this exact order. PlantMergeTask.Stage stores the
// last completed stage and acts as the resume checkpoint after a failure.
const (
	MergeStageFavorites = "favorites"
	MergeStageGardens   = "gardens"
	MergeStageReminders = "reminders"
	MergeStagePests     = "pests"
	MergeStageRetire    = "retire"
)

// MergeStageOrder lists every stage in execution order.
var MergeStageOrder = []string{
	MergeStageFavorites,
	MergeStageGardens,
	MergeStageReminders,
	MergeStagePests,
	MergeStageRetire,
}

// PlantMergeTask records one source-card into keep-card variety merge.
// The unique (keep_id, source_id) pair makes concurrent submissions
// idempotent: the first writer executes, later ones observe its result.
type PlantMergeTask struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	KeepID     uint       `gorm:"index:idx_merge_pair,unique;not null" json:"keep_id"`
	SourceID   uint       `gorm:"index:idx_merge_pair,unique;not null" json:"source_id"`
	KeepName   string     `gorm:"size:128" json:"keep_name"`
	SourceName string     `gorm:"size:128" json:"source_name"`
	Status     string     `gorm:"size:16;index;not null;default:pending" json:"status"`
	Stage      string     `gorm:"size:32" json:"stage"`
	Stats      string     `gorm:"type:json" json:"stats"`
	LastError  string     `gorm:"size:512" json:"last_error"`
	OperatorID uint       `json:"operator_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at"`
}
