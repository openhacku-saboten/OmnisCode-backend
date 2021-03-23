package service

import "time"

const format = "2006-01-02T15:04:05+09:00"

// ConvertStrToTime はstring型のtimestampをtime.Time型に変換します
func ConvertStrToTime(timeStr string) (t time.Time, err error) {
	t, err = time.Parse(format, timeStr)
	return
}

// ConvertTimeToStr はtime.Time型のtimestampをstring型に変換します
func ConvertTimeToStr(t time.Time) string {
	return t.Format(format)
}
