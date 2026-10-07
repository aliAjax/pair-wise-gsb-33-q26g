package model

import (
	"reflect"
	"testing"
)

func TestBuildMergeTaskKeyOrderInsensitive(t *testing.T) {
	k1 := BuildMergeTaskKey(1, []uint{3, 2, 5})
	k2 := BuildMergeTaskKey(1, []uint{5, 3, 2})
	if k1 != k2 {
		t.Errorf("same batch should map to same key: %q vs %q", k1, k2)
	}
}

func TestBuildMergeTaskKeyDistinct(t *testing.T) {
	cases := []struct {
		name string
		a, b string
	}{
		{"different keep", BuildMergeTaskKey(1, []uint{2, 3}), BuildMergeTaskKey(2, []uint{1, 3})},
		{"different sources", BuildMergeTaskKey(1, []uint{2, 3}), BuildMergeTaskKey(1, []uint{2, 4})},
		{"subset sources", BuildMergeTaskKey(1, []uint{2}), BuildMergeTaskKey(1, []uint{2, 3})},
	}
	for _, tc := range cases {
		if tc.a == tc.b {
			t.Errorf("%s: keys should differ, both %q", tc.name, tc.a)
		}
	}
}

func TestMergeCheckpointMarkAndDone(t *testing.T) {
	cp := ParseCheckpoint("")
	for _, phase := range MergePhases() {
		if cp.Done(phase, 7) {
			t.Errorf("phase %s should not be done initially", phase)
		}
	}
	cp.Mark(MergePhaseFavorites, 7)
	cp.Mark(MergePhaseFavorites, 7) // duplicate mark is a no-op
	if !cp.Done(MergePhaseFavorites, 7) {
		t.Error("favorites unit 7 should be done")
	}
	if len(cp.Favorites) != 1 {
		t.Errorf("duplicate mark appended: %+v", cp.Favorites)
	}
	if cp.Done(MergePhaseGardens, 7) {
		t.Error("gardens unit 7 should still be pending")
	}
	if cp.Done("unknown-phase", 7) {
		t.Error("unknown phase should never be done")
	}
}

func TestMergeCheckpointJSONRoundTrip(t *testing.T) {
	cp := ParseCheckpoint("")
	cp.Mark(MergePhaseReminders, 3)
	cp.Mark(MergePhasePests, 3)
	cp.Result.FavoritesMoved = 4
	cp.Result.RemindersMerged = 2

	restored := ParseCheckpoint(cp.JSON())
	if !restored.Done(MergePhaseReminders, 3) || !restored.Done(MergePhasePests, 3) {
		t.Errorf("restored checkpoint lost units: %+v", restored)
	}
	if restored.Result.FavoritesMoved != 4 || restored.Result.RemindersMerged != 2 {
		t.Errorf("restored checkpoint lost result counters: %+v", restored.Result)
	}
	if restored.Done(MergePhaseRetire, 3) {
		t.Error("retire unit should not be done")
	}
}

func TestParseCheckpointInvalid(t *testing.T) {
	cp := ParseCheckpoint("{not-json")
	if cp == nil || len(cp.Favorites) != 0 {
		t.Errorf("invalid json should yield empty checkpoint, got %+v", cp)
	}
}

func TestPlantMergeTaskSourceIDList(t *testing.T) {
	task := &PlantMergeTask{SourceIDs: EncodeSourceIDs([]uint{2, 3, 5})}
	got := task.SourceIDList()
	if !reflect.DeepEqual(got, []uint{2, 3, 5}) {
		t.Errorf("unexpected source ids: %v", got)
	}
	broken := &PlantMergeTask{SourceIDs: "oops"}
	if ids := broken.SourceIDList(); len(ids) != 0 {
		t.Errorf("broken json should yield empty list, got %v", ids)
	}
}

func TestPlantMergeTaskResultData(t *testing.T) {
	empty := &PlantMergeTask{}
	if empty.ResultData() != nil {
		t.Error("empty result should decode to nil")
	}
	task := &PlantMergeTask{Result: `{"favorites_moved":2,"retired_plants":1}`}
	r := task.ResultData()
	if r == nil || r.FavoritesMoved != 2 || r.RetiredPlants != 1 {
		t.Errorf("unexpected result: %+v", r)
	}
}
