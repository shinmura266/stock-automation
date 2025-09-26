package database

import (
	"fmt"
	"log/slog"
	"time"

	"stock-automation/schema"
)

// StatementsRepository 財務情報のリポジトリ
type StatementsRepository struct {
	conn *Connection
}

// NewStatementsRepository 新しいリポジトリを作成
func NewStatementsRepository(conn *Connection) *StatementsRepository {
	return &StatementsRepository{
		conn: conn,
	}
}

// SaveFinancialStatements 財務情報をデータベースに保存
func (r *StatementsRepository) SaveFinancialStatements(statements []schema.FinancialStatement) error {
	if len(statements) == 0 {
		return fmt.Errorf("保存するデータがありません")
	}

	// タイムスタンプを設定
	financialStatements := make([]schema.FinancialStatement, len(statements))
	now := time.Now()
	for i, stmt := range statements {
		financialStatements[i] = stmt
		financialStatements[i].CreatedAt = now
		financialStatements[i].UpdatedAt = now
	}

	// バッチサイズを制限（MySQLのプレースホルダー制限を回避）
	const batchSize = 100
	db := r.conn.GetGormDB()

	for i := 0; i < len(financialStatements); i += batchSize {
		end := i + batchSize
		if end > len(financialStatements) {
			end = len(financialStatements)
		}

		batch := financialStatements[i:end]
		result := db.Save(&batch)
		if result.Error != nil {
			return fmt.Errorf("データベース保存エラー (バッチ %d-%d): %v", i+1, end, result.Error)
		}

		slog.Debug("statementsバッチ保存完了", "batch", fmt.Sprintf("%d-%d", i+1, end), "count", len(batch))
	}

	slog.Debug("statements保存完了", "total_count", len(financialStatements))
	return nil
}

// GetFinancialStatements 条件に基づいて財務情報を取得
func (r *StatementsRepository) GetFinancialStatements(localCode, disclosedDate, typeOfCurrentPeriod string) ([]schema.FinancialStatement, error) {
	var statements []schema.FinancialStatement
	query := r.conn.GetGormDB().Model(&schema.FinancialStatement{})

	if localCode != "" {
		query = query.Where("local_code = ?", localCode)
	}
	if disclosedDate != "" {
		query = query.Where("disclosed_date = ?", disclosedDate)
	}
	if typeOfCurrentPeriod != "" {
		query = query.Where("type_of_current_period = ?", typeOfCurrentPeriod)
	}

	result := query.Find(&statements)
	if result.Error != nil {
		return nil, fmt.Errorf("データ取得エラー: %v", result.Error)
	}

	return statements, nil
}
