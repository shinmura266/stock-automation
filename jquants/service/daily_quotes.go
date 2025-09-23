package service

import (
	"fmt"
	"log/slog"
	"stock-automation/database"
	"stock-automation/jquants/api"
	"stock-automation/schema"
	"time"
)

// DailyQuotesService 日次株価四本値サービスクラス
type DailyQuotesService struct {
	client     *api.Client
	dbConn     *database.Connection
	repository *database.DailyQuotesRepository
}

// NewDailyQuotesService 新しい日次株価四本値サービスを作成
func NewDailyQuotesService() (*DailyQuotesService, error) {
	// データベース接続を作成
	dbConn, err := database.NewConnectionFromEnv()
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// リポジトリを作成
	repository := database.NewDailyQuotesRepository(dbConn)

	return &DailyQuotesService{
		client:     api.NewClient(),
		dbConn:     dbConn,
		repository: repository,
	}, nil
}

// UpdateDailyQuotesByCode 指定された銘柄コードの株価データを取得し、DBに保存
func (s *DailyQuotesService) UpdateDailyQuotesByCode(code string, date string) error {
	if code == "" {
		return fmt.Errorf("銘柄コードが指定されていません")
	}

	idToken, err := s.client.AuthClient.GetIdToken()
	if err != nil {
		return fmt.Errorf("IDトークン取得エラー: %v", err)
	}
	quotes, err := s.client.DailyQuotesClient.GetDailyQuotes(idToken, code, date)
	if err != nil {
		return fmt.Errorf("株価データ取得エラー: %v", err)
	}

	// データベースに保存
	if len(quotes) > 0 {
		if err := s.repository.SaveDailyQuotes(quotes); err != nil {
			return fmt.Errorf("データベース保存エラー: %v", err)
		}
		slog.Info("銘柄別株価データ保存完了", "code", code, "date", date, "count", len(quotes))
	}

	return nil
}

// UpdateDailyQuotesByDate 指定された日付の全銘柄株価データを取得し、DBに保存
func (s *DailyQuotesService) UpdateDailyQuotesByDate(date string) error {
	if date == "" {
		return fmt.Errorf("日付が指定されていません")
	}

	idToken, err := s.client.AuthClient.GetIdToken()
	if err != nil {
		return fmt.Errorf("IDトークン取得エラー: %v", err)
	}
	quotes, err := s.client.DailyQuotesClient.GetDailyQuotes(idToken, "", date)
	if err != nil {
		return fmt.Errorf("株価データ取得エラー: %v", err)
	}

	// データベースに保存
	if len(quotes) > 0 {
		if err := s.repository.SaveDailyQuotes(quotes); err != nil {
			return fmt.Errorf("データベース保存エラー: %v", err)
		}
		slog.Info("日付別全銘柄株価データ保存完了", "date", date, "count", len(quotes))
	}

	return nil
}

