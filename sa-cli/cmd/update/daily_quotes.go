package update

import (
	"fmt"

	"kabu-analysis/config"
	"kabu-analysis/database"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/spf13/cobra"
)

var (
	dailyDate string
	dailyCode string
)

// ExecuteDailyQuotes 株価四本値取得の実行関数
// jquantsService, dbService, idTokenは必須パラメータです
func ExecuteDailyQuotes(jquantsService *services.JQuantsService, dbService *services.DatabaseService, idToken, code, date string) error {
	// 必須パラメータのチェック
	if jquantsService == nil {
		return fmt.Errorf("jquantsServiceは必須です")
	}
	if dbService == nil {
		return fmt.Errorf("dbServiceは必須です")
	}
	if idToken == "" {
		return fmt.Errorf("idTokenは必須です")
	}

	// 四本値データを取得（pagination_key対応で全データ取得）
	resp, err := jquantsService.GetClient().GetDailyClient().GetAllDailyQuotes(idToken, code, date)
	if err != nil {
		return errors.DataRetrieveError(err)
	}

	// データベースに保存
	repo := database.NewDailyQuotesRepository(dbService.GetConnection())
	if err := repo.SaveDailyQuotes(resp); err != nil {
		return errors.DataSaveError(err)
	}
	return nil
}

var DailyQuotesCmd = &cobra.Command{
	Use:   "daily_quotes",
	Short: "株価四本値を取得して保存",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 設定を読み込み
		cfg, err := config.LoadFromEnv()
		if err != nil {
			return errors.ConfigError(err)
		}

		// データベースサービスを初期化
		dbService, err := services.NewDatabaseService(cfg.Database)
		if err != nil {
			return errors.DatabaseError(err)
		}
		defer dbService.Close()

		// J-Quantsサービスを初期化
		jquantsService := services.NewJQuantsService(cfg.JQuants)
		if err := jquantsService.Authenticate(); err != nil {
			return errors.JQuantsAuthError(err)
		}

		// IDトークンを取得
		idToken := jquantsService.GetIDToken()

		return ExecuteDailyQuotes(jquantsService, dbService, idToken, dailyCode, dailyDate)
	},
}

func init() {
	DailyQuotesCmd.Flags().StringVar(&dailyDate, "date", "", "基準日 (YYYY-MM-DD または YYYYMMDD)")
	DailyQuotesCmd.Flags().StringVar(&dailyCode, "code", "", "銘柄コード(4桁/5桁)")
}
