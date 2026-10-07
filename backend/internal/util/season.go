package util

import "time"

// SeasonTask returns the seasonal care highlight for a month (1-12).
func SeasonTask(month int) string {
	switch month {
	case 12, 1, 2:
		return "冬季防寒：注意保温控水，避免冻伤"
	case 3, 4, 5:
		return "春季焕新：换盆施肥，病虫害早防"
	case 6, 7, 8:
		return "夏季遮阴：增加浇水频率，防晒通风"
	default:
		return "秋季养护：减少施肥，做好越冬准备"
	}
}

// CurrentSeasonTask is a convenience wrapper using the current month.
func CurrentSeasonTask() string {
	return SeasonTask(int(time.Now().Month()))
}
