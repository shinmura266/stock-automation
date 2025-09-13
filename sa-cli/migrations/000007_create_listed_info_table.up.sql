-- 上場銘柄情報テーブルを作成
CREATE TABLE IF NOT EXISTS listed_info (
    code VARCHAR(10) NOT NULL PRIMARY KEY,
    effective_date DATE NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    company_name_english VARCHAR(255),
    market_code VARCHAR(10),
    sector17_code VARCHAR(10),
    sector33_code VARCHAR(10),
    scale_category VARCHAR(100),
    margin_code VARCHAR(10),
    margin_code_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_listed_info_market_code FOREIGN KEY (market_code) REFERENCES market_codes(code),
    CONSTRAINT fk_listed_info_sector17_code FOREIGN KEY (sector17_code) REFERENCES sector17_codes(code),
    CONSTRAINT fk_listed_info_sector33_code FOREIGN KEY (sector33_code) REFERENCES sector33_codes(code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
