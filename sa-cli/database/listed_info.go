package database

import (
	"fmt"
	"log"
	"time"

	"kabu-analysis/jquants"
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

// SaveListedInfo 上場銘柄情報をデータベースに保存
func (r *ListedInfoRepository) SaveListedInfo(listedInfo *jquants.ListedInfoResponse) error {
	if len(listedInfo.Info) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// トランザクション開始
	tx, cleanup := beginTransaction(r.conn.GetDB())
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
	updatedCount := 0

	for _, info := range listedInfo.Info {
		// 日付をパース
		date, err := time.Parse("2006-01-02", info.Date)
		if err != nil {
			log.Printf("日付パースエラー (コード: %s, 日付: %s): %v", info.Code, info.Date, err)
			continue
		}

		result, err := stmt.Exec(
			date,
			info.Code,
			info.CompanyName,
			info.CompanyNameEnglish,
			info.Sector17Code,
			info.Sector33Code,
			info.ScaleCategory,
			info.MarketCode,
			info.MarginCode,
			info.MarginCodeName,
		)

		if err != nil {
			log.Printf("データ挿入エラー (コード: %s): %v", info.Code, err)
			continue
		}

		// 挿入された行数を確認
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			insertedCount++
		} else {
			updatedCount++
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	log.Printf("データベース保存完了: 新規挿入 %d件, 更新 %d件", insertedCount, updatedCount)
	return nil
}

// DailyQuotesRepository 日次四本値のリポジトリ
type DailyQuotesRepository struct {
	conn *Connection
}

// NewDailyQuotesRepository 新しいリポジトリを作成
func NewDailyQuotesRepository(conn *Connection) *DailyQuotesRepository {
	return &DailyQuotesRepository{conn: conn}
}

// SaveDailyQuotes 四本値データを保存
func (r *DailyQuotesRepository) SaveDailyQuotes(resp *jquants.DailyQuotesResponse) error {
	if len(resp.DailyQuotes) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	tx, cleanup := beginTransaction(r.conn.GetDB())
	defer cleanup()

	stmt, err := tx.Prepare(`
        INSERT INTO daily_quotes (
            trade_date, code,
            open, high, low, close,
            volume, turnover_value,
            adjustment_open, adjustment_high, adjustment_low, adjustment_close, adjustment_volume
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            open = VALUES(open),
            high = VALUES(high),
            low = VALUES(low),
            close = VALUES(close),
            volume = VALUES(volume),
            turnover_value = VALUES(turnover_value),
            adjustment_open = VALUES(adjustment_open),
            adjustment_high = VALUES(adjustment_high),
            adjustment_low = VALUES(adjustment_low),
            adjustment_close = VALUES(adjustment_close),
            adjustment_volume = VALUES(adjustment_volume),
            updated_at = CURRENT_TIMESTAMP
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("プリペアドステートメント作成エラー: %v", err)
	}
	defer stmt.Close()

	insertedCount := 0
	updatedCount := 0

	for _, q := range resp.DailyQuotes {
		date, err := time.Parse("2006-01-02", q.Date)
		if err != nil {
			log.Printf("日付パースエラー (コード: %s, 日付: %s): %v", q.Code, q.Date, err)
			continue
		}

		result, err := stmt.Exec(
			date, q.Code,
			q.Open, q.High, q.Low, q.Close,
			q.Volume, q.TurnoverValue,
			q.AdjustmentOpen, q.AdjustmentHigh, q.AdjustmentLow, q.AdjustmentClose, q.AdjustmentVolume,
		)
		if err != nil {
			log.Printf("データ挿入エラー (コード: %s): %v", q.Code, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			insertedCount++
		} else {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	log.Printf("daily_quotes保存完了: 新規挿入 %d件, 更新 %d件", insertedCount, updatedCount)
	return nil
}

// GetListedInfoByCode 銘柄コードで銘柄情報を取得
func (r *ListedInfoRepository) GetListedInfoByCode(code string) ([]jquants.ListedInfo, error) {
	query := `
		SELECT effective_date, code, company_name, company_name_english,
			   sector17_code, sector33_code,
			   scale_category, market_code,
			   margin_code, margin_code_name
		FROM listed_info
		WHERE code = ?
		ORDER BY effective_date DESC
	`

	rows, err := r.conn.GetDB().Query(query, code)
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}
	defer rows.Close()

	var results []jquants.ListedInfo
	for rows.Next() {
		var info jquants.ListedInfo
		var date time.Time

		err := rows.Scan(
			&date, &info.Code, &info.CompanyName, &info.CompanyNameEnglish,
			&info.Sector17Code, &info.Sector33Code,
			&info.ScaleCategory, &info.MarketCode,
			&info.MarginCode, &info.MarginCodeName,
		)
		if err != nil {
			log.Printf("行スキャンエラー: %v", err)
			continue
		}

		info.Date = date.Format("2006-01-02")
		results = append(results, info)
	}

	return results, nil
}

// GetListedInfoByDate 日付で銘柄情報を取得
func (r *ListedInfoRepository) GetListedInfoByDate(date string) ([]jquants.ListedInfo, error) {
	query := `
		SELECT effective_date, code, company_name, company_name_english,
			   sector17_code, sector33_code,
			   scale_category, market_code,
			   margin_code, margin_code_name
		FROM listed_info
		WHERE effective_date = ?
		ORDER BY code
	`

	rows, err := r.conn.GetDB().Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}
	defer rows.Close()

	var results []jquants.ListedInfo
	for rows.Next() {
		var info jquants.ListedInfo
		var dateTime time.Time

		err := rows.Scan(
			&dateTime, &info.Code, &info.CompanyName, &info.CompanyNameEnglish,
			&info.Sector17Code, &info.Sector33Code,
			&info.ScaleCategory, &info.MarketCode,
			&info.MarginCode, &info.MarginCodeName,
		)
		if err != nil {
			log.Printf("行スキャンエラー: %v", err)
			continue
		}

		info.Date = dateTime.Format("2006-01-02")
		results = append(results, info)
	}

	return results, nil
}

// GetLatestListedInfo 最新の銘柄情報を取得
func (r *ListedInfoRepository) GetLatestListedInfo() ([]jquants.ListedInfo, error) {
	query := `
		SELECT effective_date, code, company_name, company_name_english,
			   sector17_code, sector33_code,
			   scale_category, market_code,
			   margin_code, margin_code_name
		FROM listed_info
		WHERE effective_date = (SELECT MAX(effective_date) FROM listed_info)
		ORDER BY code
	`

	rows, err := r.conn.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}
	defer rows.Close()

	var results []jquants.ListedInfo
	for rows.Next() {
		var info jquants.ListedInfo
		var dateTime time.Time

		err := rows.Scan(
			&dateTime, &info.Code, &info.CompanyName, &info.CompanyNameEnglish,
			&info.Sector17Code, &info.Sector33Code,
			&info.ScaleCategory, &info.MarketCode,
			&info.MarginCode, &info.MarginCodeName,
		)
		if err != nil {
			log.Printf("行スキャンエラー: %v", err)
			continue
		}

		info.Date = dateTime.Format("2006-01-02")
		results = append(results, info)
	}

	return results, nil
}

// GetListedCodesExcludingMarket 特定の市場コードを除外して銘柄コードリストを取得（昇順）
func (r *ListedInfoRepository) GetListedCodesExcludingMarket(excludeMarketCode string, startCode string, limit int) ([]string, error) {
	query := `
		SELECT DISTINCT code
		FROM listed_info
		WHERE effective_date = (SELECT MAX(effective_date) FROM listed_info)
		AND market_code != ?
		AND code >= ?
		ORDER BY code
		LIMIT ?
	`

	rows, err := r.conn.GetDB().Query(query, excludeMarketCode, startCode, limit)
	if err != nil {
		return nil, fmt.Errorf("クエリ実行エラー: %v", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			log.Printf("行スキャンエラー: %v", err)
			continue
		}
		codes = append(codes, code)
	}

	return codes, nil
}
