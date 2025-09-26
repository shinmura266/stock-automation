package database

import (
	"fmt"
	"log/slog"
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

// SaveListedInfo 上場銘柄情報をデータベースに保存
func (r *ListedInfoRepository) SaveListedInfo(listedInfos []schema.ListedInfo) error {
	if len(listedInfos) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// タイムスタンプを設定
	infos := make([]schema.ListedInfo, len(listedInfos))
	now := time.Now()
	for i, info := range listedInfos {
		infos[i] = info
		infos[i].CreatedAt = now
		infos[i].UpdatedAt = now
	}

	// バッチサイズを制限（MySQLのプレースホルダー制限を回避）
	const batchSize = 100
	db := r.conn.GetGormDB()

	for i := 0; i < len(infos); i += batchSize {
		end := i + batchSize
		if end > len(infos) {
			end = len(infos)
		}

		batch := infos[i:end]
		// 存在しないカラムを除外して保存
		result := db.Select("effective_date", "code", "company_name", "company_name_english",
			"sector17_code", "sector33_code", "scale_category", "market_code", "margin_code",
			"margin_code_name", "created_at", "updated_at").Save(&batch)
		if result.Error != nil {
			return fmt.Errorf("データベース保存エラー (バッチ %d-%d): %v", i+1, end, result.Error)
		}

		slog.Debug("listed_infoバッチ保存完了", "batch", fmt.Sprintf("%d-%d", i+1, end), "count", len(batch))
	}

	slog.Debug("listed_info保存完了", "total_count", len(infos))
	return nil
}

// GetListedInfo その他市場（0109）を除外して上場銘柄情報を取得（昇順）
func (r *ListedInfoRepository) GetListedInfo(startCode string, limit int) ([]schema.ListedInfo, error) {
	var infos []schema.ListedInfo
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

	result := query.Find(&infos)
	if result.Error != nil {
		return nil, fmt.Errorf("データ取得エラー: %v", result.Error)
	}

	return infos, nil
}
