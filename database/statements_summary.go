package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// StatementsSummary 財務情報サマリーの構造体
type StatementsSummary struct {
	LocalCode           string     `db:"local_code"`
	FiscalYearStartDate time.Time  `db:"fiscal_year_start_date"`
	FiscalYearEndDate   *time.Time `db:"fiscal_year_end_date"`
	DisclosedDate       time.Time  `db:"disclosed_date"`
	DisclosedTime       *time.Time `db:"disclosed_time"`
	TypeOfCurrentPeriod string     `db:"type_of_current_period"`
	NetSales            *int64     `db:"net_sales"`
	OperatingProfit     *int64     `db:"operating_profit"`
	OrdinaryProfit      *int64     `db:"ordinary_profit"`
	Profit              *int64     `db:"profit"`
	EPS                 *float64   `db:"eps"`
	TotalAssets         *int64     `db:"total_assets"`
	Equity              *int64     `db:"equity"`
	EquityToAssetRatio  *float64   `db:"equity_to_asset_ratio"`
	DividendPerShare    *float64   `db:"dividend_per_share"`
	IsForecast          bool       `db:"is_forecast"`
	DataType            string     `db:"data_type"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

// StatementsSummaryRepository 財務情報サマリーのリポジトリ
type StatementsSummaryRepository struct {
	conn *Connection
}

// NewStatementsSummaryRepository 新しいリポジトリを作成
func NewStatementsSummaryRepository(conn *Connection) *StatementsSummaryRepository {
	return &StatementsSummaryRepository{conn: conn}
}

// AnalyzeAndCreateSummary 財務データから各会計年度の最新サマリーを作成
func (r *StatementsSummaryRepository) AnalyzeAndCreateSummary() error {
	log.Println("財務情報サマリーの分析・作成を開始...")

	tx, cleanup := BeginTransaction(r.conn.GetDB())
	defer cleanup()

	// 現在のサマリーテーブルをクリア
	if _, err := tx.Exec("DELETE FROM statements_summary"); err != nil {
		tx.Rollback()
		return fmt.Errorf("サマリーテーブルクリアエラー: %v", err)
	}

	// 各local_codeごとに処理
	query := `
		SELECT DISTINCT local_code 
		FROM statements 
		ORDER BY local_code
	`
	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("銘柄コード取得エラー: %v", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			log.Printf("銘柄コードスキャンエラー: %v", err)
			continue
		}
		codes = append(codes, code)
	}

	log.Printf("処理対象銘柄数: %d", len(codes))

	processedCount := 0
	for _, code := range codes {
		if err := r.processLocalCode(tx, code); err != nil {
			log.Printf("銘柄 %s の処理でエラー: %v", code, err)
			continue
		}
		processedCount++
		if processedCount%100 == 0 {
			log.Printf("処理済み銘柄数: %d/%d", processedCount, len(codes))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	log.Printf("財務情報サマリーの作成完了: %d銘柄処理", processedCount)
	return nil
}

// processLocalCode 個別銘柄の財務データ処理
func (r *StatementsSummaryRepository) processLocalCode(tx *sql.Tx, localCode string) error {
	// 各会計年度の最新データを取得
	summaries, err := r.getLatestStatementsByFiscalYear(tx, localCode)
	if err != nil {
		return fmt.Errorf("会計年度別最新データ取得エラー: %v", err)
	}

	// 会計年度別に最適なデータを選択
	fiscalYearBest := r.selectBestDataPerFiscalYear(summaries)

	// サマリーデータを挿入
	for _, summary := range fiscalYearBest {
		if err := r.insertSummary(tx, summary); err != nil {
			log.Printf("サマリー挿入エラー (銘柄: %s, 年度: %s): %v",
				localCode, summary.FiscalYearStartDate.Format("2006-01-02"), err)
		}
	}

	return nil
}

// selectBestDataPerFiscalYear 会計年度別に最適なデータを選択
func (r *StatementsSummaryRepository) selectBestDataPerFiscalYear(summaries []*StatementsSummary) []*StatementsSummary {
	fiscalYearMap := make(map[string]*StatementsSummary)

	for _, summary := range summaries {
		fiscalKey := summary.FiscalYearStartDate.Format("2006-01-02")

		existing, exists := fiscalYearMap[fiscalKey]
		if !exists {
			fiscalYearMap[fiscalKey] = summary
			continue
		}

		// より優先度の高いデータを選択
		if r.isBetterData(summary, existing) {
			fiscalYearMap[fiscalKey] = summary
		}
	}

	var result []*StatementsSummary
	for _, summary := range fiscalYearMap {
		result = append(result, summary)
	}

	return result
}

// isBetterData より良いデータかを判定
func (r *StatementsSummaryRepository) isBetterData(new, existing *StatementsSummary) bool {
	// 1. より新しい開示日を優先
	if new.DisclosedDate.After(existing.DisclosedDate) {
		return true
	}
	if new.DisclosedDate.Before(existing.DisclosedDate) {
		return false
	}

	// 2. 同じ開示日の場合、データタイプ優先順位：current_actual > current_forecast > next_year_forecast
	priority := map[string]int{
		"current_actual":     3,
		"current_forecast":   2,
		"next_year_forecast": 1,
	}

	newPriority := priority[new.DataType]
	existingPriority := priority[existing.DataType]

	return newPriority > existingPriority
}

// getLatestStatementsByFiscalYear 会計年度別の最新財務データを取得
func (r *StatementsSummaryRepository) getLatestStatementsByFiscalYear(tx *sql.Tx, localCode string) ([]*StatementsSummary, error) {
	// 当期実績データ（current_fiscal_year_start_dateが存在し、確定データ）
	currentActualQuery := `
		SELECT 
			local_code,
			current_fiscal_year_start_date,
			current_fiscal_year_end_date,
			disclosed_date,
			disclosed_time,
			type_of_current_period,
			net_sales,
			operating_profit,
			ordinary_profit,
			profit,
			eps,
			total_assets,
			equity,
			equity_to_asset_ratio,
			COALESCE(result_dps_annual, 0) as dividend_per_share
		FROM statements 
		WHERE local_code = ? 
			AND current_fiscal_year_start_date IS NOT NULL
			AND (net_sales IS NOT NULL OR operating_profit IS NOT NULL)
			AND type_of_current_period IN ('FY', 'Q4')
		ORDER BY current_fiscal_year_start_date DESC, disclosed_date DESC
	`

	// 当期予想データ
	currentForecastQuery := `
		SELECT 
			local_code,
			current_fiscal_year_start_date,
			current_fiscal_year_end_date,
			disclosed_date,
			disclosed_time,
			type_of_current_period,
			COALESCE(fc_net_sales, net_sales) as net_sales,
			COALESCE(fc_operating_profit, operating_profit) as operating_profit,
			COALESCE(fc_ordinary_profit, ordinary_profit) as ordinary_profit,
			COALESCE(fc_profit, profit) as profit,
			COALESCE(fc_eps, eps) as eps,
			total_assets,
			equity,
			equity_to_asset_ratio,
			COALESCE(fc_dps_annual, result_dps_annual, 0) as dividend_per_share
		FROM statements 
		WHERE local_code = ? 
			AND current_fiscal_year_start_date IS NOT NULL
			AND (fc_net_sales IS NOT NULL OR fc_operating_profit IS NOT NULL OR net_sales IS NOT NULL)
		ORDER BY current_fiscal_year_start_date DESC, disclosed_date DESC
	`

	// 翌期予想データ
	nextYearForecastQuery := `
		SELECT 
			local_code,
			next_fiscal_year_start_date as fiscal_year_start_date,
			next_fiscal_year_end_date as fiscal_year_end_date,
			disclosed_date,
			disclosed_time,
			type_of_current_period,
			ny_fc_net_sales as net_sales,
			ny_fc_operating_profit as operating_profit,
			ny_fc_ordinary_profit as ordinary_profit,
			ny_fc_profit as profit,
			ny_fc_eps as eps,
			NULL as total_assets,
			NULL as equity,
			NULL as equity_to_asset_ratio,
			COALESCE(ny_fc_dps_fy, 0) as dividend_per_share
		FROM statements 
		WHERE local_code = ? 
			AND next_fiscal_year_start_date IS NOT NULL
			AND (ny_fc_net_sales IS NOT NULL OR ny_fc_operating_profit IS NOT NULL)
		ORDER BY next_fiscal_year_start_date DESC, disclosed_date DESC
	`

	var summaries []*StatementsSummary

	// 当期実績データを処理
	if err := r.processDataType(tx, currentActualQuery, localCode, "current_actual", false, &summaries); err != nil {
		return nil, err
	}

	// 当期予想データを処理
	if err := r.processDataType(tx, currentForecastQuery, localCode, "current_forecast", true, &summaries); err != nil {
		return nil, err
	}

	// 翌期予想データを処理
	if err := r.processDataType(tx, nextYearForecastQuery, localCode, "next_year_forecast", true, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

// processDataType データタイプ別の処理
func (r *StatementsSummaryRepository) processDataType(tx *sql.Tx, query string, localCode string, dataType string, isForecast bool, summaries *[]*StatementsSummary) error {
	rows, err := tx.Query(query, localCode)
	if err != nil {
		return fmt.Errorf("%sデータ取得エラー: %v", dataType, err)
	}
	defer rows.Close()

	// 会計年度別に最新データを管理
	fiscalYearMap := make(map[string]*StatementsSummary)

	for rows.Next() {
		summary := &StatementsSummary{
			IsForecast: isForecast,
			DataType:   dataType,
		}

		var fiscalYearStartDate, fiscalYearEndDate, disclosedDate interface{}
		var disclosedTime interface{}

		err := rows.Scan(
			&summary.LocalCode,
			&fiscalYearStartDate,
			&fiscalYearEndDate,
			&disclosedDate,
			&disclosedTime,
			&summary.TypeOfCurrentPeriod,
			&summary.NetSales,
			&summary.OperatingProfit,
			&summary.OrdinaryProfit,
			&summary.Profit,
			&summary.EPS,
			&summary.TotalAssets,
			&summary.Equity,
			&summary.EquityToAssetRatio,
			&summary.DividendPerShare,
		)
		if err != nil {
			log.Printf("データスキャンエラー: %v", err)
			continue
		}

		// 日付変換
		if fiscalYearStartDate != nil {
			if t, ok := fiscalYearStartDate.(time.Time); ok {
				summary.FiscalYearStartDate = t
			}
		}
		if fiscalYearEndDate != nil {
			if t, ok := fiscalYearEndDate.(time.Time); ok {
				summary.FiscalYearEndDate = &t
			}
		}
		if disclosedDate != nil {
			if t, ok := disclosedDate.(time.Time); ok {
				summary.DisclosedDate = t
			}
		}
		if disclosedTime != nil {
			if t, ok := disclosedTime.(time.Time); ok {
				summary.DisclosedTime = &t
			}
		}

		// 会計年度キー
		fiscalKey := summary.FiscalYearStartDate.Format("2006-01-02")

		// より新しい開示日のデータがあれば更新
		if existing, exists := fiscalYearMap[fiscalKey]; !exists || summary.DisclosedDate.After(existing.DisclosedDate) {
			fiscalYearMap[fiscalKey] = summary
		}
	}

	// マップからスライスに変換
	for _, summary := range fiscalYearMap {
		*summaries = append(*summaries, summary)
	}

	return nil
}

// insertSummary サマリーデータを挿入
func (r *StatementsSummaryRepository) insertSummary(tx *sql.Tx, summary *StatementsSummary) error {
	query := `
		INSERT INTO statements_summary (
			local_code, fiscal_year_start_date, fiscal_year_end_date,
			disclosed_date, disclosed_time, type_of_current_period,
			net_sales, operating_profit, ordinary_profit, profit, eps,
			total_assets, equity, equity_to_asset_ratio, dividend_per_share,
			is_forecast, data_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			fiscal_year_end_date = VALUES(fiscal_year_end_date),
			disclosed_date = VALUES(disclosed_date),
			disclosed_time = VALUES(disclosed_time),
			type_of_current_period = VALUES(type_of_current_period),
			net_sales = VALUES(net_sales),
			operating_profit = VALUES(operating_profit),
			ordinary_profit = VALUES(ordinary_profit),
			profit = VALUES(profit),
			eps = VALUES(eps),
			total_assets = VALUES(total_assets),
			equity = VALUES(equity),
			equity_to_asset_ratio = VALUES(equity_to_asset_ratio),
			dividend_per_share = VALUES(dividend_per_share),
			is_forecast = VALUES(is_forecast),
			data_type = VALUES(data_type),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := tx.Exec(query,
		summary.LocalCode,
		summary.FiscalYearStartDate,
		summary.FiscalYearEndDate,
		summary.DisclosedDate,
		summary.DisclosedTime,
		summary.TypeOfCurrentPeriod,
		summary.NetSales,
		summary.OperatingProfit,
		summary.OrdinaryProfit,
		summary.Profit,
		summary.EPS,
		summary.TotalAssets,
		summary.Equity,
		summary.EquityToAssetRatio,
		summary.DividendPerShare,
		summary.IsForecast,
		summary.DataType,
	)

	return err
}

// GetSummaryByCode 銘柄コード別のサマリーデータを取得
func (r *StatementsSummaryRepository) GetSummaryByCode(localCode string) ([]*StatementsSummary, error) {
	query := `
		SELECT 
			local_code, fiscal_year_start_date, fiscal_year_end_date,
			disclosed_date, disclosed_time, type_of_current_period,
			net_sales, operating_profit, ordinary_profit, profit, eps,
			total_assets, equity, equity_to_asset_ratio, dividend_per_share,
			is_forecast, data_type, created_at, updated_at
		FROM statements_summary
		WHERE local_code = ?
		ORDER BY fiscal_year_start_date DESC, data_type
	`

	rows, err := r.conn.GetDB().Query(query, localCode)
	if err != nil {
		return nil, fmt.Errorf("サマリーデータ取得エラー: %v", err)
	}
	defer rows.Close()

	var summaries []*StatementsSummary
	for rows.Next() {
		summary := &StatementsSummary{}
		err := rows.Scan(
			&summary.LocalCode,
			&summary.FiscalYearStartDate,
			&summary.FiscalYearEndDate,
			&summary.DisclosedDate,
			&summary.DisclosedTime,
			&summary.TypeOfCurrentPeriod,
			&summary.NetSales,
			&summary.OperatingProfit,
			&summary.OrdinaryProfit,
			&summary.Profit,
			&summary.EPS,
			&summary.TotalAssets,
			&summary.Equity,
			&summary.EquityToAssetRatio,
			&summary.DividendPerShare,
			&summary.IsForecast,
			&summary.DataType,
			&summary.CreatedAt,
			&summary.UpdatedAt,
		)
		if err != nil {
			log.Printf("サマリーデータスキャンエラー: %v", err)
			continue
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
