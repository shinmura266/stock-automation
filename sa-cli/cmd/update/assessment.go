package update

import (
	"fmt"
	"time"

	"kabu-analysis/config"
	"kabu-analysis/database"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/spf13/cobra"
)

// ExecuteAssessment 銘柄評価データ更新の実行関数
// dbServiceは必須パラメータです
func ExecuteAssessment(dbService *services.DatabaseService) error {
	// 必須パラメータのチェック
	if dbService == nil {
		return fmt.Errorf("dbServiceは必須です")
	}

	// リポジトリを初期化
	assessmentRepo := NewAssessmentUpdateRepository(dbService.GetConnection())

	// 古いeffective_dateを持つ銘柄のassessmentレコードを削除
	deletedCount, err := assessmentRepo.DeleteAssessmentByOldEffectiveDate()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 古いレコード削除: %d件\n", deletedCount)

	// 新しい銘柄をassessmentテーブルに追加
	addedCount, err := assessmentRepo.AddNewListedStocks()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 新規銘柄追加: %d件\n", addedCount)

	// statements_summary関連データを更新
	updatedCount1, err := assessmentRepo.UpdateFromStatementsSummary()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 財務データ更新: %d件\n", updatedCount1)

	// daily_quotes関連データを更新
	updatedCount2, err := assessmentRepo.UpdateFromDailyQuotes()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 株価データ更新: %d件\n", updatedCount2)

	// 配当利回りを更新
	updatedCount3, err := assessmentRepo.UpdateDividendYield()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 配当利回り更新: %d件\n", updatedCount3)

	// 3か月価格データと乖離率を更新
	updatedCount4, err := assessmentRepo.UpdateThreeMonthPriceData()
	if err != nil {
		return errors.DatabaseError(err)
	}
	fmt.Printf("  ✓ 3か月価格データ更新: %d件\n", updatedCount4)

	return nil
}

var AssessmentCmd = &cobra.Command{
	Use:   "assessment",
	Short: "銘柄評価データを更新",
	Long: `
assessmentテーブルの銘柄評価データを更新します。

このコマンドは以下の処理を行います：
1. listed_infoテーブルから本日より前のeffective_dateを持つ銘柄コードを抽出し、
   該当するassessmentレコードを削除
2. listed_infoテーブルから本日以降のeffective_dateを持つmarket_code 109以外の
   銘柄をassessmentテーブルに新規追加
3. statements_summaryから最新年度データでassessmentテーブルを更新:
   - 銘柄ごとの最新年度でlast_fiscal_year_end_dateを更新
   - 最新年度の配当金でlast_dividend_per_shareを更新
4. daily_quotesから最新取引データでassessmentテーブルを更新:
   - 銘柄ごとの最新取引日でlast_trade_dateを更新
   - 最新取引日の調整終値でlast_adjustment_closeを更新
5. 配当利回りを計算してlast_dividend_yieldを更新:
   - 配当金÷調整終値で配当利回りを算出
6. 3か月価格データと乖離率を更新:
   - 3か月間の最高・最低調整終値を計算
   - 最新終値から3か月最高値の乖離率: (最新終値-3か月最高値)/3か月最高値*100
   - 最新終値から3か月最低値の乖離率: (最新終値-3か月最低値)/3か月最低値*100
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("銘柄評価データ更新を開始します...")

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

		if err := ExecuteAssessment(dbService); err != nil {
			return err
		}

		fmt.Println("更新完了")
		return nil
	},
}

// AssessmentUpdateRepository 評価データ更新専用のリポジトリ
type AssessmentUpdateRepository struct {
	conn *database.Connection
}

// NewAssessmentUpdateRepository 新しい評価データ更新リポジトリを作成
func NewAssessmentUpdateRepository(conn *database.Connection) *AssessmentUpdateRepository {
	return &AssessmentUpdateRepository{conn: conn}
}

// DeleteAssessmentByOldEffectiveDate 古いeffective_dateを持つ銘柄のassessmentレコードを削除
func (r *AssessmentUpdateRepository) DeleteAssessmentByOldEffectiveDate() (int64, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		DELETE FROM assessment 
		WHERE code IN (
			SELECT DISTINCT code 
			FROM listed_info 
			WHERE effective_date < ?
		)
	`

	result, err := r.conn.GetDB().Exec(query, today)
	if err != nil {
		return 0, fmt.Errorf("古いeffective_dateを持つassessmentレコード削除エラー: %v", err)
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("削除件数取得エラー: %v", err)
	}

	return deletedCount, nil
}

