package database

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"stock-automation/schema"

	"gorm.io/gorm/clause"
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

// SaveFinancialStatements 財務情報をデータベースに保存（シンプル版）
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

	db := r.conn.GetGormDB()

	// 汎用関数を使用してシンプルに実装
	for i, stmt := range financialStatements {
		// 汎用関数で空文字列をNULLに変換
		values := ConvertEmptyStringsToNull(&stmt)

		// ON DUPLICATE KEY UPDATE を使用してUPSERT
		result := db.Model(&schema.FinancialStatement{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "disclosed_date"},
				{Name: "local_code"},
				{Name: "type_of_current_period"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"disclosed_time", "disclosure_number", "type_of_document",
				"current_period_start_date", "current_period_end_date",
				"current_fiscal_year_start_date", "current_fiscal_year_end_date",
				"next_fiscal_year_start_date", "next_fiscal_year_end_date",
				"net_sales", "operating_profit", "ordinary_profit", "profit",
				"eps", "diluted_eps", "total_assets", "equity",
				"equity_to_asset_ratio", "bvps", "cf_operating", "cf_investing",
				"cf_financing", "cash_and_equivalents", "updated_at",
			}),
		}).Create(&values)

		if result.Error != nil {
			// 外部キー制約エラー（1452）の場合はログを出力して続行
			if strings.Contains(result.Error.Error(), "Error 1452") && strings.Contains(result.Error.Error(), "foreign key constraint") {
				slog.Debug("外部キー制約エラーでスキップ",
					"local_code", stmt.LocalCode,
					"disclosed_date", stmt.DisclosedDate,
					"error", result.Error.Error())
				continue
			}
			return fmt.Errorf("データベース保存エラー (レコード %d): %v", i+1, result.Error)
		}

		if (i+1)%10 == 0 {
			slog.Debug("statements保存進捗", "progress", fmt.Sprintf("%d/%d", i+1, len(financialStatements)))
		}
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
