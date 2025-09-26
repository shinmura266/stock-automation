package cmd

import (
	"fmt"
	"stock-automation/jquants/service"

	"github.com/spf13/cobra"
)

var (
	statementsCode     string
	statementsDate     string
	statementsCount    int
	statementsInterval int
)

var StatementsCmd = &cobra.Command{
	Use:   "statements",
	Short: "財務情報取得",
	Long:  "J-Quantsの財務情報を取得して、DBへ保存する機能を提供します",
	RunE:  updateStatements,
}

func init() {
	// フラグを追加
	StatementsCmd.Flags().StringVar(&statementsCode, "code", "", "銘柄コード（指定しない場合は全銘柄）")
	StatementsCmd.Flags().StringVar(&statementsDate, "date", "", "日付（YYYY-MM-DD形式、codeともに指定しない場合は当日）")
	StatementsCmd.Flags().IntVar(&statementsCount, "count", 1, "取得する日数（指定した日付からさかのぼる日数、デフォルト: 1）")
	StatementsCmd.Flags().IntVar(&statementsInterval, "interval", 5, "インターバル（秒、デフォルト: 5）")
}

func updateStatements(cmd *cobra.Command, args []string) error {

	service, err := service.NewStatementsService(statementsInterval)
	if err != nil {
		return fmt.Errorf("財務情報サービス初期化エラー: %v", err)
	}
	defer service.Close()

	err = service.UpdateStatementsWithCount(statementsCode, statementsDate, statementsCount)
	if err != nil {
		return fmt.Errorf("財務情報データ更新エラー: %v", err)
	}

	return nil
}
