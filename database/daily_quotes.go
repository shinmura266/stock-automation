package database

import (
	"fmt"
	"log/slog"
	"time"

	"stock-automation/schema"
)

// DailyQuotesRepository 日次四本値のリポジトリ
type DailyQuotesRepository struct {
	conn *Connection
}

// NewDailyQuotesRepository 新しいリポジトリを作成
func NewDailyQuotesRepository(conn *Connection) *DailyQuotesRepository {
	return &DailyQuotesRepository{
		conn: conn,
	}
}

// SaveDailyQuotes 四本値データを保存
func (r *DailyQuotesRepository) SaveDailyQuotes(dailyQuotes []schema.DailyQuote) error {
	if len(dailyQuotes) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// タイムスタンプを設定
	quotes := make([]schema.DailyQuote, len(dailyQuotes))
	now := time.Now()
	for i, quote := range dailyQuotes {
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

// GetDailyQuotes 条件に基づいて四本値データを取得
func (r *DailyQuotesRepository) GetDailyQuotes(code, date string) ([]schema.DailyQuote, error) {
	var quotes []schema.DailyQuote
	query := r.conn.GetGormDB().Model(&schema.DailyQuote{})

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