// GetDailyQuotesMultipleDates 複数日付の株価データを取得（間隔制御付き）
func (s *DailyQuotesService) GetDailyQuotesMultipleDates(codes []string, dates []string, intervalMs int) (map[string]map[string][]schema.DailyQuote, error) {
	if len(dates) == 0 {
		return nil, fmt.Errorf("日付が指定されていません")
	}

	if intervalMs <= 0 {
		intervalMs = 1000 // デフォルト1秒間隔
	}

	result := make(map[string]map[string][]schema.DailyQuote)

	for _, date := range dates {
		slog.Info("日付別株価データ取得開始", "date", date)

		dateResult := make(map[string][]schema.DailyQuote)

		if len(codes) == 0 {
			// 全銘柄取得
			err := s.UpdateDailyQuotesByDate(date)
			if err != nil {
				slog.Error("全銘柄株価データ取得・保存エラー", "date", date, "error", err)
				continue
			}
			// 保存後は空のスライスを返す（実際のデータはDBに保存済み）
			dateResult["ALL"] = []schema.DailyQuote{}
		} else {
			// 指定銘柄取得
			for _, code := range codes {
				err := s.UpdateDailyQuotesByCode(code, date)
				if err != nil {
					slog.Error("銘柄別株価データ取得・保存エラー", "code", code, "date", date, "error", err)
					continue
				}
				// 保存後は空のスライスを返す（実際のデータはDBに保存済み）
				dateResult[code] = []schema.DailyQuote{}

				// 間隔制御
				if intervalMs > 0 {
					time.Sleep(time.Duration(intervalMs) * time.Millisecond)
				}
			}
		}

		result[date] = dateResult

		// 日付間の間隔制御
		if intervalMs > 0 {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}

	return result, nil
}

// UpdateDailyQuotesAllStocks 全銘柄の株価データを取得し、DBに保存（間隔制御付き）
func (s *DailyQuotesService) UpdateDailyQuotesAllStocks(date string, intervalMs int) error {
	if date == "" {
		date = getTodayDate()
	}

	if intervalMs <= 0 {
		intervalMs = 1000 // デフォルト1秒間隔
	}

	slog.Info("全銘柄株価データ取得・保存開始", "date", date)

	err := s.UpdateDailyQuotesByDate(date)
	if err != nil {
		return fmt.Errorf("全銘柄株価データ取得・保存エラー: %v", err)
	}

	slog.Info("全銘柄株価データ取得・保存完了", "date", date)
	return nil
}

// UpdateDailyQuotesMultipleStocks 複数銘柄の株価データを取得し、DBに保存（間隔制御付き）
func (s *DailyQuotesService) UpdateDailyQuotesMultipleStocks(codes []string, date string, intervalMs int) error {
	if len(codes) == 0 {
		return fmt.Errorf("銘柄コードが指定されていません")
	}

	if date == "" {
		date = getTodayDate()
	}

	if intervalMs <= 0 {
		intervalMs = 1000 // デフォルト1秒間隔
	}

	slog.Info("複数銘柄株価データ取得・保存開始", "date", date, "codes_count", len(codes))

	successCount := 0
	for i, code := range codes {
		slog.Info("銘柄株価データ取得・保存中", "code", code, "progress", fmt.Sprintf("%d/%d", i+1, len(codes)))

		err := s.UpdateDailyQuotesByCode(code, date)
		if err != nil {
			slog.Error("銘柄株価データ取得・保存エラー", "code", code, "error", err)
			continue
		}

		successCount++

		// 間隔制御（最後のリクエスト以外）
		if i < len(codes)-1 && intervalMs > 0 {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}

	slog.Info("複数銘柄株価データ取得・保存完了", "date", date, "success_count", successCount)
	return nil
}

// UpdateDailyQuotesMultipleDates 複数日付の株価データを取得し、DBに保存（間隔制御付き）
func (s *DailyQuotesService) UpdateDailyQuotesMultipleDates(codes []string, dates []string, intervalMs int) error {
	if len(dates) == 0 {
		return fmt.Errorf("日付が指定されていません")
	}

	if intervalMs <= 0 {
		intervalMs = 1000 // デフォルト1秒間隔
	}

	slog.Info("複数日付株価データ取得・保存開始", "dates_count", len(dates), "codes_count", len(codes))

	successCount := 0
	for i, date := range dates {
		slog.Info("日付別株価データ取得・保存中", "date", date, "progress", fmt.Sprintf("%d/%d", i+1, len(dates)))

		if len(codes) == 0 {
			// 全銘柄取得
			err := s.UpdateDailyQuotesByDate(date)
			if err != nil {
				slog.Error("全銘柄株価データ取得・保存エラー", "date", date, "error", err)
				continue
			}
		} else {
			// 指定銘柄取得
			for _, code := range codes {
				err := s.UpdateDailyQuotesByCode(code, date)
				if err != nil {
					slog.Error("銘柄別株価データ取得・保存エラー", "code", code, "date", date, "error", err)
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

	slog.Info("複数日付株価データ取得・保存完了", "success_count", successCount)
	return nil
}

// Close データベース接続を閉じる
func (s *DailyQuotesService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}

// getTodayDate 当日の日付をYYYYMMDD形式で取得
func getTodayDate() string {
	return time.Now().Format("20060102")
}
