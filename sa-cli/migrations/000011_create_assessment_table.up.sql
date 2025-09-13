-- 銘柄評価テーブルを作成
-- 各銘柄の包括的な評価指標を管理
CREATE TABLE IF NOT EXISTS assessment (
    code VARCHAR(10) NOT NULL,
    last_fiscal_year_end_date DATE,
    last_dividend_per_share DECIMAL(10,2),
    last_trade_date DATE,
    last_adjustment_close DECIMAL(10,2),
    last_dividend_yield DECIMAL(5,2),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- プライマリキー
    PRIMARY KEY (code),
    
    -- 外部キー制約
    CONSTRAINT fk_assessment_code FOREIGN KEY (code) REFERENCES listed_info(code),
    
    -- インデックス
    INDEX idx_assessment_fiscal_year (last_fiscal_year_end_date),
    INDEX idx_assessment_trade_date (last_trade_date),
    INDEX idx_assessment_dividend_yield (last_dividend_yield)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
