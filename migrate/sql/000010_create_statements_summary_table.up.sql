-- 財務情報サマリーテーブルを作成
-- local_code毎に各会計年度の最新の財務情報を管理
CREATE TABLE IF NOT EXISTS statements_summary (
    local_code VARCHAR(10) NOT NULL,
    fiscal_year_start_date DATE NOT NULL,
    fiscal_year_end_date DATE NOT NULL,
    disclosed_date DATE NOT NULL,
    disclosed_time TIME,
    type_of_current_period VARCHAR(10),
    
    -- 売上高・利益系指標
    net_sales BIGINT,
    operating_profit BIGINT,
    ordinary_profit BIGINT,
    profit BIGINT,
    eps DECIMAL(10,2),
    
    -- バランスシート系指標
    total_assets BIGINT,
    equity BIGINT,
    equity_to_asset_ratio DECIMAL(5,2),
    
    -- 配当関連指標
    dividend_per_share DECIMAL(10,2),
    
    -- 実績データか予想データかの判別フラグ
    is_forecast BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- データの種別（当期実績、当期予想、翌期予想）
    data_type ENUM('current_actual', 'current_forecast', 'next_year_forecast') NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- プライマリキー：local_code + fiscal_year_end_date
    PRIMARY KEY (local_code, fiscal_year_end_date),
    
    -- 外部キー制約
    CONSTRAINT fk_statements_summary_local_code FOREIGN KEY (local_code) REFERENCES listed_info(code),
    
    -- インデックス
    INDEX idx_statements_summary_disclosed_date (disclosed_date),
    INDEX idx_statements_summary_fiscal_year (fiscal_year_end_date),
    INDEX idx_statements_summary_code_fiscal (local_code, fiscal_year_end_date)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
