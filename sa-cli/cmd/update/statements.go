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

var (
	statementsCode     string
	statementsDate     string
	statementsCount    int
	statementsInterval int
	statementsDebug    bool
)

// ExecuteStatementsForDate 指定日付の財務情報取得の実行関数（dailyコマンド用）
// jquantsService, dbService, idTokenは必須パラメータです
func ExecuteStatementsForDate(jquantsService *services.JQuantsService, dbService *services.DatabaseService, idToken, date string, debug bool) error {
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

	// データベースリポジトリを初期化
	statementsRepo := database.NewStatementsRepository(dbService.GetConnection())

	// 財務情報を取得（pagination_key対応で全データ取得）
	resp, err := jquantsService.GetClient().GetStatementsClient().GetAllStatements(idToken, "", date, debug)
	if err != nil {
		return errors.DataRetrieveError(err)
	}

	if len(resp.Statements) == 0 {
		return nil // 0件の場合は正常終了
	}

	// データベースに保存
	if _, err := statementsRepo.SaveStatements(resp); err != nil {
		return errors.DataSaveError(err)
	}

	return nil
}

var StatementsCmd = &cobra.Command{
	Use:   "statements",
	Short: "財務情報を取得して保存",
	RunE: func(cmd *cobra.Command, args []string) error {
		// パラメータのチェック
		if statementsCode == "" && statementsDate == "" {
			return fmt.Errorf("codeまたはdateパラメータのいずれかを指定してください")
		}

		// codeとdateの両方が指定され、かつcountが指定された場合はエラー
		if statementsCode != "" && statementsDate != "" && statementsCount > 1 {
			return fmt.Errorf("codeとdateの両方を指定してcountを1より大きく設定することはできません")
		}

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

		// データベースリポジトリを初期化
		statementsRepo := database.NewStatementsRepository(dbService.GetConnection())
		listedInfoRepo := database.NewListedInfoRepository(dbService.GetConnection())

		// IDトークンを取得
		idToken := jquantsService.GetIDToken()

		// 実行パターンに応じて処理を分岐
		if statementsCode != "" && statementsDate == "" {
			// codeが指定された場合：listed_info昇順で複数銘柄取得
			return executeForCodes(jquantsService, idToken, statementsRepo, listedInfoRepo, statementsCode, statementsCount, statementsInterval, statementsDebug)
		} else if statementsDate != "" && statementsCode == "" {
			// dateが指定された場合：日付を遡って取得
			return executeForDates(jquantsService, idToken, statementsRepo, statementsDate, statementsCount, statementsInterval, statementsDebug)
		} else {
			// codeとdateの両方が指定された場合（countは1のみ）
			return executeSingle(jquantsService, idToken, statementsRepo, statementsCode, statementsDate, statementsDebug)
		}
	},
}

