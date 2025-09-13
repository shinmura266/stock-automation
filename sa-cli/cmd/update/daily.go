package update

import (
	"fmt"
	"log"
	"time"

	"kabu-analysis/config"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/spf13/cobra"
)

var (
	batchDate     string
	batchCount    int
	batchInterval int
	batchDebug    bool
)

var DailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "日次データ（listed_info、daily_quotes、statements）を一括取得",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 設定を読み込み
		cfg, err := config.LoadFromEnv()
		if err != nil {
			return errors.ConfigError(err)
		}

		fmt.Println("データベースに接続中...")
		// データベースサービスを初期化（一度だけ）
		dbService, err := services.NewDatabaseService(cfg.Database)
		if err != nil {
			return errors.DatabaseError(err)
		}
		defer dbService.Close()

		fmt.Println("J-Quants APIに認証中...")
		// J-Quantsサービスを初期化（一度だけ）
		jquantsService := services.NewJQuantsService(cfg.JQuants)
		if err := jquantsService.Authenticate(); err != nil {
			return errors.JQuantsAuthError(err)
		}

		// IDトークンを取得（一度だけ）
		idToken := jquantsService.GetIDToken()

		// 開始日を決定
		var startDate time.Time
		if batchDate == "" {
			startDate = time.Now()
		} else {
			// 日付をパース
			if date, err := time.Parse("2006-01-02", batchDate); err == nil {
				startDate = date
			} else if date, err := time.Parse("20060102", batchDate); err == nil {
				startDate = date
			} else {
				return fmt.Errorf("無効な日付形式です: %s", batchDate)
			}
		}

		fmt.Printf("日次データ更新開始: %d日分 (開始日: %s)\n", batchCount, startDate.Format("2006-01-02"))

		// 各日付について処理を実行
		for i := 0; i < batchCount; i++ {
			currentDate := startDate.AddDate(0, 0, -i) // i日前
			dateStr := currentDate.Format("2006-01-02")

			fmt.Printf("\n[%d/%d] %s のデータを取得中...\n", i+1, batchCount, dateStr)

			// 1. listed_info を更新（初回のみ）
			if i == 0 {
				fmt.Println("  - listed_info を取得中...")
				if err := ExecuteListedInfo(jquantsService, dbService, idToken, dateStr); err != nil {
					log.Printf("listed_info更新エラー: %v", err)
				} else {
					log.Printf("  ✓ listed_info: 完了")
				}
			}

			// 2. daily_quotes を更新
			fmt.Println("  - daily_quotes を取得中...")
			if err := ExecuteDailyQuotes(jquantsService, dbService, idToken, "", dateStr); err != nil {
				log.Printf("daily_quotes更新エラー (%s): %v", dateStr, err)
			} else {
				log.Printf("  ✓ daily_quotes: 完了")
			}

			// 3. statements を更新
			fmt.Println("  - statements を取得中...")
			if err := ExecuteStatementsForDate(jquantsService, dbService, idToken, dateStr, batchDebug); err != nil {
				log.Printf("statements更新エラー (%s): %v", dateStr, err)
			} else {
				log.Printf("  ✓ statements: 完了")
			}

			// 最後の日付でなければスリープ
			if i < batchCount-1 {
				if batchDebug {
					fmt.Printf("%d秒待機中...\n", batchInterval)
				}
				time.Sleep(time.Duration(batchInterval) * time.Second)
			}
		}

		// 全ての日次データ取得完了後にassessmentを更新
		fmt.Println("\n銘柄評価データ（assessment）を更新中...")
		if err := ExecuteAssessment(dbService); err != nil {
			log.Printf("assessment更新エラー: %v", err)
		} else {
			log.Printf("  ✓ assessment: 完了")
		}

		log.Printf("日次データ更新が完了しました: %d日分", batchCount)
		return nil
	},
}

func init() {
	DailyCmd.Flags().StringVar(&batchDate, "date", "", "基準日 (YYYY-MM-DD または YYYYMMDD、省略時は当日)")
	DailyCmd.Flags().IntVar(&batchCount, "count", 1, "取得する日数（指定日から過去に遡る）")
	DailyCmd.Flags().IntVar(&batchInterval, "interval", 10, "各日付間の取得間隔（秒）")
	DailyCmd.Flags().BoolVar(&batchDebug, "debug", false, "デバッグモード（詳細ログ出力）")
}
