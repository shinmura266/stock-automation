-- 基本的な市場区分コードを挿入
INSERT INTO market_codes (code, name) VALUES
('0101', '東証一部'),
('0102', '東証二部'),
('0104', 'マザーズ'),
('0105', 'TOKYO PRO MARKET'),
('0106', 'JASDAQ スタンダード'),
('0107', 'JASDAQ グロース'),
('0109', 'その他'),
('0111', 'プライム'),
('0112', 'スタンダード'),
('0113', 'グロース')
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    updated_at = CURRENT_TIMESTAMP;
