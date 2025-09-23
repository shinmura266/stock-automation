package helper

import "time"

// GetTodayDate 当日の日付をYYYY-MM-DD形式で取得
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}
