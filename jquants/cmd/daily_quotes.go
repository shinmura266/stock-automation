package cmd

import (
	"fmt"
	"stock-automation/jquants/service"

	"github.com/spf13/cobra"
)

var (
	dailyQuotesCode     string
	dailyQuotesDate     string
	dailyQuotesCount    int
	dailyQuotesInterval int
)

var DailyQuotesCmd = &cobra.Command{
	Use:   "daily_quotes",
	Short: "日次株価四本値取得",
	Long:  "日次株価四本値を取得して、DBへ保存する機能を提供します",
	RunE:  updateDailyQuotes,
}

func init() {
	// フラグを追加
	DailyQuotesCmd.Flags().StringVar(&dailyQuotesCode, "code", "", "銘柄コード（指定しない場合は全銘柄）")
	DailyQuotesCmd.Flags().StringVar(&dailyQuotesDate, "date", "", "日付（YYYY-MM-DD形式、codeともに指定しない場合は当日）")
	DailyQuotesCmd.Flags().IntVar(&dailyQuotesCount, "count", 1, "取得する日数（指定した日付からさかのぼる日数、デフォルト: 1）")
	DailyQuotesCmd.Flags().IntVar(&dailyQuotesInterval, "interval", 5, "インターバル（秒、デフォルト: 5）")
}

func updateDailyQuotes(cmd *cobra.Command, args []string) error {

	service, err := service.NewDailyQuotesService(dailyQuotesInterval)
	if err != nil {
		return fmt.Errorf("株価サービス初期化エラー: %v", err)
	}

	err = service.UpdateDailyQuotesWithCount(dailyQuotesCode, dailyQuotesDate, dailyQuotesCount)
	if err != nil {
		return fmt.Errorf("株価データ更新エラー: %v", err)
	}

	return nil
}
