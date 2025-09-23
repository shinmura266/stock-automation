package cmd

import (
	"fmt"
	"stock-automation/jquants/service"
	"time"

	"github.com/spf13/cobra"
)

var (
	code string
	date string
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
}

func updateDailyQuotes(cmd *cobra.Command, args []string) error {

	if code == "" && date == "" {
		date = getTodayDate()
	}

	service, err := service.NewDailyQuotesService()
	if err != nil {
		return fmt.Errorf("株価サービス初期化エラー: %v", err)
	}

	// if code != "" {
	// 	quotes, err := service.NewDailyQuotesService().GetDailyQuotesMultipleStocks()
	// 	if err != nil {
	// 		return fmt.Errorf("株価データ取得エラー: %v", err)
	// 	}
	// 	return nil
	// }

	err = service.UpdateDailyQuotesByDate(date)
	if err != nil {
		return fmt.Errorf("株価データ更新エラー: %v", err)
	}
	fmt.Println("株価データ更新完了")

	return nil
}

// getTodayDate 当日の日付をYYYYMMDD形式で取得
func getTodayDate() string {
	return time.Now().Format("2006-01-02")
}
