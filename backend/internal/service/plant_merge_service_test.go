package service

import (
	"reflect"
	"testing"
)

func TestNormalizeSourceIDs(t *testing.T) {
	cases := []struct {
		name       string
		keepID     uint
		sourceIDs  []uint
		wantIDs    []uint
		wantIssues int
	}{
		{"sorted and deduped", 1, []uint{5, 3, 3, 2}, []uint{2, 3, 5}, 0},
		{"keep card rejected", 1, []uint{1, 2}, []uint{2}, 1},
		{"zero id rejected", 1, []uint{0, 2}, []uint{2}, 1},
		{"empty sources", 1, nil, []uint{}, 1},
		{"all filtered out", 1, []uint{1, 0}, []uint{}, 3},
		{"single source", 9, []uint{4}, []uint{4}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ids, issues := normalizeSourceIDs(tc.keepID, tc.sourceIDs)
			if !reflect.DeepEqual(ids, tc.wantIDs) {
				t.Errorf("ids = %v, want %v", ids, tc.wantIDs)
			}
			if len(issues) != tc.wantIssues {
				t.Errorf("issues = %v, want %d entries", issues, tc.wantIssues)
			}
		})
	}
}
