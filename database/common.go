package database

import (
	"database/sql"
	"strconv"
)

// NullIfEmpty 空文字列の場合はnilを返すヘルパー関数（整数値用）
func NullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return nil
}

// NullIfEmptyFloat 空文字列の場合はnilを返すヘルパー関数（小数値用）
func NullIfEmptyFloat(s string) interface{} {
	if s == "" {
		return nil
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return nil
}

// BeginTransaction トランザクションを開始し、パニック時の自動ロールバックを設定
func BeginTransaction(db *sql.DB) (*sql.Tx, func()) {
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
