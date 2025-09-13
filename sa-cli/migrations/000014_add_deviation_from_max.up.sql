-- assessmentテーブルに3か月最高値乖離率カラムを追加
ALTER TABLE assessment ADD COLUMN deviation_from_max DECIMAL(7,2) COMMENT '3か月最高値からの乖離率(%)';
