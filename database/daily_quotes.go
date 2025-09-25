package database

import (
	"fmt"
	"log/slog"
	"stock-automation/schema"
	"time"
)

// DailyQuotesRepository 日次四本値のリポジトリ
type DailyQuotesRepository struct {
	conn *Connection
}

// NewDailyQuotesRepository 新しいリポジトリを作成
func NewDailyQuotesRepository(conn *Connection) *DailyQuotesRepository {
	return &DailyQuotesRepository{conn: conn}
}

// SaveDailyQuotes 四本値データを保存（バッチ処理版）
func (r *DailyQuotesRepository) SaveDailyQuotes(resp []schema.DailyQuote) error {
	if len(resp) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// タイムスタンプを設定
	quotes := make([]schema.DailyQuote, len(resp))
	now := time.Now()
	for i, quote := range resp {
		quotes[i] = quote
		quotes[i].CreatedAt = now
		quotes[i].UpdatedAt = now
	}

	// バッチサイズを制限（MySQLのプレースホルダー制限を回避）
	const batchSize = 100
	db := r.conn.GetGormDB()

	for i := 0; i < len(quotes); i += batchSize {
		end := i + batchSize
		if end > len(quotes) {
			end = len(quotes)
		}

		batch := quotes[i:end]
		result := db.Save(&batch)
		if result.Error != nil {
			return fmt.Errorf("データベース保存エラー (バッチ %d-%d): %v", i+1, end, result.Error)
		}

		slog.Debug("daily_quotesバッチ保存完了", "batch", fmt.Sprintf("%d-%d", i+1, end), "count", len(batch))
	}

	slog.Debug("daily_quotes保存完了", "total_count", len(quotes))
	return nil
}

// SaveDailyQuotesBatch バッチで四本値データを保存（より効率的なUPSERT）
func (r *DailyQuotesRepository) SaveDailyQuotesBatch(quotes []schema.DailyQuote) error {
	if len(quotes) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// タイムスタンプを設定
	now := time.Now()
	for i := range quotes {
		quotes[i].CreatedAt = now
		quotes[i].UpdatedAt = now
	}

	// GORMでバッチUPSERT実行
	db := r.conn.GetGormDB()
	result := db.Save(&quotes)
	if result.Error != nil {
		return fmt.Errorf("データベース保存エラー: %v", result.Error)
	}

	slog.Debug("daily_quotesバッチ保存完了", "count", len(quotes))
	return nil
}

// GetDailyQuotes 条件に基づいて四本値データを取得
func (r *DailyQuotesRepository) GetDailyQuotes(code, date string) ([]schema.DailyQuote, error) {
	var quotes []schema.DailyQuote
	db := r.conn.GetGormDB()

	query := db.Model(&schema.DailyQuote{})

	if code != "" {
		query = query.Where("code = ?", code)
	}
	if date != "" {
		query = query.Where("trade_date = ?", date)
	}

	result := query.Find(&quotes)
	if result.Error != nil {
		return nil, fmt.Errorf("データ取得エラー: %v", result.Error)
	}

	return quotes, nil
}

// DeleteDailyQuotes 条件に基づいて四本値データを削除
func (r *DailyQuotesRepository) DeleteDailyQuotes(code, date string) error {
	db := r.conn.GetGormDB()

	query := db.Model(&schema.DailyQuote{})

	if code != "" {
		query = query.Where("code = ?", code)
	}
	if date != "" {
		query = query.Where("trade_date = ?", date)
	}

	result := query.Delete(&schema.DailyQuote{})
	if result.Error != nil {
		return fmt.Errorf("データ削除エラー: %v", result.Error)
	}

	slog.Debug("daily_quotes削除完了", "affected_rows", result.RowsAffected)
	return nil
}
