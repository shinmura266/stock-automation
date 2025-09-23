package schema

import "time"

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
	Code      string `json:"Code" gorm:"column:code;primaryKey"`
	Name      string `json:"Name" gorm:"column:name"`
	CreatedAt string `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt string `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (MarketCode) TableName() string {
	return "market_codes"
}

// Sector17Code 17業種コード
type Sector17Code struct {
	Code      string `json:"Code" gorm:"column:code;primaryKey"`
	Name      string `json:"Name" gorm:"column:name"`
	CreatedAt string `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt string `json:"UpdatedAt" gorm:"column:updated_at"`
}

// TableName GORMのテーブル名を指定
func (Sector17Code) TableName() string {
	return "sector17_codes"
}

// Sector33Code 33業種コード
type Sector33Code struct {
	Code      string `json:"Code" gorm:"column:code;primaryKey"`
	Name      string `json:"Name" gorm:"column:name"`
	CreatedAt string `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt string `json:"UpdatedAt" gorm:"column:updated_at"`
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
	Date               string `json:"Date" gorm:"column:effective_date"`
	Code               string `json:"Code" gorm:"column:code;primaryKey"`
	CompanyName        string `json:"CompanyName" gorm:"column:company_name"`
	CompanyNameEnglish string `json:"CompanyNameEnglish" gorm:"column:company_name_english"`
	Sector17Code       string `json:"Sector17Code" gorm:"column:sector17_code"`
	Sector17CodeName   string `json:"Sector17CodeName"`
	Sector33Code       string `json:"Sector33Code" gorm:"column:sector33_code"`
	Sector33CodeName   string `json:"Sector33CodeName"`
	ScaleCategory      string `json:"ScaleCategory" gorm:"column:scale_category"`
	MarketCode         string `json:"MarketCode" gorm:"column:market_code"`
	MarketCodeName     string `json:"MarketCodeName"`
	MarginCode         string `json:"MarginCode" gorm:"column:margin_code"`
	MarginCodeName     string `json:"MarginCodeName" gorm:"column:margin_code_name"`
	CreatedAt          string `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt          string `json:"UpdatedAt" gorm:"column:updated_at"`
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
	DisclosedDate                                                                string `json:"DisclosedDate"`
	DisclosedTime                                                                string `json:"DisclosedTime"`
	LocalCode                                                                    string `json:"LocalCode"`
	DisclosureNumber                                                             string `json:"DisclosureNumber"`
	TypeOfDocument                                                               string `json:"TypeOfDocument"`
	TypeOfCurrentPeriod                                                          string `json:"TypeOfCurrentPeriod"`
	CurrentPeriodStartDate                                                       string `json:"CurrentPeriodStartDate"`
	CurrentPeriodEndDate                                                         string `json:"CurrentPeriodEndDate"`
	CurrentFiscalYearStartDate                                                   string `json:"CurrentFiscalYearStartDate"`
	CurrentFiscalYearEndDate                                                     string `json:"CurrentFiscalYearEndDate"`
	NextFiscalYearStartDate                                                      string `json:"NextFiscalYearStartDate"`
	NextFiscalYearEndDate                                                        string `json:"NextFiscalYearEndDate"`
	NetSales                                                                     string `json:"NetSales"`
	OperatingProfit                                                              string `json:"OperatingProfit"`
	OrdinaryProfit                                                               string `json:"OrdinaryProfit"`
	Profit                                                                       string `json:"Profit"`
	EarningsPerShare                                                             string `json:"EarningsPerShare"`
	DilutedEarningsPerShare                                                      string `json:"DilutedEarningsPerShare"`
	TotalAssets                                                                  string `json:"TotalAssets"`
	Equity                                                                       string `json:"Equity"`
	EquityToAssetRatio                                                           string `json:"EquityToAssetRatio"`
	BookValuePerShare                                                            string `json:"BookValuePerShare"`
	CashFlowsFromOperatingActivities                                             string `json:"CashFlowsFromOperatingActivities"`
	CashFlowsFromInvestingActivities                                             string `json:"CashFlowsFromInvestingActivities"`
	CashFlowsFromFinancingActivities                                             string `json:"CashFlowsFromFinancingActivities"`
	CashAndEquivalents                                                           string `json:"CashAndEquivalents"`
	ResultDividendPerShare1StQuarter                                             string `json:"ResultDividendPerShare1stQuarter"`
	ResultDividendPerShare2NdQuarter                                             string `json:"ResultDividendPerShare2ndQuarter"`
	ResultDividendPerShare3RdQuarter                                             string `json:"ResultDividendPerShare3rdQuarter"`
	ResultDividendPerShareFY                                                     string `json:"ResultDividendPerShareFiscalYearEnd"`
	ResultDividendPerShareAnnual                                                 string `json:"ResultDividendPerShareAnnual"`
	DistributionsPerUnitREIT                                                     string `json:"DistributionsPerUnit(REIT)"`
	ResultTotalDividendPaidAnnual                                                string `json:"ResultTotalDividendPaidAnnual"`
	ResultPayoutRatioAnnual                                                      string `json:"ResultPayoutRatioAnnual"`
	ForecastDividendPerShare1StQuarter                                           string `json:"ForecastDividendPerShare1stQuarter"`
	ForecastDividendPerShare2NdQuarter                                           string `json:"ForecastDividendPerShare2ndQuarter"`
	ForecastDividendPerShare3RdQuarter                                           string `json:"ForecastDividendPerShare3rdQuarter"`
	ForecastDividendPerShareFY                                                   string `json:"ForecastDividendPerShareFiscalYearEnd"`
	ForecastDividendPerShareAnnual                                               string `json:"ForecastDividendPerShareAnnual"`
	ForecastDistributionsPerUnitREIT                                             string `json:"ForecastDistributionsPerUnit(REIT)"`
	ForecastTotalDividendPaidAnnual                                              string `json:"ForecastTotalDividendPaidAnnual"`
	ForecastPayoutRatioAnnual                                                    string `json:"ForecastPayoutRatioAnnual"`
	NextYearForecastDividendPerShare1StQuarter                                   string `json:"NextYearForecastDividendPerShare1stQuarter"`
	NextYearForecastDividendPerShare2NdQuarter                                   string `json:"NextYearForecastDividendPerShare2ndQuarter"`
	NextYearForecastDividendPerShare3RdQuarter                                   string `json:"NextYearForecastDividendPerShare3rdQuarter"`
	NextYearForecastDividendPerShareFY                                           string `json:"NextYearForecastDividendPerShareFY"`
	NextYearForecastDistributionsPerUnitREIT                                     string `json:"NextYearForecastDistributionsPerUnit(REIT)"`
	NextYearForecastPayoutRatioAnnual                                            string `json:"NextYearForecastPayoutRatioAnnual"`
	ForecastNetSales2NdQuarter                                                   string `json:"ForecastNetSales2ndQuarter"`
	ForecastOperatingProfit2NdQuarter                                            string `json:"ForecastOperatingProfit2ndQuarter"`
	ForecastOrdinaryProfit2NdQuarter                                             string `json:"ForecastOrdinaryProfit2ndQuarter"`
	ForecastProfit2NdQuarter                                                     string `json:"ForecastProfit2ndQuarter"`
	ForecastEarningsPerShare2NdQuarter                                           string `json:"ForecastEarningsPerShare2ndQuarter"`
	NextYearForecastNetSales2NdQuarter                                           string `json:"NextYearForecastNetSales2ndQuarter"`
	NextYearForecastOperatingProfit2NdQuarter                                    string `json:"NextYearForecastOperatingProfit2ndQuarter"`
	NextYearForecastOrdinaryProfit2NdQuarter                                     string `json:"NextYearForecastOrdinaryProfit2ndQuarter"`
	NextYearForecastProfit2NdQuarter                                             string `json:"NextYearForecastProfit2ndQuarter"`
	NextYearForecastEarningsPerShare2NdQuarter                                   string `json:"NextYearForecastEarningsPerShare2ndQuarter"`
	ForecastNetSales                                                             string `json:"ForecastNetSales"`
	ForecastOperatingProfit                                                      string `json:"ForecastOperatingProfit"`
	ForecastOrdinaryProfit                                                       string `json:"ForecastOrdinaryProfit"`
	ForecastProfit                                                               string `json:"ForecastProfit"`
	ForecastEarningsPerShare                                                     string `json:"ForecastEarningsPerShare"`
	NextYearForecastNetSales                                                     string `json:"NextYearForecastNetSales"`
	NextYearForecastOperatingProfit                                              string `json:"NextYearForecastOperatingProfit"`
	NextYearForecastOrdinaryProfit                                               string `json:"NextYearForecastOrdinaryProfit"`
	NextYearForecastProfit                                                       string `json:"NextYearForecastProfit"`
	NextYearForecastEarningsPerShare                                             string `json:"NextYearForecastEarningsPerShare"`
	MaterialChangesInSubsidiaries                                                string `json:"MaterialChangesInSubsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  string `json:"SignificantChangesInTheScopeOfConsolidation"`
	ChangesBasedOnRevisionsOfAccountingStandard                                  string `json:"ChangesBasedOnRevisionsOfAccountingStandard"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     string `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"`
	ChangesInAccountingEstimates                                                 string `json:"ChangesInAccountingEstimates"`
	RetrospectiveRestatement                                                     string `json:"RetrospectiveRestatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock string `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    string `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear"`
	AverageNumberOfShares                                                        string `json:"AverageNumberOfShares"`
	NonConsolidatedNetSales                                                      string `json:"NonConsolidatedNetSales"`
	NonConsolidatedOperatingProfit                                               string `json:"NonConsolidatedOperatingProfit"`
	NonConsolidatedOrdinaryProfit                                                string `json:"NonConsolidatedOrdinaryProfit"`
	NonConsolidatedProfit                                                        string `json:"NonConsolidatedProfit"`
	NonConsolidatedEarningsPerShare                                              string `json:"NonConsolidatedEarningsPerShare"`
	NonConsolidatedTotalAssets                                                   string `json:"NonConsolidatedTotalAssets"`
	NonConsolidatedEquity                                                        string `json:"NonConsolidatedEquity"`
	NonConsolidatedEquityToAssetRatio                                            string `json:"NonConsolidatedEquityToAssetRatio"`
	NonConsolidatedBookValuePerShare                                             string `json:"NonConsolidatedBookValuePerShare"`
	// 非連結予想データ
	ForecastNonConsolidatedNetSales2NdQuarter                 string `json:"ForecastNonConsolidatedNetSales2ndQuarter"`
	ForecastNonConsolidatedOperatingProfit2NdQuarter          string `json:"ForecastNonConsolidatedOperatingProfit2ndQuarter"`
	ForecastNonConsolidatedOrdinaryProfit2NdQuarter           string `json:"ForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	ForecastNonConsolidatedProfit2NdQuarter                   string `json:"ForecastNonConsolidatedProfit2ndQuarter"`
	ForecastNonConsolidatedEarningsPerShare2NdQuarter         string `json:"ForecastNonConsolidatedEarningsPerShare2ndQuarter"`
	NextYearForecastNonConsolidatedNetSales2NdQuarter         string `json:"NextYearForecastNonConsolidatedNetSales2ndQuarter"`
	NextYearForecastNonConsolidatedOperatingProfit2NdQuarter  string `json:"NextYearForecastNonConsolidatedOperatingProfit2ndQuarter"`
	NextYearForecastNonConsolidatedOrdinaryProfit2NdQuarter   string `json:"NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	NextYearForecastNonConsolidatedProfit2NdQuarter           string `json:"NextYearForecastNonConsolidatedProfit2ndQuarter"`
	NextYearForecastNonConsolidatedEarningsPerShare2NdQuarter string `json:"NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter"`
	ForecastNonConsolidatedNetSales                           string `json:"ForecastNonConsolidatedNetSales"`
	ForecastNonConsolidatedOperatingProfit                    string `json:"ForecastNonConsolidatedOperatingProfit"`
	ForecastNonConsolidatedOrdinaryProfit                     string `json:"ForecastNonConsolidatedOrdinaryProfit"`
	ForecastNonConsolidatedProfit                             string `json:"ForecastNonConsolidatedProfit"`
	ForecastNonConsolidatedEarningsPerShare                   string `json:"ForecastNonConsolidatedEarningsPerShare"`
	NextYearForecastNonConsolidatedNetSales                   string `json:"NextYearForecastNonConsolidatedNetSales"`
	NextYearForecastNonConsolidatedOperatingProfit            string `json:"NextYearForecastNonConsolidatedOperatingProfit"`
	NextYearForecastNonConsolidatedOrdinaryProfit             string `json:"NextYearForecastNonConsolidatedOrdinaryProfit"`
	NextYearForecastNonConsolidatedProfit                     string `json:"NextYearForecastNonConsolidatedProfit"`
	NextYearForecastNonConsolidatedEarningsPerShare           string `json:"NextYearForecastNonConsolidatedEarningsPerShare"`
}
