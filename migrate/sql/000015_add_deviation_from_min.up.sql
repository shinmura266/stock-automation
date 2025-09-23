-- assessmentテーブルに3か月最低値乖離率カラムを追加
ALTER TABLE assessment ADD COLUMN deviation_from_min DECIMAL(7,2) COMMENT '3か月最低値からの乖離率(%)';
