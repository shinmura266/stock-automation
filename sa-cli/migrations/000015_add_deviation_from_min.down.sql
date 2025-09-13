-- assessmentテーブルから3か月最低値乖離率カラムを削除
ALTER TABLE assessment DROP COLUMN deviation_from_min;
