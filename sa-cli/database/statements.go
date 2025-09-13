package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"kabu-analysis/jquants"
)

// StatementsRepository 財務情報のリポジトリ
type StatementsRepository struct {
	conn *Connection
}

// NewStatementsRepository 新しいリポジトリを作成
func NewStatementsRepository(conn *Connection) *StatementsRepository {
	return &StatementsRepository{conn: conn}
}

// SaveStatements 財務情報をデータベースに保存
func (r *StatementsRepository) SaveStatements(resp *jquants.FinancialStatementsResponse) (int, error) {
	if len(resp.Statements) == 0 {
		return 0, fmt.Errorf("保存するデータがありません")
	}

	tx, cleanup := beginTransaction(r.conn.GetDB())
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

	for _, stmt_data := range resp.Statements {
		// 日付をパース
		disclosedDate, err := time.Parse("2006-01-02", stmt_data.DisclosedDate)
		if err != nil {
			log.Printf("開示日パースエラー (コード: %s, 日付: %s): %v", stmt_data.LocalCode, stmt_data.DisclosedDate, err)
			continue
		}

		// 時刻をパース（空の場合はnilを設定）
		var disclosedTime interface{}
		if stmt_data.DisclosedTime != "" {
			if t, err := time.Parse("15:04:05", stmt_data.DisclosedTime); err == nil {
				disclosedTime = t.Format("15:04:05")
			}
		}

		// 期間開始日・終了日をパース
		var currentPeriodStartDate, currentPeriodEndDate interface{}
		var currentFiscalYearStartDate, currentFiscalYearEndDate interface{}
		var nextFiscalYearStartDate, nextFiscalYearEndDate interface{}

		if stmt_data.CurrentPeriodStartDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.CurrentPeriodStartDate); err == nil {
				currentPeriodStartDate = d
			}
		}
		if stmt_data.CurrentPeriodEndDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.CurrentPeriodEndDate); err == nil {
				currentPeriodEndDate = d
			}
		}
		if stmt_data.CurrentFiscalYearStartDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.CurrentFiscalYearStartDate); err == nil {
				currentFiscalYearStartDate = d
			}
		}
		if stmt_data.CurrentFiscalYearEndDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.CurrentFiscalYearEndDate); err == nil {
				currentFiscalYearEndDate = d
			}
		}
		if stmt_data.NextFiscalYearStartDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.NextFiscalYearStartDate); err == nil {
				nextFiscalYearStartDate = d
			}
		}
		if stmt_data.NextFiscalYearEndDate != "" {
			if d, err := time.Parse("2006-01-02", stmt_data.NextFiscalYearEndDate); err == nil {
				nextFiscalYearEndDate = d
			}
		}

		result, err := stmt.Exec(
			disclosedDate, disclosedTime, stmt_data.LocalCode, stmt_data.DisclosureNumber, stmt_data.TypeOfDocument,
			stmt_data.TypeOfCurrentPeriod, currentPeriodStartDate, currentPeriodEndDate,
			currentFiscalYearStartDate, currentFiscalYearEndDate,
			nextFiscalYearStartDate, nextFiscalYearEndDate,
			nullIfEmpty(stmt_data.NetSales), nullIfEmpty(stmt_data.OperatingProfit),
			nullIfEmpty(stmt_data.OrdinaryProfit), nullIfEmpty(stmt_data.Profit),
			nullIfEmptyFloat(stmt_data.EarningsPerShare), nullIfEmptyFloat(stmt_data.DilutedEarningsPerShare),
			nullIfEmpty(stmt_data.TotalAssets), nullIfEmpty(stmt_data.Equity),
			nullIfEmptyFloat(stmt_data.EquityToAssetRatio), nullIfEmptyFloat(stmt_data.BookValuePerShare),
			nullIfEmpty(stmt_data.CashFlowsFromOperatingActivities), nullIfEmpty(stmt_data.CashFlowsFromInvestingActivities),
			nullIfEmpty(stmt_data.CashFlowsFromFinancingActivities), nullIfEmpty(stmt_data.CashAndEquivalents),
			nullIfEmptyFloat(stmt_data.ResultDividendPerShare1StQuarter), nullIfEmptyFloat(stmt_data.ResultDividendPerShare2NdQuarter),
			nullIfEmptyFloat(stmt_data.ResultDividendPerShare3RdQuarter), nullIfEmptyFloat(stmt_data.ResultDividendPerShareFY),
			nullIfEmptyFloat(stmt_data.ResultDividendPerShareAnnual),
			nullIfEmptyFloat(stmt_data.DistributionsPerUnitREIT), nullIfEmpty(stmt_data.ResultTotalDividendPaidAnnual),
			nullIfEmptyFloat(stmt_data.ResultPayoutRatioAnnual),
			nullIfEmptyFloat(stmt_data.ForecastDividendPerShare1StQuarter), nullIfEmptyFloat(stmt_data.ForecastDividendPerShare2NdQuarter),
			nullIfEmptyFloat(stmt_data.ForecastDividendPerShare3RdQuarter), nullIfEmptyFloat(stmt_data.ForecastDividendPerShareFY),
			nullIfEmptyFloat(stmt_data.ForecastDividendPerShareAnnual),
			nullIfEmptyFloat(stmt_data.ForecastDistributionsPerUnitREIT), nullIfEmpty(stmt_data.ForecastTotalDividendPaidAnnual),
			nullIfEmptyFloat(stmt_data.ForecastPayoutRatioAnnual),
			nullIfEmptyFloat(stmt_data.NextYearForecastDividendPerShare1StQuarter), nullIfEmptyFloat(stmt_data.NextYearForecastDividendPerShare2NdQuarter),
			nullIfEmptyFloat(stmt_data.NextYearForecastDividendPerShare3RdQuarter), nullIfEmptyFloat(stmt_data.NextYearForecastDividendPerShareFY),
			nullIfEmptyFloat(stmt_data.NextYearForecastDistributionsPerUnitREIT), nullIfEmptyFloat(stmt_data.NextYearForecastPayoutRatioAnnual),
			nullIfEmpty(stmt_data.ForecastNetSales2NdQuarter), nullIfEmpty(stmt_data.ForecastOperatingProfit2NdQuarter),
			nullIfEmpty(stmt_data.ForecastOrdinaryProfit2NdQuarter), nullIfEmpty(stmt_data.ForecastProfit2NdQuarter),
			nullIfEmptyFloat(stmt_data.ForecastEarningsPerShare2NdQuarter),
			nullIfEmpty(stmt_data.NextYearForecastNetSales2NdQuarter), nullIfEmpty(stmt_data.NextYearForecastOperatingProfit2NdQuarter),
			nullIfEmpty(stmt_data.NextYearForecastOrdinaryProfit2NdQuarter), nullIfEmpty(stmt_data.NextYearForecastProfit2NdQuarter),
			nullIfEmptyFloat(stmt_data.NextYearForecastEarningsPerShare2NdQuarter),
			nullIfEmpty(stmt_data.ForecastNetSales), nullIfEmpty(stmt_data.ForecastOperatingProfit),
			nullIfEmpty(stmt_data.ForecastOrdinaryProfit), nullIfEmpty(stmt_data.ForecastProfit),
			nullIfEmptyFloat(stmt_data.ForecastEarningsPerShare),
			nullIfEmpty(stmt_data.NextYearForecastNetSales), nullIfEmpty(stmt_data.NextYearForecastOperatingProfit),
			nullIfEmpty(stmt_data.NextYearForecastOrdinaryProfit), nullIfEmpty(stmt_data.NextYearForecastProfit),
			nullIfEmptyFloat(stmt_data.NextYearForecastEarningsPerShare),
			stmt_data.MaterialChangesInSubsidiaries, stmt_data.SignificantChangesInTheScopeOfConsolidation,
			stmt_data.ChangesBasedOnRevisionsOfAccountingStandard, stmt_data.ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard,
			stmt_data.ChangesInAccountingEstimates, stmt_data.RetrospectiveRestatement,
			nullIfEmpty(stmt_data.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock),
			nullIfEmpty(stmt_data.NumberOfTreasuryStockAtTheEndOfFiscalYear), nullIfEmpty(stmt_data.AverageNumberOfShares),
			nullIfEmpty(stmt_data.NonConsolidatedNetSales), nullIfEmpty(stmt_data.NonConsolidatedOperatingProfit),
			nullIfEmpty(stmt_data.NonConsolidatedOrdinaryProfit), nullIfEmpty(stmt_data.NonConsolidatedProfit),
			nullIfEmptyFloat(stmt_data.NonConsolidatedEarningsPerShare), nullIfEmpty(stmt_data.NonConsolidatedTotalAssets),
			nullIfEmpty(stmt_data.NonConsolidatedEquity), nullIfEmptyFloat(stmt_data.NonConsolidatedEquityToAssetRatio),
			nullIfEmptyFloat(stmt_data.NonConsolidatedBookValuePerShare),
			// 非連結予想データ（新規追加分）
			nullIfEmpty(stmt_data.ForecastNonConsolidatedNetSales2NdQuarter), nullIfEmpty(stmt_data.ForecastNonConsolidatedOperatingProfit2NdQuarter),
			nullIfEmpty(stmt_data.ForecastNonConsolidatedOrdinaryProfit2NdQuarter), nullIfEmpty(stmt_data.ForecastNonConsolidatedProfit2NdQuarter),
			nullIfEmptyFloat(stmt_data.ForecastNonConsolidatedEarningsPerShare2NdQuarter),
			nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedNetSales2NdQuarter), nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedOperatingProfit2NdQuarter),
			nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedOrdinaryProfit2NdQuarter), nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedProfit2NdQuarter),
			nullIfEmptyFloat(stmt_data.NextYearForecastNonConsolidatedEarningsPerShare2NdQuarter),
			nullIfEmpty(stmt_data.ForecastNonConsolidatedNetSales), nullIfEmpty(stmt_data.ForecastNonConsolidatedOperatingProfit),
			nullIfEmpty(stmt_data.ForecastNonConsolidatedOrdinaryProfit), nullIfEmpty(stmt_data.ForecastNonConsolidatedProfit),
			nullIfEmptyFloat(stmt_data.ForecastNonConsolidatedEarningsPerShare),
			nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedNetSales), nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedOperatingProfit),
			nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedOrdinaryProfit), nullIfEmpty(stmt_data.NextYearForecastNonConsolidatedProfit),
			nullIfEmptyFloat(stmt_data.NextYearForecastNonConsolidatedEarningsPerShare),
		)

		if err != nil {
			// local_codeの外部キー制約エラーの場合はログ出力をスキップ
			if !strings.Contains(err.Error(), "fk_statements_local_code") {
				log.Printf("データ挿入エラー (コード: %s): %v", stmt_data.LocalCode, err)
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
