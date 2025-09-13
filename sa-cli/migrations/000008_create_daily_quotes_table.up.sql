CREATE TABLE IF NOT EXISTS daily_quotes (
    trade_date DATE NOT NULL,
    code VARCHAR(10) NOT NULL,
    open DECIMAL(10,2) NULL,
    high DECIMAL(10,2) NULL,
    low DECIMAL(10,2) NULL,
    close DECIMAL(10,2) NULL,
    upper_limit VARCHAR(1) NULL,
    lower_limit VARCHAR(1) NULL,
    volume DECIMAL(15,0) NULL,
    turnover_value DECIMAL(20,0) NULL,
    adjustment_factor DECIMAL(10,6) NULL,
    adjustment_open DECIMAL(10,2) NULL,
    adjustment_high DECIMAL(10,2) NULL,
    adjustment_low DECIMAL(10,2) NULL,
    adjustment_close DECIMAL(10,2) NULL,
    adjustment_volume DECIMAL(15,0) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (trade_date, code)
);

