package cmd

import (
	"fmt"
	"stock-automation/jquants/service"
	"time"

	"github.com/spf13/cobra"
)

var (
	statementsCode  string
	statementsDate  string
	statementsCount int
)

var StatementsCmd = &cobra.Command{
	Use:   "statements",
	Short: "財務情報取得",
	Long:  "J-Quantsの財務情報を取得して、DBへ保存する機能を提供します",
	RunE:  updateStatements,
}

func init() {
	// フラグを追加
	StatementsCmd.Flags().StringVar(&statementsCode, "code", "", "銘柄コード（4桁/5桁）")
	StatementsCmd.Flags().StringVar(&statementsDate, "date", "", "開示日（YYYY-MM-DD形式）")
	StatementsCmd.Flags().IntVar(&statementsCount, "count", 1, "取得回数（指定した日付からさかのぼる日数、デフォルト: 1）")
}

func updateStatements(cmd *cobra.Command, args []string) error {
	// パラメータのチェック
	if statementsCode == "" && statementsDate == "" {
		return fmt.Errorf("codeまたはdateパラメータのいずれかを指定してください")
	}

	// codeとdateの両方が指定され、かつcountが指定された場合はエラー
	if statementsCode != "" && statementsDate != "" && statementsCount > 1 {
		return fmt.Errorf("codeとdateの両方を指定してcountを1より大きく設定することはできません")
	}

	service, err := service.NewStatementsService()
	if err != nil {
		return fmt.Errorf("財務情報サービス初期化エラー: %v", err)
	}
	defer service.Close()

	// 実行パターンに応じて処理を分岐
	if statementsCode != "" && statementsDate == "" {
		// codeが指定された場合：指定銘柄の複数日付データを取得
		dates := generateBackwardDatesForStatements(getTodayDateForStatements(), statementsCount)
		err = service.UpdateStatementsMultipleDates([]string{statementsCode}, dates, 1000, false)
		if err != nil {
			return fmt.Errorf("財務情報データ更新エラー: %v", err)
		}
		fmt.Printf("財務情報データ更新完了（銘柄: %s, %d日分）\n", statementsCode, len(dates))
	} else if statementsDate != "" && statementsCode == "" {
		// dateが指定された場合：日付を遡って取得
		dates := generateBackwardDatesForStatements(statementsDate, statementsCount)
		err = service.UpdateStatementsMultipleDates([]string{}, dates, 1000, false)
		if err != nil {
			return fmt.Errorf("財務情報データ更新エラー: %v", err)
		}
		fmt.Printf("財務情報データ更新完了（%d日分）\n", len(dates))
	} else {
		// codeとdateの両方が指定された場合（countは1のみ）
		err = service.UpdateStatementsByCodeAndDate(statementsCode, statementsDate, false)
		if err != nil {
			return fmt.Errorf("財務情報データ更新エラー: %v", err)
		}
		fmt.Printf("財務情報データ更新完了（銘柄: %s, 日付: %s）\n", statementsCode, statementsDate)
	}

	return nil
}

// getTodayDateForStatements 当日の日付をYYYY-MM-DD形式で取得
func getTodayDateForStatements() string {
	return time.Now().Format("2006-01-02")
}

// generateBackwardDatesForStatements 指定した日付からcount数分の日付をさかのぼって生成
func generateBackwardDatesForStatements(startDate string, count int) []string {
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
