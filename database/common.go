package database

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"
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

// ConvertEmptyStringsToNull 構造体の空文字列フィールドをNULLに変換する汎用関数
func ConvertEmptyStringsToNull(obj interface{}) map[string]interface{} {
	values := make(map[string]interface{})

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	// ポインタの場合は要素を取得
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// gormタグからカラム名を取得
		gormTag := fieldType.Tag.Get("gorm")
		if gormTag == "" {
			continue
		}

		// カラム名を抽出（例: "column:net_sales" -> "net_sales"）
		columnName := ExtractColumnName(gormTag)
		if columnName == "" {
			continue
		}

		// 文字列フィールドの場合のみ処理
		if field.Kind() == reflect.String {
			if field.String() == "" {
				values[columnName] = nil
			} else {
				values[columnName] = field.String()
			}
		} else {
			// 文字列以外はそのまま
			values[columnName] = field.Interface()
		}
	}

	return values
}

// ExtractColumnName gormタグからカラム名を抽出
func ExtractColumnName(gormTag string) string {
	// "column:net_sales;primaryKey" から "net_sales" を抽出
	if len(gormTag) > 7 && gormTag[:7] == "column:" {
		// セミコロンで分割して最初の部分（カラム名）を取得
		parts := strings.Split(gormTag[7:], ";")
		return parts[0]
	}
	return ""
}
