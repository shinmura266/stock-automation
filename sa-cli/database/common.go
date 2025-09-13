package database

import (
	"database/sql"
	"strconv"
)

// nullIfEmpty 空文字列の場合はnilを返すヘルパー関数（整数値用）
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return nil
}

// nullIfEmptyFloat 空文字列の場合はnilを返すヘルパー関数（小数値用）
func nullIfEmptyFloat(s string) interface{} {
	if s == "" {
		return nil
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return nil
}

// beginTransaction トランザクションを開始し、パニック時の自動ロールバックを設定
func beginTransaction(db *sql.DB) (*sql.Tx, func()) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}

	return tx, cleanup
}
