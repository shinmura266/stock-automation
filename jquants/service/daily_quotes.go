package service

import (
	"fmt"
	"log/slog"
	"stock-automation/database"
	"stock-automation/helper"
	"stock-automation/jquants/api"
	"stock-automation/schema"
	"time"
)

// DailyQuotesService 日次株価四本値サービスクラス
type DailyQuotesService struct {
	client     *api.Client
	dbConn     *database.Connection
	repository *database.DailyQuotesRepository
	interval   int // インターバル（秒）
}

// NewDailyQuotesService 新しい日次株価四本値サービスを作成
func NewDailyQuotesService(interval int) (*DailyQuotesService, error) {
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
		interval:   interval,
	}, nil
}

// UpdateDailyQuotes 株価データを取得し、DBに保存
// code: 銘柄コード（空の場合は全銘柄）
// date: 日付（空の場合は当日、ただしcodeが指定されている場合は全期間）
func (s *DailyQuotesService) UpdateDailyQuotes(code, date string) error {
	// codeもdateも両方とも空文字の場合は当日を使用
	if code == "" && date == "" {
		date = helper.GetTodayDate()
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

// GetDailyQuotesMultipleDates 複数日付の株価データを取得（間隔制御付き）
func (s *DailyQuotesService) GetDailyQuotesMultipleDates(codes []string, dates []string, intervalMs int) (map[string]map[string][]schema.DailyQuote, error) {
	if len(dates) == 0 {
		return nil, fmt.Errorf("日付が指定されていません")
	}

	if intervalMs <= 0 {
		intervalMs = 5000 // デフォルト5秒間隔
	}

	result := make(map[string]map[string][]schema.DailyQuote)

	for _, date := range dates {
		slog.Info("日付別株価データ取得開始", "date", date)

		dateResult := make(map[string][]schema.DailyQuote)

		if len(codes) == 0 {
			// 全銘柄取得
			err := s.UpdateDailyQuotes("", date)
			if err != nil {
				slog.Error("全銘柄株価データ取得・保存エラー", "date", date, "error", err)
				continue
			}
			// 保存後は空のスライスを返す（実際のデータはDBに保存済み）
			dateResult["ALL"] = []schema.DailyQuote{}
		} else {
			// 指定銘柄取得
			for _, code := range codes {
				err := s.UpdateDailyQuotes(code, date)
				if err != nil {
					slog.Error("銘柄別株価データ取得・保存エラー", "code", code, "date", date, "error", err)
					continue
				}
				// 保存後は空のスライスを返す（実際のデータはDBに保存済み）
				dateResult[code] = []schema.DailyQuote{}

				// 銘柄間の間隔制御
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
		date = helper.GetTodayDate()
	}

	if intervalMs <= 0 {
		intervalMs = 5000 // デフォルト5秒間隔
	}

	slog.Info("全銘柄株価データ取得・保存開始", "date", date)

	err := s.UpdateDailyQuotes("", date)
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
		date = helper.GetTodayDate()
	}

	if intervalMs <= 0 {
		intervalMs = 5000 // デフォルト5秒間隔
	}

	slog.Info("複数銘柄株価データ取得・保存開始", "date", date, "codes_count", len(codes))

	successCount := 0
	for i, code := range codes {
		slog.Info("銘柄株価データ取得・保存中", "code", code, "progress", fmt.Sprintf("%d/%d", i+1, len(codes)))

		err := s.UpdateDailyQuotes(code, date)
		if err != nil {
			slog.Error("銘柄株価データ取得・保存エラー", "code", code, "error", err)
			continue
		}

		successCount++

		// 銘柄間の間隔制御（最後のリクエスト以外）
		if i < len(codes)-1 && intervalMs > 0 {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}

	slog.Info("複数銘柄株価データ取得・保存完了", "date", date, "success_count", successCount)
	return nil
}

// UpdateDailyQuotesMultipleDates 複数日付の株価データを取得し、DBに保存（間隔制御付き）
// date: 開始日付
// count: 取得する日数
func (s *DailyQuotesService) UpdateDailyQuotesMultipleDates(date string, count int) error {
	if count <= 0 {
		return fmt.Errorf("countが0以下です")
	}

	slog.Info("複数日付株価データ取得・保存開始", "start_date", date, "count", count)

	for i := 0; i < count; i++ {
		// 指定日付からi日分さかのぼった日付を計算
		currentDate := helper.SubDate(date, i)
		slog.Info("日付別株価データ取得・保存中", "date", currentDate, "progress", fmt.Sprintf("%d/%d", i+1, count))

		err := s.UpdateDailyQuotes("", currentDate)
		if err != nil {
			return fmt.Errorf("全銘柄株価データ取得・保存エラー (date: %s): %v", currentDate, err)
		}

		// 日付間の間隔制御（最後の日付以外）
		if i < count-1 && s.interval > 0 {
			time.Sleep(time.Duration(s.interval) * time.Second)
		}
	}

	slog.Info("複数日付株価データ取得・保存完了", "count", count)
	return nil
}

// UpdateDailyQuotesWithCount 株価データを取得し、DBに保存（count対応版）
// code: 銘柄コード（空の場合は全銘柄）
// date: 日付（空の場合は当日）
// count: 取得する日数または銘柄数
func (s *DailyQuotesService) UpdateDailyQuotesWithCount(code, date string, count int) error {
	// codeもdateも両方とも空文字の場合は当日を使用
	if code == "" && date == "" {
		date = helper.GetTodayDate()
	}

	// countが2以上の場合の処理
	if count >= 2 && code != "" && date != "" {
		return fmt.Errorf("countが2以上の場合は、codeとdateのどちらか一方のみを指定してください")
	}

	// countが2未満の場合は単一実行
	if count < 2 {
		err := s.UpdateDailyQuotes(code, date)
		if err != nil {
			return fmt.Errorf("株価データ更新エラー: %v", err)
		}
		slog.Info("株価データ更新完了")
		return nil
	}

	if date != "" {
		// 指定日付からcount日数分をさかのぼって繰り返し実行
		err := s.UpdateDailyQuotesMultipleDates(date, count)
		if err != nil {
			return fmt.Errorf("株価データ更新エラー: %v", err)
		}
		slog.Info("株価データ更新完了", "days", count)
	} else {
		// 指定コードからcount分のコードを昇順で取得して繰り返し実行
		codes, err := s.getCodesFromListedInfo(code, count)
		if err != nil {
			return fmt.Errorf("銘柄コード取得エラー: %v", err)
		}

		for _, c := range codes {
			err := s.UpdateDailyQuotes(c, "")
			if err != nil {
				slog.Error("銘柄株価データ更新エラー", "code", c, "error", err)
				continue
			}

			// インターバル制御
			if s.interval > 0 {
				time.Sleep(time.Duration(s.interval) * time.Second)
			}
		}
		slog.Info("株価データ更新完了", "codes", len(codes))
	}

	return nil
}

// getCodesFromListedInfo 指定されたコードからcount分のコードを昇順で取得
func (s *DailyQuotesService) getCodesFromListedInfo(startCode string, count int) ([]string, error) {
	// リポジトリを作成
	repository := database.NewListedInfoRepository(s.dbConn)

	// 指定されたコードから昇順でcount分のコードを取得
	codes, err := repository.GetListedCodesExcludingMarket("", startCode, count)
	if err != nil {
		return nil, fmt.Errorf("銘柄コード取得エラー: %v", err)
	}

	return codes, nil
}

// Close データベース接続を閉じる
func (s *DailyQuotesService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}
