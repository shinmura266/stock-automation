package service

import (
	"fmt"
	"log/slog"
	"stock-automation/database"
	"stock-automation/helper"
	"stock-automation/jquants/api"
	"time"
)

// StatementsService 財務情報サービスクラス
type StatementsService struct {
	client     *api.Client
	dbConn     *database.Connection
	repository *database.StatementsRepository
}

// NewStatementsService 新しい財務情報サービスを作成
func NewStatementsService() (*StatementsService, error) {
	// データベース接続を作成
	dbConn, err := database.NewConnectionFromEnv()
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// リポジトリを作成
	repository := database.NewStatementsRepository(dbConn)

	return &StatementsService{
		client:     api.NewClient(),
		dbConn:     dbConn,
		repository: repository,
	}, nil
}

// UpdateStatements 財務情報を取得し、DBに保存
// code: 銘柄コード（空の場合は全銘柄）
// date: 日付（空の場合は当日、ただしcodeが指定されている場合は全期間）
func (s *StatementsService) UpdateStatements(code, date string, debug bool) error {
	// codeもdateも両方とも空文字の場合は当日を使用
	if code == "" && date == "" {
		date = helper.GetTodayDate()
	}

	idToken, err := s.client.AuthClient.GetIdToken()
	if err != nil {
		return fmt.Errorf("IDトークン取得エラー: %v", err)
	}

	statements, err := s.client.StatementsClient.GetAllStatements(idToken, code, date, debug)
	if err != nil {
		return fmt.Errorf("財務情報取得エラー: %v", err)
	}

	// データベースに保存
	if len(statements.Statements) > 0 {
		if count, err := s.repository.SaveFinancialStatements(statements); err != nil {
			return fmt.Errorf("データベース保存エラー: %v", err)
		} else {
			if code != "" {
				slog.Info("銘柄別財務情報保存完了", "code", code, "date", date, "count", count)
			} else {
				slog.Info("日付別財務情報保存完了", "date", date, "count", count)
			}
		}
	}

	return nil
}

// UpdateStatementsByCodeAndDate 指定された銘柄コードと日付の財務情報を取得し、DBに保存
func (s *StatementsService) UpdateStatementsByCodeAndDate(code, date string, debug bool) error {
	if code == "" {
		return fmt.Errorf("銘柄コードが指定されていません")
	}
	if date == "" {
		return fmt.Errorf("日付が指定されていません")
	}

	return s.UpdateStatements(code, date, debug)
}

// UpdateStatementsMultipleDates 複数日付の財務情報を取得し、DBに保存（間隔制御付き）
func (s *StatementsService) UpdateStatementsMultipleDates(codes []string, dates []string, intervalMs int, debug bool) error {
	if len(dates) == 0 {
		return fmt.Errorf("日付が指定されていません")
	}

	if intervalMs <= 0 {
		intervalMs = 1000 // デフォルト1秒間隔
	}

	slog.Info("複数日付財務情報取得・保存開始", "dates_count", len(dates), "codes_count", len(codes))

	successCount := 0
	for i, date := range dates {
		slog.Info("日付別財務情報取得・保存中", "date", date, "progress", fmt.Sprintf("%d/%d", i+1, len(dates)))

		if len(codes) == 0 {
			// 全銘柄取得
			err := s.UpdateStatements("", date, debug)
			if err != nil {
				slog.Error("全銘柄財務情報取得・保存エラー", "date", date, "error", err)
				continue
			}
		} else {
			// 指定銘柄取得
			for _, code := range codes {
				err := s.UpdateStatements(code, date, debug)
				if err != nil {
					slog.Error("銘柄別財務情報取得・保存エラー", "code", code, "date", date, "error", err)
					continue
				}

				// 銘柄間の間隔制御
				if intervalMs > 0 {
					time.Sleep(time.Duration(intervalMs) * time.Millisecond)
				}
			}
		}

		successCount++

		// 日付間の間隔制御（最後の日付以外）
		if i < len(dates)-1 && intervalMs > 0 {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}

	slog.Info("複数日付財務情報取得・保存完了", "success_count", successCount)
	return nil
}

// Close データベース接続を閉じる
func (s *StatementsService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}
