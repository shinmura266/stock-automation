package cmd

import (
	"fmt"
	"stock-automation/jquants/service"
	"time"

	"github.com/spf13/cobra"
)

var (
	code  string
	date  string
	count int
)

var DailyQuotesCmd = &cobra.Command{
	Use:   "daily-quotes",
	Short: "日次株価四本値取得",
	Long:  "日次株価四本値を取得して、DBへ保存する機能を提供します",
	RunE:  updateDailyQuotes,
}

func init() {
	// フラグを追加
	DailyQuotesCmd.Flags().StringVar(&code, "code", "", "銘柄コード（指定しない場合は全銘柄）")
	DailyQuotesCmd.Flags().StringVar(&date, "date", "", "日付（YYYY-MM-DD形式、codeともに指定しない場合は当日）")
	DailyQuotesCmd.Flags().IntVar(&count, "count", 1, "取得する日数（指定した日付からさかのぼる日数、デフォルト: 1）")
}

func updateDailyQuotes(cmd *cobra.Command, args []string) error {

	if code == "" && date == "" {
		date = getTodayDate()
	}

	service, err := service.NewDailyQuotesService()
	if err != nil {
		return fmt.Errorf("株価サービス初期化エラー: %v", err)
	}

	// count数分の日付リストを生成
	dates := generateBackwardDates(date, count)

	if code != "" {
		// 指定銘柄の複数日付データを取得
		codes := []string{code}
		err = service.UpdateDailyQuotesMultipleDates(codes, dates, 1000)
		if err != nil {
			return fmt.Errorf("株価データ更新エラー: %v", err)
		}
	} else {
		// 全銘柄の複数日付データを取得
		err = service.UpdateDailyQuotesMultipleDates([]string{}, dates, 1000)
		if err != nil {
			return fmt.Errorf("株価データ更新エラー: %v", err)
		}
	}

	fmt.Printf("株価データ更新完了（%d日分）\n", len(dates))

	return nil
}

// getTodayDate 当日の日付をYYYYMMDD形式で取得
func getTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// generateBackwardDates 指定した日付からcount数分の日付をさかのぼって生成
func generateBackwardDates(startDate string, count int) []string {
	if count <= 0 {
		count = 1
	}

	dates := make([]string, count)
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		// パースエラーの場合は当日を使用
		start = time.Now()
	}

	for i := 0; i < count; i++ {
		dates[i] = start.AddDate(0, 0, -i).Format("2006-01-02")
	}

	return dates
}
