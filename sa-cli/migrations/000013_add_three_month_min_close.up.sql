-- assessmentテーブルに3か月最低調整終値カラムを追加
ALTER TABLE assessment ADD COLUMN three_month_min_close DECIMAL(10,2) COMMENT '3か月間の最低調整終値';
