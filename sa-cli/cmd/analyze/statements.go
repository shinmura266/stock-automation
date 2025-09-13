package analyze

import (
	"fmt"

	"kabu-analysis/config"
	"kabu-analysis/database"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/spf13/cobra"
)

var AnalyzeStatementsCmd = &cobra.Command{
	Use:   "statements",
	Short: "財務情報を分析してサマリーを作成",
	Long: `
財務情報から各銘柄の会計年度別最新データを分析し、サマリーテーブルを作成します。

処理内容:
- disclosed_dateの新しいデータで古いデータを更新
- current_fiscal_year_start_date毎にfiscal_yearとして財務情報を整理
- next_fiscal_year_start_dateのデータも翌年度データとして作成
- 実績データと予想データを分離して管理

分析対象項目:
- net_sales (売上高)
- operating_profit (営業利益)
- ordinary_profit (経常利益)
- profit (当期純利益)
- eps (1株当たり当期純利益)
- total_assets (総資産)
- equity (自己資本)
- equity_to_asset_ratio (自己資本比率)
- dividend_per_share (1株当たり配当)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("財務情報の分析処理を開始します...")

		// 設定を読み込み
		cfg, err := config.LoadFromEnv()
		if err != nil {
			return errors.ConfigError(err)
		}

		fmt.Println("データベースに接続中...")
		// データベースサービスを初期化
		dbService, err := services.NewDatabaseService(cfg.Database)
		if err != nil {
			return errors.DatabaseError(err)
		}
		defer dbService.Close()

		// リポジトリを初期化
		summaryRepo := database.NewStatementsSummaryRepository(dbService.GetConnection())

		fmt.Println("財務データの分析・サマリー作成中...")
		if err := summaryRepo.AnalyzeAndCreateSummary(); err != nil {
			return errors.DatabaseError(err)
		}

		fmt.Println("財務情報の分析処理が正常に完了しました")
		return nil
	},
}
