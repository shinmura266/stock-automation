package analyze

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

// DividendYieldResult 配当利回り分析結果
type DividendYieldResult struct {
	LocalCode             string     `db:"local_code"`
	CompanyName           string     `db:"company_name"`
	LastFiscalYearEndDate *time.Time `db:"last_fiscal_year_end_date"`
	LastDividendPerShare  *float64   `db:"last_dividend_per_share"`
	LastTradeDate         *time.Time `db:"last_trade_date"`
	LastAdjustmentClose   *float64   `db:"last_adjustment_close"`
	LastDividendYield     *float64   `db:"last_dividend_yield"`
	DeviationFromMin      *float64   `db:"deviation_from_min"`
	DeviationFromMax      *float64   `db:"deviation_from_max"`
}

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "配当利回り4%以上かつ最高値乖離率-10%以下の銘柄を3か月最安値乖離率順で分析・表示",
	Long: `
assessmentテーブルから配当利回り4%以上かつ3か月最高値乖離率-10%以下の銘柄を取得し、
3か月最安値からの乖離率が小さい順に上位20銘柄を表示します。

処理内容:
- assessmentテーブルから配当利回り4%以上の銘柄を取得
- 3か月最高値乖離率-10%以下の銘柄に絞り込み
- listed_infoテーブルから企業名を取得
- deviation_from_minの昇順で上位20銘柄を表示

表示項目:
- 銘柄コード
- 企業名
- 最終会計年度終了日
- 最終配当金
- 最終取引日
- 最終調整終値
- 最終配当利回り (%)
- 3か月最安値乖離率 (%)
- 3か月最高値乖離率 (%)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("配当利回り分析を開始します...")

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
		queryRepo := NewQueryRepository(dbService.GetConnection())

		fmt.Println("配当利回り4%以上、3か月最安値乖離率の小さい銘柄を分析中...")
		results, err := queryRepo.GetTopDividendYieldStocks(20)
		if err != nil {
			return errors.DatabaseError(err)
		}

		// 結果を表示
		fmt.Println("\n=== 配当利回り4%以上・最高値乖離率-10%以下・3か月最安値乖離率順20銘柄 ===")
		fmt.Printf("%-10s %-30s %-15s %-10s %-12s %-10s %-8s %-8s %-8s\n",
			"銘柄コード", "企業名", "会計年度終了日", "配当金", "取引日", "調整終値", "利回り%", "最安乖離%", "最高乖離%")
		fmt.Println("----------------------------------------------------------------------------------------------------------------")

		for i, result := range results {
			var companyName, fiscalYearStr, dividendStr, priceStr, tradeDateStr, yieldStr, deviationMinStr, deviationMaxStr string

			// 企業名の表示（30文字制限）
			if len(result.CompanyName) > 30 {
				companyName = result.CompanyName[:27] + "..."
			} else {
				companyName = result.CompanyName
			}

			// 会計年度終了日
			if result.LastFiscalYearEndDate != nil {
				fiscalYearStr = result.LastFiscalYearEndDate.Format("2006-01-02")
			} else {
				fiscalYearStr = "-"
			}

			// 配当金
			if result.LastDividendPerShare != nil {
				dividendStr = fmt.Sprintf("%.2f", *result.LastDividendPerShare)
			} else {
				dividendStr = "-"
			}

			// 調整終値
			if result.LastAdjustmentClose != nil {
				priceStr = fmt.Sprintf("%.0f", *result.LastAdjustmentClose)
			} else {
				priceStr = "-"
			}

			// 取引日
			if result.LastTradeDate != nil {
				tradeDateStr = result.LastTradeDate.Format("2006-01-02")
			} else {
				tradeDateStr = "-"
			}

			// 配当利回り
			if result.LastDividendYield != nil {
				yieldStr = fmt.Sprintf("%.2f", *result.LastDividendYield)
			} else {
				yieldStr = "-"
			}

			// 3か月最安値乖離率
			if result.DeviationFromMin != nil {
				deviationMinStr = fmt.Sprintf("%.2f", *result.DeviationFromMin)
			} else {
				deviationMinStr = "-"
			}

			// 3か月最高値乖離率
			if result.DeviationFromMax != nil {
				deviationMaxStr = fmt.Sprintf("%.2f", *result.DeviationFromMax)
			} else {
				deviationMaxStr = "-"
			}

			// 銘柄コードを4桁に変換（最後の0を削除）
			displayCode := result.LocalCode
			if len(displayCode) == 5 && displayCode[4] == '0' {
				displayCode = displayCode[:4]
			}

			fmt.Printf("%2d. %-8s %-30s %-15s %-10s %-12s %-10s %-8s %-8s %-8s\n",
				i+1, displayCode, companyName, fiscalYearStr, dividendStr, tradeDateStr, priceStr, yieldStr, deviationMinStr, deviationMaxStr)
		}

		fmt.Printf("\n分析完了: 配当利回り4%%以上で%d銘柄の乖離率データを表示しました\n", len(results))
		return nil
	},
}

// QueryRepository クエリ専用のリポジトリ
type QueryRepository struct {
	conn *database.Connection
}

// NewQueryRepository 新しいクエリリポジトリを作成
func NewQueryRepository(conn *database.Connection) *QueryRepository {
	return &QueryRepository{conn: conn}
}

// GetTopDividendYieldStocks 配当利回り4%以上かつ最高値乖離率-10%以下で3か月最安値乖離率が小さい銘柄を取得
func (r *QueryRepository) GetTopDividendYieldStocks(limit int) ([]*DividendYieldResult, error) {
	query := `
		SELECT 
			a.code,
			COALESCE(li.company_name, '') as company_name,
			a.last_fiscal_year_end_date,
			a.last_dividend_per_share,
			a.last_trade_date,
			a.last_adjustment_close,
			a.last_dividend_yield,
			a.deviation_from_min,
			a.deviation_from_max
		FROM assessment a
		LEFT JOIN (
			-- 最新の銘柄情報を取得
			SELECT 
				code,
				company_name,
				ROW_NUMBER() OVER (PARTITION BY code ORDER BY effective_date DESC) as rn
			FROM listed_info
		) li ON a.code = li.code AND li.rn = 1
		WHERE a.last_dividend_yield IS NOT NULL 
		  AND a.last_dividend_yield >= 4.0
		  AND a.deviation_from_min IS NOT NULL
		  AND a.deviation_from_max IS NOT NULL
		  AND a.deviation_from_max <= -10.0
		ORDER BY a.deviation_from_min ASC
		LIMIT ?
	`

	rows, err := r.conn.GetDB().Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("配当利回り4%%以上最高値乖離率-10%%以下銘柄取得エラー: %v", err)
	}
	defer rows.Close()

	var results []*DividendYieldResult
	for rows.Next() {
		result := &DividendYieldResult{}

		err := rows.Scan(
			&result.LocalCode,
			&result.CompanyName,
			&result.LastFiscalYearEndDate,
			&result.LastDividendPerShare,
			&result.LastTradeDate,
			&result.LastAdjustmentClose,
			&result.LastDividendYield,
			&result.DeviationFromMin,
			&result.DeviationFromMax,
		)
		if err != nil {
			log.Printf("配当利回りデータスキャンエラー: %v", err)
			continue
		}

		results = append(results, result)
	}

	return results, nil
}
