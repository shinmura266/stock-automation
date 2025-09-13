-- assessmentテーブルに3か月最高調整終値カラムを追加
ALTER TABLE assessment ADD COLUMN three_month_max_close DECIMAL(10,2) COMMENT '3か月間の最高調整終値';
