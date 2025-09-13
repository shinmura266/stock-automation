package update

import (
	"fmt"
	"log"
	"time"

	"kabu-analysis/config"
	"kabu-analysis/database"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/spf13/cobra"
)

// ExecuteListedInfo 上場銘柄一覧取得の実行関数
// jquantsService, dbService, idTokenは必須パラメータです
func ExecuteListedInfo(jquantsService *services.JQuantsService, dbService *services.DatabaseService, idToken, date string) error {
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

	// 上場銘柄一覧を取得
	targetDate := date
	if targetDate == "" {
		targetDate = time.Now().Format("2006-01-02")
	}

	// 上場銘柄一覧を取得（pagination_key対応で全データ取得）
	listedInfo, err := jquantsService.GetClient().GetListedClient().GetAllListedInfo(idToken, "", targetDate)
	if err != nil {
		return errors.DataRetrieveError(err)
	}

	// データベースに保存
	repository := database.NewListedInfoRepository(dbService.GetConnection())
	if err := repository.SaveListedInfo(listedInfo); err != nil {
		return errors.DataSaveError(err)
	}

	return nil
}

var dateFlag string

var ListedInfoCmd = &cobra.Command{
	Use:   "listed_info",
	Short: "上場銘柄一覧を取得して保存",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("listed_infoコマンドを開始します")

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

		fmt.Println("J-Quants APIに認証中...")
		// J-Quantsサービスを初期化
		jquantsService := services.NewJQuantsService(cfg.JQuants)
		if err := jquantsService.Authenticate(); err != nil {
			return errors.JQuantsAuthError(err)
		}

		// IDトークンを取得
		idToken := jquantsService.GetIDToken()

		if err := ExecuteListedInfo(jquantsService, dbService, idToken, dateFlag); err != nil {
			return err
		}

		log.Println("listed_infoコマンドが完了しました")
		return nil
	},
}

func init() {
	ListedInfoCmd.Flags().StringVar(&dateFlag, "date", "", "対象日付 (YYYY-MM-DD形式、省略時は今日)")
}
