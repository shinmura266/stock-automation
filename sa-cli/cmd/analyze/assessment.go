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

// AssessmentResult 評価結果
type AssessmentResult struct {
	Code                  string     `db:"code"`
	CompanyName           string     `db:"company_name"`
	LastFiscalYearEndDate *time.Time `db:"last_fiscal_year_end_date"`
	LastDividendPerShare  *float64   `db:"last_dividend_per_share"`
	LastTradeDate         *time.Time `db:"last_trade_date"`
	LastAdjustmentClose   *float64   `db:"last_adjustment_close"`
	LastDividendYield     *float64   `db:"last_dividend_yield"`
}

var AssessmentCmd = &cobra.Command{
	Use:   "assessment",
	Short: "銘柄評価データを表示",
	Long: `
assessmentテーブルから銘柄評価データを表示します。

表示項目:
- 銘柄コード
- 企業名  
- 最終決算期末日
- 最終配当金（1株当たり）
- 最終取引日
- 最終調整終値
- 最終配当利回り
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("銘柄評価データ分析を開始します...")

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
		assessmentRepo := NewAssessmentRepository(dbService.GetConnection())

		fmt.Println("評価データを取得中...")
		results, err := assessmentRepo.GetAssessmentData()
		if err != nil {
			return errors.DatabaseError(err)
		}

		// 結果を表示
		fmt.Println("\n=== 銘柄評価データ ===")
		fmt.Printf("%-10s %-30s %-15s %-12s %-12s %-12s %-10s\n",
			"銘柄コード", "企業名", "決算期末日", "配当金", "取引日", "調整終値", "利回り%")
		fmt.Println("------------------------------------------------------------------------------------------------")

		for i, result := range results {
			var companyName, fiscalEndStr, dividendStr, tradeDateStr, closeStr, yieldStr string

			// 企業名の表示（30文字制限）
			if len(result.CompanyName) > 30 {
				companyName = result.CompanyName[:27] + "..."
			} else {
				companyName = result.CompanyName
			}

			// 決算期末日
			if result.LastFiscalYearEndDate != nil {
				fiscalEndStr = result.LastFiscalYearEndDate.Format("2006-01-02")
			} else {
				fiscalEndStr = "-"
			}

			// 配当金
			if result.LastDividendPerShare != nil {
				dividendStr = fmt.Sprintf("%.2f", *result.LastDividendPerShare)
			} else {
				dividendStr = "-"
			}

			// 取引日
			if result.LastTradeDate != nil {
				tradeDateStr = result.LastTradeDate.Format("2006-01-02")
			} else {
				tradeDateStr = "-"
			}

			// 調整終値
			if result.LastAdjustmentClose != nil {
				closeStr = fmt.Sprintf("%.0f", *result.LastAdjustmentClose)
			} else {
				closeStr = "-"
			}

			// 配当利回り
			if result.LastDividendYield != nil {
				yieldStr = fmt.Sprintf("%.4f", *result.LastDividendYield)
			} else {
				yieldStr = "-"
			}

			fmt.Printf("%3d. %-8s %-30s %-15s %-12s %-12s %-12s %-10s\n",
				i+1, result.Code, companyName, fiscalEndStr, dividendStr, tradeDateStr, closeStr, yieldStr)
		}

		fmt.Printf("\n分析完了: %d銘柄の評価データを表示しました\n", len(results))
		return nil
	},
}

// AssessmentRepository 評価データ専用のリポジトリ
type AssessmentRepository struct {
	conn *database.Connection
}

// NewAssessmentRepository 新しい評価リポジトリを作成
func NewAssessmentRepository(conn *database.Connection) *AssessmentRepository {
	return &AssessmentRepository{conn: conn}
}

// GetAssessmentData 評価データを取得
func (r *AssessmentRepository) GetAssessmentData() ([]*AssessmentResult, error) {
	query := `
		SELECT 
			a.code,
			COALESCE(li.company_name, '') as company_name,
			a.last_fiscal_year_end_date,
			a.last_dividend_per_share,
			a.last_trade_date,
			a.last_adjustment_close,
			a.last_dividend_yield
		FROM assessment a
		LEFT JOIN (
			-- 最新の銘柄情報を取得
			SELECT 
				code,
				company_name,
				ROW_NUMBER() OVER (PARTITION BY code ORDER BY effective_date DESC) as rn
			FROM listed_info
		) li ON a.code = li.code AND li.rn = 1
		ORDER BY a.code
	`

	rows, err := r.conn.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("評価データ取得エラー: %v", err)
	}
	defer rows.Close()

	var results []*AssessmentResult
	for rows.Next() {
		result := &AssessmentResult{}

		err := rows.Scan(
			&result.Code,
			&result.CompanyName,
			&result.LastFiscalYearEndDate,
			&result.LastDividendPerShare,
			&result.LastTradeDate,
			&result.LastAdjustmentClose,
			&result.LastDividendYield,
		)
		if err != nil {
			log.Printf("評価データスキャンエラー: %v", err)
			continue
		}

		results = append(results, result)
	}

	return results, nil
}