// AddNewListedStocks 本日以降のeffective_dateを持つmarket_code 109以外の銘柄をassessmentテーブルに追加
func (r *AssessmentUpdateRepository) AddNewListedStocks() (int64, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO assessment (code, created_at, updated_at)
		SELECT DISTINCT li.code, NOW(), NOW()
		FROM listed_info li
		WHERE li.effective_date >= ?
		  AND li.market_code != '0109'
		  AND li.code NOT IN (SELECT code FROM assessment)
	`

	result, err := r.conn.GetDB().Exec(query, today)
	if err != nil {
		return 0, fmt.Errorf("新規銘柄のassessment追加エラー: %v", err)
	}

	addedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("追加件数取得エラー: %v", err)
	}

	return addedCount, nil
}

// UpdateFromStatementsSummary statements_summary関連データでassessmentテーブルを更新
func (r *AssessmentUpdateRepository) UpdateFromStatementsSummary() (int64, error) {
	query := `
		UPDATE assessment a
		LEFT JOIN (
			SELECT
				local_code,
				MAX(fiscal_year_end_date) AS fiscal_year_end_date
			FROM statements_summary
			GROUP BY local_code
		) ss_max ON a.code = ss_max.local_code
		LEFT JOIN statements_summary ss ON a.code = ss.local_code 
		                                AND ss_max.fiscal_year_end_date = ss.fiscal_year_end_date
		SET 
			a.last_fiscal_year_end_date = ss_max.fiscal_year_end_date,
			a.last_dividend_per_share = ss.dividend_per_share,
			a.updated_at = NOW()
	`

	result, err := r.conn.GetDB().Exec(query)
	if err != nil {
		return 0, fmt.Errorf("statements_summary関連データ更新エラー: %v", err)
	}

	updatedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("更新件数取得エラー: %v", err)
	}

	return updatedCount, nil
}

// UpdateFromDailyQuotes daily_quotes関連データでassessmentテーブルを更新
func (r *AssessmentUpdateRepository) UpdateFromDailyQuotes() (int64, error) {
	query := `
		UPDATE assessment a
		LEFT JOIN (
			SELECT
				code,
				MAX(trade_date) AS max_trade_date
			FROM daily_quotes
			GROUP BY code
		) dq_max ON a.code = dq_max.code
		LEFT JOIN daily_quotes dq ON a.code = dq.code 
		                          AND dq_max.max_trade_date = dq.trade_date
		SET 
			a.last_trade_date = dq.trade_date,
			a.last_adjustment_close = dq.adjustment_close,
			a.updated_at = NOW()
	`

	result, err := r.conn.GetDB().Exec(query)
	if err != nil {
		return 0, fmt.Errorf("daily_quotes関連データ更新エラー: %v", err)
	}

	updatedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("更新件数取得エラー: %v", err)
	}

	return updatedCount, nil
}

// UpdateDividendYield 配当利回りを計算してlast_dividend_yieldを更新
func (r *AssessmentUpdateRepository) UpdateDividendYield() (int64, error) {
	query := `
		UPDATE assessment 
		SET last_dividend_yield = CASE 
			WHEN last_dividend_per_share IS NOT NULL 
				AND last_dividend_per_share > 0 
				AND last_adjustment_close IS NOT NULL 
				AND last_adjustment_close > 0
			THEN (last_dividend_per_share / last_adjustment_close) * 100
			ELSE NULL
		END,
		updated_at = NOW()
		WHERE last_dividend_per_share IS NOT NULL 
		   OR last_adjustment_close IS NOT NULL
	`

	result, err := r.conn.GetDB().Exec(query)
	if err != nil {
		return 0, fmt.Errorf("配当利回り計算・更新エラー: %v", err)
	}

	updatedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("更新件数取得エラー: %v", err)
	}

	return updatedCount, nil
}

// UpdateThreeMonthPriceData 3か月間の最高・最低調整終値と乖離率を更新
func (r *AssessmentUpdateRepository) UpdateThreeMonthPriceData() (int64, error) {
	query := `
		UPDATE assessment a
		LEFT JOIN (
			SELECT 
				code,
				MAX(adjustment_close) AS max_close_3m,
				MIN(adjustment_close) AS min_close_3m
			FROM daily_quotes 
			WHERE trade_date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)
			  AND adjustment_close IS NOT NULL
			GROUP BY code
		) dq_3m ON a.code = dq_3m.code
		SET 
			a.three_month_max_close = dq_3m.max_close_3m,
			a.three_month_min_close = dq_3m.min_close_3m,
			a.deviation_from_max = CASE 
				WHEN a.last_adjustment_close IS NOT NULL 
					AND dq_3m.max_close_3m IS NOT NULL 
					AND dq_3m.max_close_3m > 0
				THEN ((a.last_adjustment_close - dq_3m.max_close_3m) / dq_3m.max_close_3m) * 100
				ELSE NULL
			END,
			a.deviation_from_min = CASE 
				WHEN a.last_adjustment_close IS NOT NULL 
					AND dq_3m.min_close_3m IS NOT NULL 
					AND dq_3m.min_close_3m > 0
				THEN ((a.last_adjustment_close - dq_3m.min_close_3m) / dq_3m.min_close_3m) * 100
				ELSE NULL
			END,
			a.updated_at = NOW()
		WHERE dq_3m.code IS NOT NULL
	`

	result, err := r.conn.GetDB().Exec(query)
	if err != nil {
		return 0, fmt.Errorf("3か月価格データ更新エラー: %v", err)
	}

	updatedCount, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("更新件数取得エラー: %v", err)
	}

	return updatedCount, nil
}

func init() {
	// 今後、必要に応じてフラグを追加予定
	// AssessmentCmd.Flags().StringVar(&assessmentCode, "code", "", "更新対象の銘柄コード（指定しない場合は全銘柄）")
	// AssessmentCmd.Flags().BoolVar(&assessmentDebug, "debug", false, "デバッグモードで実行")
}
