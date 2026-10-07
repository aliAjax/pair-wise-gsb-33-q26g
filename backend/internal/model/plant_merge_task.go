package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Plant merge task statuses.
const (
	MergeTaskPending    = "pending"
	MergeTaskProcessing = "processing"
	MergeTaskSucceeded  = "succeeded"
	MergeTaskFailed     = "failed"
)

// Plant merge phases, executed in order for every source card.
const (
	MergePhaseFavorites = "favorites"
	MergePhaseGardens   = "gardens"
	MergePhaseReminders = "reminders"
	MergePhasePests     = "pests"
	MergePhaseRetire    = "retired"
)

// MergePhases returns the ordered phase list of a merge task.
func MergePhases() []string {
	return []string{MergePhaseFavorites, MergePhaseGardens, MergePhaseReminders, MergePhasePests, MergePhaseRetire}
}

// MergeResult accumulates the migration counters of a merge task.
type MergeResult struct {
	FavoritesMoved   int64 `json:"favorites_moved"`
	FavoritesDeduped int64 `json:"favorites_deduped"`
	GardensMoved     int64 `json:"gardens_moved"`
	GardensDeduped   int64 `json:"gardens_deduped"`
	RemindersMoved   int64 `json:"reminders_moved"`
	RemindersMerged  int64 `json:"reminders_merged"`
	PestsMoved       int64 `json:"pests_moved"`
	RetiredPlants    int64 `json:"retired_plants"`
}

// MergePrecheckStats estimates how many rows each association type will touch.
type MergePrecheckStats struct {
	FavoritesTotal  int64 `json:"favorites_total"`
	FavoritesDedupe int64 `json:"favorites_dedupe"`
	GardensTotal    int64 `json:"gardens_total"`
	GardensDedupe   int64 `json:"gardens_dedupe"`
	RemindersTotal  int64 `json:"reminders_total"`
	RemindersMerge  int64 `json:"reminders_merge"`
	PestsTotal      int64 `json:"pests_total"`
}

// MergePrecheck is the pre-migration snapshot persisted with the task.
type MergePrecheck struct {
	KeepID      uint               `json:"keep_id"`
	KeepName    string             `json:"keep_name"`
	SourceIDs   []uint             `json:"source_ids"`
	SourceNames []string           `json:"source_names"`
	Issues      []string           `json:"issues"`
	Stats       MergePrecheckStats `json:"stats"`
}

// MergeCheckpoint records which (phase, source) units are already done so a
// retry after failure only patches unfinished objects instead of repeating work.
type MergeCheckpoint struct {
	Favorites []uint      `json:"favorites"`
	Gardens   []uint      `json:"gardens"`
	Reminders []uint      `json:"reminders"`
	Pests     []uint      `json:"pests"`
	Retired   []uint      `json:"retired"`
	Result    MergeResult `json:"result"`
}

// phaseList returns the source-id list backing the given phase.
func (c *MergeCheckpoint) phaseList(phase string) *[]uint {
	switch phase {
	case MergePhaseFavorites:
		return &c.Favorites
	case MergePhaseGardens:
		return &c.Gardens
	case MergePhaseReminders:
		return &c.Reminders
	case MergePhasePests:
		return &c.Pests
	case MergePhaseRetire:
		return &c.Retired
	}
	return nil
}

// Done reports whether the (phase, source) unit is already checkpointed.
func (c *MergeCheckpoint) Done(phase string, sourceID uint) bool {
	list := c.phaseList(phase)
	if list == nil {
		return false
	}
	for _, id := range *list {
		if id == sourceID {
			return true
		}
	}
	return false
}

// Mark records the (phase, source) unit as completed.
func (c *MergeCheckpoint) Mark(phase string, sourceID uint) {
	list := c.phaseList(phase)
	if list == nil || c.Done(phase, sourceID) {
		return
	}
	*list = append(*list, sourceID)
}

// JSON serializes the checkpoint for persistence.
func (c *MergeCheckpoint) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ParseCheckpoint decodes a checkpoint JSON string; empty input yields a fresh checkpoint.
func ParseCheckpoint(raw string) *MergeCheckpoint {
	cp := &MergeCheckpoint{}
	if raw == "" {
		return cp
	}
	if err := json.Unmarshal([]byte(raw), cp); err != nil {
		return &MergeCheckpoint{}
	}
	return cp
}

// PlantMergeTask is a plant species merge job with checkpoint recovery.
type PlantMergeTask struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskKey    string    `gorm:"size:128;uniqueIndex;not null" json:"task_key"`
	KeepID     uint      `gorm:"index;not null" json:"keep_id"`
	SourceIDs  string    `gorm:"type:json;not null" json:"source_ids"`
	Status     string    `gorm:"size:16;index;not null;default:pending" json:"status"`
	Checkpoint string    `gorm:"type:json" json:"checkpoint"`
	Precheck   string    `gorm:"type:json" json:"precheck"`
	Result     string    `gorm:"type:json" json:"result"`
	LastError  string    `gorm:"size:512" json:"last_error"`
	OperatorID uint      `json:"operator_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CheckpointData decodes the persisted checkpoint.
func (t *PlantMergeTask) CheckpointData() *MergeCheckpoint {
	return ParseCheckpoint(t.Checkpoint)
}

// SourceIDList decodes the source card id list.
func (t *PlantMergeTask) SourceIDList() []uint {
	ids := []uint{}
	if err := json.Unmarshal([]byte(t.SourceIDs), &ids); err != nil {
		return []uint{}
	}
	return ids
}

// ResultData decodes the persisted migration result; it is nil until the
// task finishes and stores real counters.
func (t *PlantMergeTask) ResultData() *MergeResult {
	if t.Result == "" || t.Result == "{}" {
		return nil
	}
	r := &MergeResult{}
	if err := json.Unmarshal([]byte(t.Result), r); err != nil {
		return nil
	}
	return r
}

// EncodeSourceIDs serializes source ids for the SourceIDs column.
func EncodeSourceIDs(ids []uint) string {
	b, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// BuildMergeTaskKey builds the idempotency key of a merge batch: the same
// keep card plus the same source set always maps to the same key, so a
// duplicate submission hits the unique index and returns the existing task.
func BuildMergeTaskKey(keepID uint, sourceIDs []uint) string {
	ids := append([]uint(nil), sourceIDs...)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatUint(uint64(id), 10)
	}
	return fmt.Sprintf("K%d-S%s", keepID, strings.Join(parts, "."))
}
