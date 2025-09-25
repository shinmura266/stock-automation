package helper

import "time"

// GetTodayDate 当日の日付をYYYY-MM-DD形式で取得
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// SubDate 指定した日付から指定日数分さかのぼった日付を計算
func SubDate(date string, days int) string {
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		// パースエラーの場合は当日を使用
		start = time.Now()
	}

	return start.AddDate(0, 0, -days).Format("2006-01-02")
}
