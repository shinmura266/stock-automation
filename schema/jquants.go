package schema

import (
	"time"
)

// リフレッシュトークンリクエスト
// https://api.jquants.com/v1/token/auth_user
type AuthUserRequest struct {
	Mailaddress string `json:"mailaddress"`
	Password    string `json:"password"`
}

// リフレッシュトークンレスポンス (有効期限: 1週間)
// https://api.jquants.com/v1/token/auth_user
type AuthUserResponse struct {
	RefreshToken string `json:"refreshToken"`
}

// IDトークンレスポンス (有効期限: 24時間)
// https://api.jquants.com/v1/token/auth_refresh
type IdTokenResponse struct {
	IdToken string `json:"idToken"`
}

// MarketCode 市場区分コード
type MarketCode struct {
	Code      string    `json:"Code" gorm:"column:code;primaryKey"`
	Name      string    `json:"Name" gorm:"column:name"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (MarketCode) TableName() string {
	return "market_codes"
}

// Sector17Code 17業種コード
type Sector17Code struct {
	Code      string    `json:"Code" gorm:"column:code;primaryKey"`
	Name      string    `json:"Name" gorm:"column:name"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (Sector17Code) TableName() string {
	return "sector17_codes"
}

// Sector33Code 33業種コード
type Sector33Code struct {
	Code      string    `json:"Code" gorm:"column:code;primaryKey"`
	Name      string    `json:"Name" gorm:"column:name"`
	CreatedAt time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (Sector33Code) TableName() string {
	return "sector33_codes"
}

// ListedInfoResponse 上場銘柄一覧レスポンス
// https://api.jquants.com/v1/listed/info
type ListedInfoResponse struct {
	Info          []ListedInfo `json:"info"`
	PaginationKey string       `json:"pagination_key"`
}

// ListedInfo 上場銘柄情報
type ListedInfo struct {
	Date               string    `json:"Date" gorm:"column:effective_date"`
	Code               string    `json:"Code" gorm:"column:code;primaryKey"`
	CompanyName        string    `json:"CompanyName" gorm:"column:company_name"`
	CompanyNameEnglish string    `json:"CompanyNameEnglish" gorm:"column:company_name_english"`
	Sector17Code       string    `json:"Sector17Code" gorm:"column:sector17_code"`
	Sector17CodeName   string    `json:"Sector17CodeName"`
	Sector33Code       string    `json:"Sector33Code" gorm:"column:sector33_code"`
	Sector33CodeName   string    `json:"Sector33CodeName"`
	ScaleCategory      string    `json:"ScaleCategory" gorm:"column:scale_category"`
	MarketCode         string    `json:"MarketCode" gorm:"column:market_code"`
	MarketCodeName     string    `json:"MarketCodeName"`
	MarginCode         string    `json:"MarginCode" gorm:"column:margin_code"`
	MarginCodeName     string    `json:"MarginCodeName" gorm:"column:margin_code_name"`
	CreatedAt          time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (ListedInfo) TableName() string {
	return "listed_info"
}

// DailyQuotesResponse 株価四本値レスポンス
// https://api.jquants.com/v1/prices/daily_quotes
type DailyQuotesResponse struct {
	DailyQuotes   []DailyQuote `json:"daily_quotes"`
	PaginationKey string       `json:"pagination_key"`
}

// DailyQuote 四本値1レコード
type DailyQuote struct {
	Date             string    `json:"Date" gorm:"column:trade_date;primaryKey"`
	Code             string    `json:"Code" gorm:"column:code;primaryKey"`
	Open             float64   `json:"Open" gorm:"column:open"`
	High             float64   `json:"High" gorm:"column:high"`
	Low              float64   `json:"Low" gorm:"column:low"`
	Close            float64   `json:"Close" gorm:"column:close"`
	UpperLimit       string    `json:"UpperLimit" gorm:"column:upper_limit"`
	LowerLimit       string    `json:"LowerLimit" gorm:"column:lower_limit"`
	Volume           float64   `json:"Volume" gorm:"column:volume"`
	TurnoverValue    float64   `json:"TurnoverValue" gorm:"column:turnover_value"`
	AdjustmentFactor float64   `json:"AdjustmentFactor" gorm:"column:adjustment_factor"`
	AdjustmentOpen   float64   `json:"AdjustmentOpen" gorm:"column:adjustment_open"`
	AdjustmentHigh   float64   `json:"AdjustmentHigh" gorm:"column:adjustment_high"`
	AdjustmentLow    float64   `json:"AdjustmentLow" gorm:"column:adjustment_low"`
	AdjustmentClose  float64   `json:"AdjustmentClose" gorm:"column:adjustment_close"`
	AdjustmentVolume float64   `json:"AdjustmentVolume" gorm:"column:adjustment_volume"`
	CreatedAt        time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (DailyQuote) TableName() string {
	return "daily_quotes"
}

// FinancialStatementsResponse 財務情報レスポンス
// https://api.jquants.com/v1/fins/statements
type FinancialStatementsResponse struct {
	Statements    []FinancialStatement `json:"statements"`
	PaginationKey string               `json:"pagination_key"`
}

// FinancialStatement 財務情報1レコード
type FinancialStatement struct {
	DisclosedDate                                                                string `json:"DisclosedDate" gorm:"column:disclosed_date;primaryKey"`
	DisclosedTime                                                                string `json:"DisclosedTime" gorm:"column:disclosed_time"`
	LocalCode                                                                    string `json:"LocalCode" gorm:"column:local_code;primaryKey"`
	DisclosureNumber                                                             string `json:"DisclosureNumber" gorm:"column:disclosure_number"`
	TypeOfDocument                                                               string `json:"TypeOfDocument" gorm:"column:type_of_document"`
	TypeOfCurrentPeriod                                                          string `json:"TypeOfCurrentPeriod" gorm:"column:type_of_current_period;primaryKey"`
	CurrentPeriodStartDate                                                       string `json:"CurrentPeriodStartDate" gorm:"column:current_period_start_date"`
	CurrentPeriodEndDate                                                         string `json:"CurrentPeriodEndDate" gorm:"column:current_period_end_date"`
	CurrentFiscalYearStartDate                                                   string `json:"CurrentFiscalYearStartDate" gorm:"column:current_fiscal_year_start_date"`
	CurrentFiscalYearEndDate                                                     string `json:"CurrentFiscalYearEndDate" gorm:"column:current_fiscal_year_end_date"`
	NextFiscalYearStartDate                                                      string `json:"NextFiscalYearStartDate" gorm:"column:next_fiscal_year_start_date"`
	NextFiscalYearEndDate                                                        string `json:"NextFiscalYearEndDate" gorm:"column:next_fiscal_year_end_date"`
	NetSales                                                                     string `json:"NetSales" gorm:"column:net_sales"`
	OperatingProfit                                                              string `json:"OperatingProfit" gorm:"column:operating_profit"`
	OrdinaryProfit                                                               string `json:"OrdinaryProfit" gorm:"column:ordinary_profit"`
	Profit                                                                       string `json:"Profit" gorm:"column:profit"`
	EarningsPerShare                                                             string `json:"EarningsPerShare" gorm:"column:eps"`
	DilutedEarningsPerShare                                                      string `json:"DilutedEarningsPerShare" gorm:"column:diluted_eps"`
	TotalAssets                                                                  string `json:"TotalAssets" gorm:"column:total_assets"`
	Equity                                                                       string `json:"Equity" gorm:"column:equity"`
	EquityToAssetRatio                                                           string `json:"EquityToAssetRatio" gorm:"column:equity_to_asset_ratio"`
	BookValuePerShare                                                            string `json:"BookValuePerShare" gorm:"column:bvps"`
	CashFlowsFromOperatingActivities                                             string `json:"CashFlowsFromOperatingActivities" gorm:"column:cf_operating"`
	CashFlowsFromInvestingActivities                                             string `json:"CashFlowsFromInvestingActivities" gorm:"column:cf_investing"`
	CashFlowsFromFinancingActivities                                             string `json:"CashFlowsFromFinancingActivities" gorm:"column:cf_financing"`
	CashAndEquivalents                                                           string `json:"CashAndEquivalents" gorm:"column:cash_and_equivalents"`
	ResultDividendPerShare1StQuarter                                             string `json:"ResultDividendPerShare1stQuarter" gorm:"column:result_dps_1q"`
	ResultDividendPerShare2NdQuarter                                             string `json:"ResultDividendPerShare2ndQuarter" gorm:"column:result_dps_2q"`
	ResultDividendPerShare3RdQuarter                                             string `json:"ResultDividendPerShare3rdQuarter" gorm:"column:result_dps_3q"`
	ResultDividendPerShareFY                                                     string `json:"ResultDividendPerShareFiscalYearEnd" gorm:"column:result_dps_fy"`
	ResultDividendPerShareAnnual                                                 string `json:"ResultDividendPerShareAnnual" gorm:"column:result_dps_annual"`
	DistributionsPerUnitREIT                                                     string `json:"DistributionsPerUnit(REIT)" gorm:"column:distributions_per_unit_reit"`
	ResultTotalDividendPaidAnnual                                                string `json:"ResultTotalDividendPaidAnnual" gorm:"column:result_total_dividend_annual"`
	ResultPayoutRatioAnnual                                                      string `json:"ResultPayoutRatioAnnual" gorm:"column:result_payout_ratio_annual"`
	ForecastDividendPerShare1StQuarter                                           string `json:"ForecastDividendPerShare1stQuarter" gorm:"column:fc_dps_1q"`
	ForecastDividendPerShare2NdQuarter                                           string `json:"ForecastDividendPerShare2ndQuarter" gorm:"column:fc_dps_2q"`
	ForecastDividendPerShare3RdQuarter                                           string `json:"ForecastDividendPerShare3rdQuarter" gorm:"column:fc_dps_3q"`
	ForecastDividendPerShareFY                                                   string `json:"ForecastDividendPerShareFiscalYearEnd" gorm:"column:fc_dps_fy"`
	ForecastDividendPerShareAnnual                                               string `json:"ForecastDividendPerShareAnnual" gorm:"column:fc_dps_annual"`
	ForecastDistributionsPerUnitREIT                                             string `json:"ForecastDistributionsPerUnit(REIT)" gorm:"column:fc_distributions_per_unit_reit"`
	ForecastTotalDividendPaidAnnual                                              string `json:"ForecastTotalDividendPaidAnnual" gorm:"column:fc_total_dividend_annual"`
	ForecastPayoutRatioAnnual                                                    string `json:"ForecastPayoutRatioAnnual" gorm:"column:fc_payout_ratio_annual"`
	NextYearForecastDividendPerShare1StQuarter                                   string `json:"NextYearForecastDividendPerShare1stQuarter" gorm:"column:ny_fc_dps_1q"`
	NextYearForecastDividendPerShare2NdQuarter                                   string `json:"NextYearForecastDividendPerShare2ndQuarter" gorm:"column:ny_fc_dps_2q"`
	NextYearForecastDividendPerShare3RdQuarter                                   string `json:"NextYearForecastDividendPerShare3rdQuarter" gorm:"column:ny_fc_dps_3q"`
	NextYearForecastDividendPerShareFY                                           string `json:"NextYearForecastDividendPerShareFY" gorm:"column:ny_fc_dps_fy"`
	NextYearForecastDistributionsPerUnitREIT                                     string `json:"NextYearForecastDistributionsPerUnit(REIT)" gorm:"column:ny_fc_distributions_per_unit_reit"`
	NextYearForecastPayoutRatioAnnual                                            string `json:"NextYearForecastPayoutRatioAnnual" gorm:"column:ny_fc_payout_ratio_annual"`
	ForecastNetSales2NdQuarter                                                   string `json:"ForecastNetSales2ndQuarter" gorm:"column:fc_net_sales_2q"`
	ForecastOperatingProfit2NdQuarter                                            string `json:"ForecastOperatingProfit2ndQuarter" gorm:"column:fc_operating_profit_2q"`
	ForecastOrdinaryProfit2NdQuarter                                             string `json:"ForecastOrdinaryProfit2ndQuarter" gorm:"column:fc_ordinary_profit_2q"`
	ForecastProfit2NdQuarter                                                     string `json:"ForecastProfit2ndQuarter" gorm:"column:fc_profit_2q"`
	ForecastEarningsPerShare2NdQuarter                                           string `json:"ForecastEarningsPerShare2ndQuarter" gorm:"column:fc_eps_2q"`
	NextYearForecastNetSales2NdQuarter                                           string `json:"NextYearForecastNetSales2ndQuarter" gorm:"column:ny_fc_net_sales_2q"`
	NextYearForecastOperatingProfit2NdQuarter                                    string `json:"NextYearForecastOperatingProfit2ndQuarter" gorm:"column:ny_fc_operating_profit_2q"`
	NextYearForecastOrdinaryProfit2NdQuarter                                     string `json:"NextYearForecastOrdinaryProfit2ndQuarter" gorm:"column:ny_fc_ordinary_profit_2q"`
	NextYearForecastProfit2NdQuarter                                             string `json:"NextYearForecastProfit2NdQuarter" gorm:"column:ny_fc_profit_2q"`
	NextYearForecastEarningsPerShare2NdQuarter                                   string `json:"NextYearForecastEarningsPerShare2NdQuarter" gorm:"column:ny_fc_eps_2q"`
	ForecastNetSales                                                             string `json:"ForecastNetSales" gorm:"column:fc_net_sales"`
	ForecastOperatingProfit                                                      string `json:"ForecastOperatingProfit" gorm:"column:fc_operating_profit"`
	ForecastOrdinaryProfit                                                       string `json:"ForecastOrdinaryProfit" gorm:"column:fc_ordinary_profit"`
	ForecastProfit                                                               string `json:"ForecastProfit" gorm:"column:fc_profit"`
	ForecastEarningsPerShare                                                     string `json:"ForecastEarningsPerShare" gorm:"column:fc_eps"`
	NextYearForecastNetSales                                                     string `json:"NextYearForecastNetSales" gorm:"column:ny_fc_net_sales"`
	NextYearForecastOperatingProfit                                              string `json:"NextYearForecastOperatingProfit" gorm:"column:ny_fc_operating_profit"`
	NextYearForecastOrdinaryProfit                                               string `json:"NextYearForecastOrdinaryProfit" gorm:"column:ny_fc_ordinary_profit"`
	NextYearForecastProfit                                                       string `json:"NextYearForecastProfit" gorm:"column:ny_fc_profit"`
	NextYearForecastEarningsPerShare                                             string `json:"NextYearForecastEarningsPerShare" gorm:"column:ny_fc_eps"`
	MaterialChangesInSubsidiaries                                                string `json:"MaterialChangesInSubsidiaries" gorm:"column:material_changes_subsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  string `json:"SignificantChangesInTheScopeOfConsolidation" gorm:"column:significant_changes_consolidation_scope"`
	ChangesBasedOnRevisionsOfAccountingStandard                                  string `json:"ChangesBasedOnRevisionsOfAccountingStandard" gorm:"column:changes_accounting_std_revisions"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     string `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard" gorm:"column:changes_accounting_std_other"`
	ChangesInAccountingEstimates                                                 string `json:"ChangesInAccountingEstimates" gorm:"column:changes_accounting_estimates"`
	RetrospectiveRestatement                                                     string `json:"RetrospectiveRestatement" gorm:"column:retrospective_restatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock string `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock" gorm:"column:issued_shares_end_fy_incl_treasury"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    string `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear" gorm:"column:treasury_shares_end_fy"`
	AverageNumberOfShares                                                        string `json:"AverageNumberOfShares" gorm:"column:avg_shares"`
	NonConsolidatedNetSales                                                      string `json:"NonConsolidatedNetSales" gorm:"column:nc_net_sales"`
	NonConsolidatedOperatingProfit                                               string `json:"NonConsolidatedOperatingProfit" gorm:"column:nc_operating_profit"`
	NonConsolidatedOrdinaryProfit                                                string `json:"NonConsolidatedOrdinaryProfit" gorm:"column:nc_ordinary_profit"`
	NonConsolidatedProfit                                                        string `json:"NonConsolidatedProfit" gorm:"column:nc_profit"`
	NonConsolidatedEarningsPerShare                                              string `json:"NonConsolidatedEarningsPerShare" gorm:"column:nc_eps"`
	NonConsolidatedTotalAssets                                                   string `json:"NonConsolidatedTotalAssets" gorm:"column:nc_total_assets"`
	NonConsolidatedEquity                                                        string `json:"NonConsolidatedEquity" gorm:"column:nc_equity"`
	NonConsolidatedEquityToAssetRatio                                            string `json:"NonConsolidatedEquityToAssetRatio" gorm:"column:nc_equity_to_asset_ratio"`
	NonConsolidatedBookValuePerShare                                             string `json:"NonConsolidatedBookValuePerShare" gorm:"column:nc_bvps"`
	// 非連結予想データ
	ForecastNonConsolidatedNetSales2NdQuarter                 string    `json:"ForecastNonConsolidatedNetSales2ndQuarter" gorm:"column:fc_nc_net_sales_2q"`
	ForecastNonConsolidatedOperatingProfit2NdQuarter          string    `json:"ForecastNonConsolidatedOperatingProfit2ndQuarter" gorm:"column:fc_nc_operating_profit_2q"`
	ForecastNonConsolidatedOrdinaryProfit2NdQuarter           string    `json:"ForecastNonConsolidatedOrdinaryProfit2ndQuarter" gorm:"column:fc_nc_ordinary_profit_2q"`
	ForecastNonConsolidatedProfit2NdQuarter                   string    `json:"ForecastNonConsolidatedProfit2ndQuarter" gorm:"column:fc_nc_profit_2q"`
	ForecastNonConsolidatedEarningsPerShare2NdQuarter         string    `json:"ForecastNonConsolidatedEarningsPerShare2ndQuarter" gorm:"column:fc_nc_eps_2q"`
	NextYearForecastNonConsolidatedNetSales2NdQuarter         string    `json:"NextYearForecastNonConsolidatedNetSales2ndQuarter" gorm:"column:ny_fc_nc_net_sales_2q"`
	NextYearForecastNonConsolidatedOperatingProfit2NdQuarter  string    `json:"NextYearForecastNonConsolidatedOperatingProfit2NdQuarter" gorm:"column:ny_fc_nc_operating_profit_2q"`
	NextYearForecastNonConsolidatedOrdinaryProfit2NdQuarter   string    `json:"NextYearForecastNonConsolidatedOrdinaryProfit2NdQuarter" gorm:"column:ny_fc_nc_ordinary_profit_2q"`
	NextYearForecastNonConsolidatedProfit2NdQuarter           string    `json:"NextYearForecastNonConsolidatedProfit2NdQuarter" gorm:"column:ny_fc_nc_profit_2q"`
	NextYearForecastNonConsolidatedEarningsPerShare2NdQuarter string    `json:"NextYearForecastNonConsolidatedEarningsPerShare2NdQuarter" gorm:"column:ny_fc_nc_eps_2q"`
	ForecastNonConsolidatedNetSales                           string    `json:"ForecastNonConsolidatedNetSales" gorm:"column:fc_nc_net_sales"`
	ForecastNonConsolidatedOperatingProfit                    string    `json:"ForecastNonConsolidatedOperatingProfit" gorm:"column:fc_nc_operating_profit"`
	ForecastNonConsolidatedOrdinaryProfit                     string    `json:"ForecastNonConsolidatedOrdinaryProfit" gorm:"column:fc_nc_ordinary_profit"`
	ForecastNonConsolidatedProfit                             string    `json:"ForecastNonConsolidatedProfit" gorm:"column:fc_nc_profit"`
	ForecastNonConsolidatedEarningsPerShare                   string    `json:"ForecastNonConsolidatedEarningsPerShare" gorm:"column:fc_nc_eps"`
	NextYearForecastNonConsolidatedNetSales                   string    `json:"NextYearForecastNonConsolidatedNetSales" gorm:"column:ny_fc_nc_net_sales"`
	NextYearForecastNonConsolidatedOperatingProfit            string    `json:"NextYearForecastNonConsolidatedOperatingProfit" gorm:"column:ny_fc_nc_operating_profit"`
	NextYearForecastNonConsolidatedOrdinaryProfit             string    `json:"NextYearForecastNonConsolidatedOrdinaryProfit" gorm:"column:ny_fc_nc_ordinary_profit"`
	NextYearForecastNonConsolidatedProfit                     string    `json:"NextYearForecastNonConsolidatedProfit" gorm:"column:ny_fc_nc_profit"`
	NextYearForecastNonConsolidatedEarningsPerShare           string    `json:"NextYearForecastNonConsolidatedEarningsPerShare" gorm:"column:ny_fc_nc_eps"`
	CreatedAt                                                 time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt                                                 time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (FinancialStatement) TableName() string {
	return "statements"
}
