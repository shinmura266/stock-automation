package database

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"stock-automation/schema"
)

// ListedInfoRepository 上場銘柄情報のリポジトリ
type ListedInfoRepository struct {
	conn *Connection
}

// NewListedInfoRepository 新しいリポジトリを作成
func NewListedInfoRepository(conn *Connection) *ListedInfoRepository {
	return &ListedInfoRepository{
		conn: conn,
	}
}

// SaveListedInfo 上場銘柄情報をデータベースに保存（ジェネリック対応）
func (r *ListedInfoRepository) SaveListedInfo(listedInfo interface{}) error {
	// リフレクションを使用して動的に型を処理
	v := reflect.ValueOf(listedInfo)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Info配列を取得
	infoField := v.FieldByName("Info")
	if !infoField.IsValid() || infoField.Kind() != reflect.Slice {
		return fmt.Errorf("Infoフィールドが見つからないか、スライスではありません")
	}

	if infoField.Len() == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// トランザクション開始
	tx, cleanup := BeginTransaction(r.conn.GetDB())
	defer cleanup()

	// バッチ挿入用のプリペアドステートメント（スキーマ準拠）
	// 既存データのeffective_dateより新しい場合のみ更新
	stmt, err := tx.Prepare(`
		INSERT INTO listed_info (
			effective_date, code, company_name, company_name_english,
			sector17_code, sector33_code, scale_category,
			market_code, margin_code, margin_code_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			effective_date = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(effective_date)
				ELSE effective_date
			END,
			company_name = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(company_name)
				ELSE company_name
			END,
			company_name_english = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(company_name_english)
				ELSE company_name_english
			END,
			sector17_code = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(sector17_code)
				ELSE sector17_code
			END,
			sector33_code = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(sector33_code)
				ELSE sector33_code
			END,
			scale_category = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(scale_category)
				ELSE scale_category
			END,
			market_code = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(market_code)
				ELSE market_code
			END,
			margin_code = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(margin_code)
				ELSE margin_code
			END,
			margin_code_name = CASE 
				WHEN VALUES(effective_date) > effective_date THEN VALUES(margin_code_name)
				ELSE margin_code_name
			END,
			updated_at = CASE 
				WHEN VALUES(effective_date) > effective_date THEN CURRENT_TIMESTAMP
				ELSE updated_at
			END
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("プリペアドステートメント作成エラー: %v", err)
	}
	defer stmt.Close()

	// データをバッチ挿入
	insertedCount := 0

	for i := 0; i < infoField.Len(); i++ {
		info := infoField.Index(i)

		// 日付をパース
		dateStr := GetStringField(info, "Date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			slog.Warn("日付パースエラー",
				"code", GetStringField(info, "Code"),
				"date", dateStr,
				"error", err)
			continue
		}

		result, err := stmt.Exec(
			date,
			GetStringField(info, "Code"),
			GetStringField(info, "CompanyName"),
			GetStringField(info, "CompanyNameEnglish"),
			GetStringField(info, "Sector17Code"),
			GetStringField(info, "Sector33Code"),
			GetStringField(info, "ScaleCategory"),
			GetStringField(info, "MarketCode"),
			GetStringField(info, "MarginCode"),
			GetStringField(info, "MarginCodeName"),
		)

		if err != nil {
			slog.Error("データ挿入エラー",
				"code", GetStringField(info, "Code"),
				"error", err)
			continue
		}

		// 挿入された行数を確認
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			insertedCount++
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	slog.Info("データベース保存完了", "inserted_count", insertedCount)
	return nil
}

// GetListedInfo その他市場（0109）を除外して上場銘柄情報を取得（昇順）
func (r *ListedInfoRepository) GetListedInfo(startCode string, limit int) ([]schema.ListedInfo, error) {
	var listedInfos []schema.ListedInfo
	query := r.conn.GetGormDB().Model(&schema.ListedInfo{}).
		Where("market_code != ?", "0109")

	// startCodeが空文字列でない場合のみフィルターを適用
	if startCode != "" {
		query = query.Where("code >= ?", startCode)
	}

	query = query.Order("code")

	// limitが0より大きい場合のみ制限を適用
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&listedInfos).Error
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}

	return listedInfos, nil
}
