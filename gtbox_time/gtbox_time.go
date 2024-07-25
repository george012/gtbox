/*
Package gtbox_time 时间相关工具
*/
package gtbox_time

import (
	"fmt"
	"time"
)

// GTToolsTimeGetCurrentTimeWithUTC	时间字符串 转 UTC时间
func GTToolsTimeGetCurrentTimeWithUTC() time.Time {
	return time.Now().UTC()
}

// GTToolsTimeStringCovertToUTCTime	时间字符串 转 UTC时间
func GTToolsTimeStringCovertToUTCTime(timeString string) time.Time {
	const layout = "2006-01-02T15:04:05Z"
	at, _ := time.Parse(layout, timeString)
	return at
}

// GTToolsTimestampCovertToBeijing	通过时间戳获取北京时间支持10位、13位时间戳
func GTToolsTimestampCovertToBeijing(timestamp float64) time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海

	timestampLength := len(fmt.Sprintf("%.0f", timestamp))
	var utcTime time.Time

	switch timestampLength {
	case 10:
		utcTime = time.Unix(int64(timestamp), 0).UTC()
	case 13:
		utcTime = time.UnixMilli(int64(timestamp)).UTC()
	case 16:
		// 微秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)*int64(time.Microsecond)).UTC()
	case 19:
		// 纳秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)).UTC()
	}

	beijingTIme := utcTime.In(beijingLoc)
	return beijingTIme
}

// GTToolsTimestampCovertToUTC 通过时间戳获取 UTC 时间支持10位、13位时间戳
func GTToolsTimestampCovertToUTC(timestamp float64) time.Time {
	timestampLength := len(fmt.Sprintf("%.0f", timestamp))
	var utcTime time.Time

	switch timestampLength {
	case 10:
		utcTime = time.Unix(int64(timestamp), 0).UTC()
	case 13:
		utcTime = time.UnixMilli(int64(timestamp)).UTC()
	case 16:
		// 微秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)*int64(time.Microsecond)).UTC()
	case 19:
		// 纳秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)).UTC()
	}

	return utcTime
}

// GTToolsTimesGetBeijingTime	普通时间 转 北京时间
func GTToolsTimesGetBeijingTime() time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	beijingTIme := time.Now().In(beijingLoc)
	return beijingTIme
}

// GTToolsTimeUTCCovertToBeijing	UTC时间 转 北京时间
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

// NowUTC --[☑]--Required
/*
	en: get utc format now time;
	zh-CN: 获取UTC 当前时间;
	@return [☑] en:  ;zh-CN: UTC 时间;
*/
func NowUTC() time.Time {
	return time.Now().UTC()
}
