package database

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// StatementsRepository 財務情報のリポジトリ
type StatementsRepository struct {
	conn *Connection
}

// NewStatementsRepository 新しいリポジトリを作成
func NewStatementsRepository(conn *Connection) *StatementsRepository {
	return &StatementsRepository{conn: conn}
}

// SaveFinancialStatements 財務情報をデータベースに保存（ジェネリック対応）
func (r *StatementsRepository) SaveFinancialStatements(statements interface{}) (int, error) {
	// リフレクションを使用して動的に型を処理
	v := reflect.ValueOf(statements)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Statements配列を取得
	statementsField := v.FieldByName("Statements")
	if !statementsField.IsValid() || statementsField.Kind() != reflect.Slice {
		return 0, fmt.Errorf("Statementsフィールドが見つからないか、スライスではありません")
	}

	if statementsField.Len() == 0 {
		return 0, fmt.Errorf("保存するデータがありません")
	}

	tx, cleanup := BeginTransaction(r.conn.GetDB())
	defer cleanup()

	stmt, err := tx.Prepare(`
		INSERT INTO statements (
			disclosed_date, disclosed_time, local_code, disclosure_number, type_of_document,
			type_of_current_period, current_period_start_date, current_period_end_date,
			current_fiscal_year_start_date, current_fiscal_year_end_date,
			next_fiscal_year_start_date, next_fiscal_year_end_date,
			net_sales, operating_profit, ordinary_profit, profit,
			eps, diluted_eps, total_assets, equity,
			equity_to_asset_ratio, bvps,
			cf_operating, cf_investing,
			cf_financing, cash_and_equivalents,
			result_dps_1q, result_dps_2q,
			result_dps_3q, result_dps_fy,
			result_dps_annual,
			distributions_per_unit_reit, result_total_dividend_annual, result_payout_ratio_annual,
			fc_dps_1q, fc_dps_2q,
			fc_dps_3q, fc_dps_fy,
			fc_dps_annual,
			fc_distributions_per_unit_reit, fc_total_dividend_annual, fc_payout_ratio_annual,
			ny_fc_dps_1q, ny_fc_dps_2q,
			ny_fc_dps_3q, ny_fc_dps_fy,
			ny_fc_distributions_per_unit_reit, ny_fc_payout_ratio_annual,
			fc_net_sales_2q, fc_operating_profit_2q,
			fc_ordinary_profit_2q, fc_profit_2q, fc_eps_2q,
			ny_fc_net_sales_2q, ny_fc_operating_profit_2q,
			ny_fc_ordinary_profit_2q, ny_fc_profit_2q,
			ny_fc_eps_2q,
			fc_net_sales, fc_operating_profit, fc_ordinary_profit,
			fc_profit, fc_eps,
			ny_fc_net_sales, ny_fc_operating_profit,
			ny_fc_ordinary_profit, ny_fc_profit, ny_fc_eps,
			material_changes_subsidiaries, significant_changes_consolidation_scope,
			changes_accounting_std_revisions, changes_accounting_std_other,
			changes_accounting_estimates, retrospective_restatement,
			issued_shares_end_fy_incl_treasury,
			treasury_shares_end_fy, avg_shares,
			nc_net_sales, nc_operating_profit, nc_ordinary_profit,
			nc_profit, nc_eps, nc_total_assets,
			nc_equity, nc_equity_to_asset_ratio, nc_bvps,
			fc_nc_net_sales_2q, fc_nc_operating_profit_2q, fc_nc_ordinary_profit_2q,
			fc_nc_profit_2q, fc_nc_eps_2q,
			ny_fc_nc_net_sales_2q, ny_fc_nc_operating_profit_2q, ny_fc_nc_ordinary_profit_2q,
			ny_fc_nc_profit_2q, ny_fc_nc_eps_2q,
			fc_nc_net_sales, fc_nc_operating_profit, fc_nc_ordinary_profit,
			fc_nc_profit, fc_nc_eps,
			ny_fc_nc_net_sales, ny_fc_nc_operating_profit, ny_fc_nc_ordinary_profit,
			ny_fc_nc_profit, ny_fc_nc_eps
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			disclosed_time = VALUES(disclosed_time),
			disclosure_number = VALUES(disclosure_number),
			type_of_document = VALUES(type_of_document),
			current_period_start_date = VALUES(current_period_start_date),
			current_period_end_date = VALUES(current_period_end_date),
			current_fiscal_year_start_date = VALUES(current_fiscal_year_start_date),
			current_fiscal_year_end_date = VALUES(current_fiscal_year_end_date),
			next_fiscal_year_start_date = VALUES(next_fiscal_year_start_date),
			next_fiscal_year_end_date = VALUES(next_fiscal_year_end_date),
			net_sales = VALUES(net_sales),
			operating_profit = VALUES(operating_profit),
			ordinary_profit = VALUES(ordinary_profit),
			profit = VALUES(profit),
			eps = VALUES(eps),
			diluted_eps = VALUES(diluted_eps),
			total_assets = VALUES(total_assets),
			equity = VALUES(equity),
			equity_to_asset_ratio = VALUES(equity_to_asset_ratio),
			bvps = VALUES(bvps),
			cf_operating = VALUES(cf_operating),
			cf_investing = VALUES(cf_investing),
			cf_financing = VALUES(cf_financing),
			cash_and_equivalents = VALUES(cash_and_equivalents),
			result_dps_1q = VALUES(result_dps_1q),
			result_dps_2q = VALUES(result_dps_2q),
			result_dps_3q = VALUES(result_dps_3q),
			result_dps_fy = VALUES(result_dps_fy),
			result_dps_annual = VALUES(result_dps_annual),
			distributions_per_unit_reit = VALUES(distributions_per_unit_reit),
			result_total_dividend_annual = VALUES(result_total_dividend_annual),
			result_payout_ratio_annual = VALUES(result_payout_ratio_annual),
			fc_dps_1q = VALUES(fc_dps_1q),
			fc_dps_2q = VALUES(fc_dps_2q),
			fc_dps_3q = VALUES(fc_dps_3q),
			fc_dps_fy = VALUES(fc_dps_fy),
			fc_dps_annual = VALUES(fc_dps_annual),
			fc_distributions_per_unit_reit = VALUES(fc_distributions_per_unit_reit),
			fc_total_dividend_annual = VALUES(fc_total_dividend_annual),
			fc_payout_ratio_annual = VALUES(fc_payout_ratio_annual),
			ny_fc_dps_1q = VALUES(ny_fc_dps_1q),
			ny_fc_dps_2q = VALUES(ny_fc_dps_2q),
			ny_fc_dps_3q = VALUES(ny_fc_dps_3q),
			ny_fc_dps_fy = VALUES(ny_fc_dps_fy),
			ny_fc_distributions_per_unit_reit = VALUES(ny_fc_distributions_per_unit_reit),
			ny_fc_payout_ratio_annual = VALUES(ny_fc_payout_ratio_annual),
			fc_net_sales_2q = VALUES(fc_net_sales_2q),
			fc_operating_profit_2q = VALUES(fc_operating_profit_2q),
			fc_ordinary_profit_2q = VALUES(fc_ordinary_profit_2q),
			fc_profit_2q = VALUES(fc_profit_2q),
			fc_eps_2q = VALUES(fc_eps_2q),
			ny_fc_net_sales_2q = VALUES(ny_fc_net_sales_2q),
			ny_fc_operating_profit_2q = VALUES(ny_fc_operating_profit_2q),
			ny_fc_ordinary_profit_2q = VALUES(ny_fc_ordinary_profit_2q),
			ny_fc_profit_2q = VALUES(ny_fc_profit_2q),
			ny_fc_eps_2q = VALUES(ny_fc_eps_2q),
			fc_net_sales = VALUES(fc_net_sales),
			fc_operating_profit = VALUES(fc_operating_profit),
			fc_ordinary_profit = VALUES(fc_ordinary_profit),
			fc_profit = VALUES(fc_profit),
			fc_eps = VALUES(fc_eps),
			ny_fc_net_sales = VALUES(ny_fc_net_sales),
			ny_fc_operating_profit = VALUES(ny_fc_operating_profit),
			ny_fc_ordinary_profit = VALUES(ny_fc_ordinary_profit),
			ny_fc_profit = VALUES(ny_fc_profit),
			ny_fc_eps = VALUES(ny_fc_eps),
			material_changes_subsidiaries = VALUES(material_changes_subsidiaries),
			significant_changes_consolidation_scope = VALUES(significant_changes_consolidation_scope),
			changes_accounting_std_revisions = VALUES(changes_accounting_std_revisions),
			changes_accounting_std_other = VALUES(changes_accounting_std_other),
			changes_accounting_estimates = VALUES(changes_accounting_estimates),
			retrospective_restatement = VALUES(retrospective_restatement),
			issued_shares_end_fy_incl_treasury = VALUES(issued_shares_end_fy_incl_treasury),
			treasury_shares_end_fy = VALUES(treasury_shares_end_fy),
			avg_shares = VALUES(avg_shares),
			nc_net_sales = VALUES(nc_net_sales),
			nc_operating_profit = VALUES(nc_operating_profit),
			nc_ordinary_profit = VALUES(nc_ordinary_profit),
			nc_profit = VALUES(nc_profit),
			nc_eps = VALUES(nc_eps),
			nc_total_assets = VALUES(nc_total_assets),
			nc_equity = VALUES(nc_equity),
			nc_equity_to_asset_ratio = VALUES(nc_equity_to_asset_ratio),
			nc_bvps = VALUES(nc_bvps),
			fc_nc_net_sales_2q = VALUES(fc_nc_net_sales_2q),
			fc_nc_operating_profit_2q = VALUES(fc_nc_operating_profit_2q),
			fc_nc_ordinary_profit_2q = VALUES(fc_nc_ordinary_profit_2q),
			fc_nc_profit_2q = VALUES(fc_nc_profit_2q),
			fc_nc_eps_2q = VALUES(fc_nc_eps_2q),
			ny_fc_nc_net_sales_2q = VALUES(ny_fc_nc_net_sales_2q),
			ny_fc_nc_operating_profit_2q = VALUES(ny_fc_nc_operating_profit_2q),
			ny_fc_nc_ordinary_profit_2q = VALUES(ny_fc_nc_ordinary_profit_2q),
			ny_fc_nc_profit_2q = VALUES(ny_fc_nc_profit_2q),
			ny_fc_nc_eps_2q = VALUES(ny_fc_nc_eps_2q),
			fc_nc_net_sales = VALUES(fc_nc_net_sales),
			fc_nc_operating_profit = VALUES(fc_nc_operating_profit),
			fc_nc_ordinary_profit = VALUES(fc_nc_ordinary_profit),
			fc_nc_profit = VALUES(fc_nc_profit),
			fc_nc_eps = VALUES(fc_nc_eps),
			ny_fc_nc_net_sales = VALUES(ny_fc_nc_net_sales),
			ny_fc_nc_operating_profit = VALUES(ny_fc_nc_operating_profit),
			ny_fc_nc_ordinary_profit = VALUES(ny_fc_nc_ordinary_profit),
			ny_fc_nc_profit = VALUES(ny_fc_nc_profit),
			ny_fc_nc_eps = VALUES(ny_fc_nc_eps),
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("プリペアドステートメント作成エラー: %v", err)
	}
	defer stmt.Close()

	insertedCount := 0
	updatedCount := 0

	// 各ステートメントを処理
	for i := 0; i < statementsField.Len(); i++ {
		stmtItem := statementsField.Index(i)

		// リフレクションでフィールドを取得
		disclosedDate := GetStringField(stmtItem, "DisclosedDate")
		disclosedTime := GetStringField(stmtItem, "DisclosedTime")
		localCode := GetStringField(stmtItem, "LocalCode")

		// 日付をパース
		parsedDate, err := time.Parse("2006-01-02", disclosedDate)
		if err != nil {
			log.Printf("開示日パースエラー (コード: %s, 日付: %s): %v", localCode, disclosedDate, err)
			continue
		}

		// 時刻をパース（空の場合はnilを設定）
		var parsedTime interface{}
		if disclosedTime != "" {
			if t, err := time.Parse("15:04:05", disclosedTime); err == nil {
				parsedTime = t.Format("15:04:05")
			}
		}

		// 期間開始日・終了日をパース
		var currentPeriodStartDate, currentPeriodEndDate interface{}
		var currentFiscalYearStartDate, currentFiscalYearEndDate interface{}
		var nextFiscalYearStartDate, nextFiscalYearEndDate interface{}

		if currentPSD := GetStringField(stmtItem, "CurrentPeriodStartDate"); currentPSD != "" {
			if d, err := time.Parse("2006-01-02", currentPSD); err == nil {
				currentPeriodStartDate = d
			}
		}
		if currentPED := GetStringField(stmtItem, "CurrentPeriodEndDate"); currentPED != "" {
			if d, err := time.Parse("2006-01-02", currentPED); err == nil {
				currentPeriodEndDate = d
			}
		}
		if currentFYSD := GetStringField(stmtItem, "CurrentFiscalYearStartDate"); currentFYSD != "" {
			if d, err := time.Parse("2006-01-02", currentFYSD); err == nil {
				currentFiscalYearStartDate = d
			}
		}
		if currentFYED := GetStringField(stmtItem, "CurrentFiscalYearEndDate"); currentFYED != "" {
			if d, err := time.Parse("2006-01-02", currentFYED); err == nil {
				currentFiscalYearEndDate = d
			}
		}
		if nextFYSD := GetStringField(stmtItem, "NextFiscalYearStartDate"); nextFYSD != "" {
			if d, err := time.Parse("2006-01-02", nextFYSD); err == nil {
				nextFiscalYearStartDate = d
			}
		}
		if nextFYED := GetStringField(stmtItem, "NextFiscalYearEndDate"); nextFYED != "" {
			if d, err := time.Parse("2006-01-02", nextFYED); err == nil {
				nextFiscalYearEndDate = d
			}
		}

		result, err := stmt.Exec(
			parsedDate, parsedTime, localCode,
			GetStringField(stmtItem, "DisclosureNumber"),
			GetStringField(stmtItem, "TypeOfDocument"),
			GetStringField(stmtItem, "TypeOfCurrentPeriod"),
			currentPeriodStartDate, currentPeriodEndDate,
			currentFiscalYearStartDate, currentFiscalYearEndDate,
			nextFiscalYearStartDate, nextFiscalYearEndDate,
			NullIfEmpty(GetStringField(stmtItem, "NetSales")),
			NullIfEmpty(GetStringField(stmtItem, "OperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "OrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "Profit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "EarningsPerShare")),
			NullIfEmptyFloat(GetStringField(stmtItem, "DilutedEarningsPerShare")),
			NullIfEmpty(GetStringField(stmtItem, "TotalAssets")),
			NullIfEmpty(GetStringField(stmtItem, "Equity")),
			NullIfEmptyFloat(GetStringField(stmtItem, "EquityToAssetRatio")),
			NullIfEmptyFloat(GetStringField(stmtItem, "BookValuePerShare")),
			NullIfEmpty(GetStringField(stmtItem, "CashFlowsFromOperatingActivities")),
			NullIfEmpty(GetStringField(stmtItem, "CashFlowsFromInvestingActivities")),
			NullIfEmpty(GetStringField(stmtItem, "CashFlowsFromFinancingActivities")),
			NullIfEmpty(GetStringField(stmtItem, "CashAndEquivalents")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultDividendPerShare1StQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultDividendPerShare2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultDividendPerShare3RdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultDividendPerShareFY")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultDividendPerShareAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "DistributionsPerUnitREIT")),
			NullIfEmpty(GetStringField(stmtItem, "ResultTotalDividendPaidAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ResultPayoutRatioAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDividendPerShare1StQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDividendPerShare2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDividendPerShare3RdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDividendPerShareFY")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDividendPerShareAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastDistributionsPerUnitREIT")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastTotalDividendPaidAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastPayoutRatioAnnual")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastDividendPerShare1StQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastDividendPerShare2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastDividendPerShare3RdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastDividendPerShareFY")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastDistributionsPerUnitREIT")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastPayoutRatioAnnual")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNetSales2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastOperatingProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastOrdinaryProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastProfit2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastEarningsPerShare2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNetSales2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastOperatingProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastOrdinaryProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastProfit2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastEarningsPerShare2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNetSales")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastOperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastOrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastProfit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastEarningsPerShare")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNetSales")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastOperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastOrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastProfit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastEarningsPerShare")),
			GetStringField(stmtItem, "MaterialChangesInSubsidiaries"),
			GetStringField(stmtItem, "SignificantChangesInTheScopeOfConsolidation"),
			GetStringField(stmtItem, "ChangesBasedOnRevisionsOfAccountingStandard"),
			GetStringField(stmtItem, "ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"),
			GetStringField(stmtItem, "ChangesInAccountingEstimates"),
			GetStringField(stmtItem, "RetrospectiveRestatement"),
			NullIfEmpty(GetStringField(stmtItem, "NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock")),
			NullIfEmpty(GetStringField(stmtItem, "NumberOfTreasuryStockAtTheEndOfFiscalYear")),
			NullIfEmpty(GetStringField(stmtItem, "AverageNumberOfShares")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedNetSales")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedOperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedOrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedProfit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NonConsolidatedEarningsPerShare")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedTotalAssets")),
			NullIfEmpty(GetStringField(stmtItem, "NonConsolidatedEquity")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NonConsolidatedEquityToAssetRatio")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NonConsolidatedBookValuePerShare")),
			// 非連結予想データ（新規追加分）
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedNetSales2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedOperatingProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedOrdinaryProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedProfit2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastNonConsolidatedEarningsPerShare2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedNetSales2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedOperatingProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedOrdinaryProfit2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedProfit2NdQuarter")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastNonConsolidatedEarningsPerShare2NdQuarter")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedNetSales")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedOperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedOrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "ForecastNonConsolidatedProfit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "ForecastNonConsolidatedEarningsPerShare")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedNetSales")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedOperatingProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedOrdinaryProfit")),
			NullIfEmpty(GetStringField(stmtItem, "NextYearForecastNonConsolidatedProfit")),
			NullIfEmptyFloat(GetStringField(stmtItem, "NextYearForecastNonConsolidatedEarningsPerShare")),
		)

		if err != nil {
			// local_codeの外部キー制約エラーの場合はログ出力をスキップ
			if !strings.Contains(err.Error(), "fk_statements_local_code") {
				log.Printf("データ挿入エラー (コード: %s): %v", localCode, err)
			}
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			insertedCount++
		} else {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	totalCount := insertedCount + updatedCount
	return totalCount, nil
}