// executeForCodes codeが指定された場合の処理（複数銘柄を昇順で取得）
func executeForCodes(jquantsService *services.JQuantsService, idToken string, statementsRepo *database.StatementsRepository, listedInfoRepo *database.ListedInfoRepository, startCode string, count int, interval int, debug bool) error {
	// 市場コード0109を除外して銘柄コードリストを取得
	codes, err := listedInfoRepo.GetListedCodesExcludingMarket("0109", startCode, count)
	if err != nil {
		return errors.DatabaseError(err)
	}

	if len(codes) == 0 {
		return fmt.Errorf("指定された条件で取得可能な銘柄が見つかりません")
	}

	fmt.Printf("取得対象銘柄: %d件 (開始銘柄: %s)\n", len(codes), startCode)

	for i, code := range codes {
		if debug {
			fmt.Printf("[%d/%d] 財務情報を取得中... (code: %s)\n", i+1, len(codes), code)
		}

		resp, err := jquantsService.GetClient().GetStatementsClient().GetAllStatements(idToken, code, "", debug)
		if err != nil {
			log.Printf("銘柄 %s の取得でエラー: %v", code, err)
		} else {
			if debug {
				fmt.Printf("取得完了: %d件の財務情報データ (code: %s)\n", len(resp.Statements), code)
			}

			if len(resp.Statements) == 0 {
				log.Printf("[%d/%d] 銘柄 %s: 0件処理", i+1, len(codes), code)
			} else {
				if count, err := statementsRepo.SaveStatements(resp); err != nil {
					log.Printf("[%d/%d] 銘柄 %s の保存でエラー: %v", i+1, len(codes), code, err)
				} else {
					log.Printf("[%d/%d] 銘柄 %s: %d件処理", i+1, len(codes), code, count)
				}
			}
		}

		// 最後の銘柄でなければスリープ
		if i < len(codes)-1 {
			if debug {
				fmt.Printf("%d秒待機中...\n", interval)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	log.Printf("statements（複数銘柄）の更新が完了しました: %d件", len(codes))
	return nil
}

// executeForDates dateが指定された場合の処理（日付を遡って取得）
func executeForDates(jquantsService *services.JQuantsService, idToken string, statementsRepo *database.StatementsRepository, startDate string, count int, interval int, debug bool) error {
	// 日付をパース
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		// YYYYMMDD形式もサポート
		if date, err = time.Parse("20060102", startDate); err != nil {
			return fmt.Errorf("無効な日付形式です: %s", startDate)
		}
	}

	fmt.Printf("取得対象期間: %d日分 (開始日: %s)\n", count, date.Format("2006-01-02"))

	for i := 0; i < count; i++ {
		currentDate := date.AddDate(0, 0, -i) // i日前
		dateStr := currentDate.Format("2006-01-02")

		if debug {
			fmt.Printf("[%d/%d] 財務情報を取得中... (date: %s)\n", i+1, count, dateStr)
		}

		resp, err := jquantsService.GetClient().GetStatementsClient().GetAllStatements(idToken, "", dateStr, debug)
		if err != nil {
			log.Printf("日付 %s の取得でエラー: %v", dateStr, err)
		} else {
			if debug {
				fmt.Printf("取得完了: %d件の財務情報データ (date: %s)\n", len(resp.Statements), dateStr)
			}

			if len(resp.Statements) == 0 {
				log.Printf("[%d/%d] 日付 %s: 0件処理", i+1, count, dateStr)
			} else {
				if processedCount, err := statementsRepo.SaveStatements(resp); err != nil {
					log.Printf("[%d/%d] 日付 %s の保存でエラー: %v", i+1, count, dateStr, err)
				} else {
					log.Printf("[%d/%d] 日付 %s: %d件処理", i+1, count, dateStr, processedCount)
				}
			}
		}

		// 最後の日付でなければスリープ
		if i < count-1 {
			if debug {
				fmt.Printf("%d秒待機中...\n", interval)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	log.Printf("statements（複数日付）の更新が完了しました: %d日分", count)
	return nil
}

// executeSingle codeとdateの両方が指定された場合の処理（1回のみ）
func executeSingle(jquantsService *services.JQuantsService, idToken string, statementsRepo *database.StatementsRepository, code string, date string, debug bool) error {
	fmt.Printf("財務情報を取得中... (code: %s, date: %s)\n", code, date)

	resp, err := jquantsService.GetClient().GetStatementsClient().GetAllStatements(idToken, code, date, debug)
	if err != nil {
		return errors.DataRetrieveError(err)
	}

	if debug {
		fmt.Printf("取得完了: %d件の財務情報データ\n", len(resp.Statements))
	}

	if count, err := statementsRepo.SaveStatements(resp); err != nil {
		return errors.DataSaveError(err)
	} else {
		log.Printf("単体実行: %d件処理", count)
	}

	log.Println("statementsの更新が完了しました")
	return nil
}

func init() {
	StatementsCmd.Flags().StringVar(&statementsCode, "code", "", "銘柄コード(4桁/5桁)")
	StatementsCmd.Flags().StringVar(&statementsDate, "date", "", "開示日 (YYYY-MM-DD または YYYYMMDD)")
	StatementsCmd.Flags().IntVar(&statementsCount, "count", 1, "取得回数")
	StatementsCmd.Flags().IntVar(&statementsInterval, "interval", 5, "取得間隔（秒）")
	StatementsCmd.Flags().BoolVar(&statementsDebug, "debug", false, "デバッグモード（詳細ログ出力）")
}
