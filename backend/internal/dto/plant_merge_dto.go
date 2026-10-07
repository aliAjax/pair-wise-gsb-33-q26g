package dto

import (
	"time"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
)

// PlantMergeRequest selects the keep card and the source cards to merge away.
type PlantMergeRequest struct {
	KeepID    uint   `json:"keep_id" binding:"required,min=1"`
	SourceIDs []uint `json:"source_ids" binding:"required,min=1,dive,min=1"`
}

// MergeTaskView is the API representation of a plant merge task.
type MergeTaskView struct {
	ID         uint               `json:"id"`
	TaskKey    string             `json:"task_key"`
	KeepID     uint               `json:"keep_id"`
	SourceIDs  []uint             `json:"source_ids"`
	Status     string             `json:"status"`
	Result     *model.MergeResult `json:"result,omitempty"`
	LastError  string             `json:"last_error"`
	OperatorID uint               `json:"operator_id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// NewMergeTaskView builds a MergeTaskView from the persisted task.
func NewMergeTaskView(t *model.PlantMergeTask) MergeTaskView {
	return MergeTaskView{
		ID:         t.ID,
		TaskKey:    t.TaskKey,
		KeepID:     t.KeepID,
		SourceIDs:  t.SourceIDList(),
		Status:     t.Status,
		Result:     t.ResultData(),
		LastError:  t.LastError,
		OperatorID: t.OperatorID,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
