package gtbox_time

import (
	"time"
)

func GTToolsTimeStringCovertToUTCTime(tiemString string) time.Time {
	utcloc, _ := time.LoadLocation("UTC")
	aTime, _ := time.ParseInLocation("2006-01-02T15:04:05Z", tiemString, utcloc)
	bTime := time.Unix(aTime.Unix(), 0).UTC()
	return bTime
}

func GTToolsTimestampCovertToBeijing(timestamp float64) time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	utcTIme := time.Unix(int64(timestamp), 0)
	beijingTIme := utcTIme.In(beijingLoc)
	return beijingTIme
}

func GTToolsTimesGetBeijingTime() time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	beijingTIme := time.Now().In(beijingLoc)
	return beijingTIme
}

func GTToolsTimeUTCCovertToBeijing(inTime time.Time) time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	beijingTIme := inTime.In(beijingLoc)
	return beijingTIme
}

// GTDateGetNowYearMoonDay 获取当前年月日字符串 时间格式为"2006-01-02"
func GTDateGetNowYearMoonDay() string {
	return time.Now().Format("2006-01-02")
}

// GTDateGetYearMoonDayFromTime 获取当前年月日字符串 时间格式为"2006-01-02"
func GTDateGetYearMoonDayFromTime(aTime time.Time) string {
	return aTime.Format("2006-01-02")
}

func GTDateEqualYearMoonDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GTDateEqualYearMoonDayHours(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	h1 := date1.Hour()
	y2, m2, d2 := date2.Date()
	h2 := date2.Hour()
	return y1 == y2 && m1 == m2 && d1 == d2 && h1 == h2
}

// GTGetTodayCustomHoursAndMinuteWithBeijing 获取北京时间今天指定的小时、分钟
func GTGetTodayCustomHoursAndMinuteWithBeijing(aHours int, aMinute int) time.Time {
	aNow := time.Now()
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海

	return time.Date(aNow.Year(), aNow.Month(), aNow.Day(), aHours, aMinute, 0, 0, beijingLoc)
}
