package util

import "testing"

func TestSeasonTask(t *testing.T) {
	cases := []struct {
		month int
		want  string
	}{
		{1, "冬季防寒：注意保温控水，避免冻伤"},
		{4, "春季焕新：换盆施肥，病虫害早防"},
		{7, "夏季遮阴：增加浇水频率，防晒通风"},
		{10, "秋季养护：减少施肥，做好越冬准备"},
	}
	for _, c := range cases {
		if got := SeasonTask(c.month); got != c.want {
			t.Errorf("SeasonTask(%d) = %q, want %q", c.month, got, c.want)
		}
	}
}
