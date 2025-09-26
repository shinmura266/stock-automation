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
	interval   int // インターバル（秒）
}

// NewStatementsService 新しい財務情報サービスを作成
func NewStatementsService(interval int, verbose bool) (*StatementsService, error) {
	// データベース接続を作成
	dbConn, err := database.NewConnectionFromEnv(verbose)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// リポジトリを作成
	repository := database.NewStatementsRepository(dbConn)

	return &StatementsService{
		client:     api.NewClient(),
		dbConn:     dbConn,
		repository: repository,
		interval:   interval,
	}, nil
}

// UpdateStatements 財務情報を取得し、DBに保存
// code: 銘柄コード（空の場合は全銘柄）
// date: 日付（空の場合は当日、ただしcodeが指定されている場合は全期間）
func (s *StatementsService) UpdateStatements(code, date string) error {
	// codeもdateも両方とも空文字の場合は当日を使用
	if code == "" && date == "" {
		date = helper.GetTodayDate()
	}

	idToken, err := s.client.AuthClient.GetIdToken()
	if err != nil {
		return fmt.Errorf("IDトークン取得エラー: %v", err)
	}

	statements, err := s.client.StatementsClient.GetStatements(idToken, code, date)
	if err != nil {
		return fmt.Errorf("財務情報取得エラー: %v", err)
	}

	// データベースに保存
	if len(statements) > 0 {
		if err := s.repository.SaveFinancialStatements(statements); err != nil {
			return fmt.Errorf("データベース保存エラー: %v", err)
		}
		slog.Info("財務情報保存完了", "code", code, "date", date, "count", len(statements))
	}

	return nil
}

// UpdateStatementsMultipleDates 複数日付の財務情報を取得し、DBに保存（間隔制御付き）
// date: 開始日付
// count: 取得する日数
func (s *StatementsService) UpdateStatementsMultipleDates(date string, count int) error {
	if count <= 0 {
		return fmt.Errorf("countが0以下です")
	}

	slog.Info("複数日付財務情報取得・保存開始", "start_date", date, "count", count)

	for i := 0; i < count; i++ {
		// 指定日付からi日分さかのぼった日付を計算
		currentDate := helper.SubDate(date, i)
		slog.Info("日付別財務情報取得・保存中", "date", currentDate, "progress", fmt.Sprintf("%d/%d", i+1, count))

		err := s.UpdateStatements("", currentDate)
		if err != nil {
			return fmt.Errorf("全銘柄財務情報取得・保存エラー (date: %s): %v", currentDate, err)
		}

		// 日付間の間隔制御（最後の日付以外）
		if i < count-1 && s.interval > 0 {
			time.Sleep(time.Duration(s.interval) * time.Second)
		}
	}

	slog.Info("複数日付財務情報取得・保存完了", "count", count)
	return nil
}

// UpdateStatementsMultipleCodes 複数銘柄の財務情報を取得し、DBに保存（間隔制御付き）
// code: 銘柄コード（空の場合は全銘柄）
// count: 取得する銘柄数
func (s *StatementsService) UpdateStatementsMultipleCodes(code string, count int) error {
	if count <= 0 {
		return fmt.Errorf("countが0以下です")
	}

	slog.Info("複数銘柄財務情報取得・保存開始", "code", code, "count", count)

	// リポジトリを作成
	repository := database.NewListedInfoRepository(s.dbConn)

	// 指定されたコードから昇順でcount分の銘柄情報を取得
	listedInfos, err := repository.GetListedInfo(code, count)
	if err != nil {
		return fmt.Errorf("銘柄情報取得エラー: %v", err)
	}

	for _, info := range listedInfos {
		err := s.UpdateStatements(info.Code, "")
		if err != nil {
			slog.Error("銘柄財務情報更新エラー", "code", info.Code, "error", err)
			continue
		}

		// インターバル制御
		if s.interval > 0 {
			time.Sleep(time.Duration(s.interval) * time.Second)
		}
	}
	slog.Info("財務情報更新完了", "count", len(listedInfos))
	return nil
}

// UpdateStatementsWithCount 財務情報を取得し、DBに保存（count対応版）
// code: 銘柄コード（空の場合は全銘柄）
// date: 日付（空の場合は当日）
// count: 取得する日数または銘柄数
func (s *StatementsService) UpdateStatementsWithCount(code, date string, count int) error {
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
		err := s.UpdateStatements(code, date)
		if err != nil {
			return fmt.Errorf("財務情報更新エラー: %v", err)
		}
		slog.Info("財務情報更新完了")
		return nil
	}

	if date != "" {
		// 指定日付からcount日数分をさかのぼって繰り返し実行
		err := s.UpdateStatementsMultipleDates(date, count)
		if err != nil {
			return fmt.Errorf("財務情報更新エラー: %v", err)
		}
	} else {
		// 指定コードからcount分のコードを昇順で取得して繰り返し実行
		err := s.UpdateStatementsMultipleCodes(code, count)
		if err != nil {
			return fmt.Errorf("財務情報更新エラー: %v", err)
		}
	}
	slog.Info("財務情報更新完了", "count", count)

	return nil
}

// Close データベース接続を閉じる
func (s *StatementsService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}
