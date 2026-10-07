package util

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	d := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	if got := FormatDate(d); got != "2026-08-16" {
		t.Errorf("FormatDate = %s", got)
	}
}

func TestPlantTypeText(t *testing.T) {
	cases := []struct{ in, want string }{
		{"flower", "观花"},
		{"foliage", "观叶"},
		{"succulent", "多肉"},
		{"aquatic", "水生"},
		{"unknown", "未知"},
	}
	for _, c := range cases {
		if got := PlantTypeText(c.in); got != c.want {
			t.Errorf("PlantTypeText(%s) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestReminderStatusText(t *testing.T) {
	if got := ReminderStatusText("overdue"); got != "已逾期" {
		t.Errorf("ReminderStatusText = %s", got)
	}
}
