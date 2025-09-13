-- assessmentテーブルから3か月最高値乖離率カラムを削除
ALTER TABLE assessment DROP COLUMN deviation_from_max;
