package service

import (
	"fmt"
	"log/slog"
	"stock-automation/database"
	"stock-automation/helper"
	"stock-automation/jquants/api"
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

// UpdateDailyQuotesMultipleCodes 複数銘柄の株価データを取得し、DBに保存（間隔制御付き）
// code: 銘柄コード（空の場合は全銘柄）
// count: 取得する銘柄数
func (s *DailyQuotesService) UpdateDailyQuotesMultipleCodes(code string, count int) error {
	if count <= 0 {
		return fmt.Errorf("countが0以下です")
	}

	slog.Info("複数銘柄株価データ取得・保存開始", "code", code, "count", count)

	// リポジトリを作成
	repository := database.NewListedInfoRepository(s.dbConn)

	// 指定されたコードから昇順でcount分の銘柄情報を取得
	listedInfos, err := repository.GetListedInfo(code, count)
	if err != nil {
		return fmt.Errorf("銘柄情報取得エラー: %v", err)
	}

	for _, info := range listedInfos {
		err := s.UpdateDailyQuotes(info.Code, "")
		if err != nil {
			slog.Error("銘柄株価データ更新エラー", "code", info.Code, "error", err)
			continue
		}

		// インターバル制御
		if s.interval > 0 {
			time.Sleep(time.Duration(s.interval) * time.Second)
		}
	}
	slog.Info("株価データ更新完了", "count", len(listedInfos))
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
	} else {
		// 指定コードからcount分のコードを昇順で取得して繰り返し実行
		err := s.UpdateDailyQuotesMultipleCodes(code, count)
		if err != nil {
			return fmt.Errorf("株価データ更新エラー: %v", err)
		}
	}
	slog.Info("株価データ更新完了", "count", count)

	return nil
}

// Close データベース接続を閉じる
func (s *DailyQuotesService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}
