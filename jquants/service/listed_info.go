package service

import (
	"fmt"
	"log/slog"
	"stock-automation/database"
	"stock-automation/jquants/api"
)

// ListedInfoService 上場銘柄情報サービスクラス
type ListedInfoService struct {
	client     *api.Client
	dbConn     *database.Connection
	repository *database.ListedInfoRepository
}

// NewListedInfoService 新しい上場銘柄情報サービスを作成
func NewListedInfoService() (*ListedInfoService, error) {
	// データベース接続を作成
	dbConn, err := database.NewConnectionFromEnv()
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// リポジトリを作成
	repository := database.NewListedInfoRepository(dbConn)

	return &ListedInfoService{
		client:     api.NewClient(),
		dbConn:     dbConn,
		repository: repository,
	}, nil
}

// Close サービスを閉じる
func (s *ListedInfoService) Close() error {
	if s.dbConn != nil {
		return s.dbConn.Close()
	}
	return nil
}

// UpdateListedInfo 上場銘柄情報を取得し、DBに保存
// date: 日付（空の場合はAPIの最新日付）
func (s *ListedInfoService) UpdateListedInfo(date string) error {
	idToken, err := s.client.AuthClient.GetIdToken()
	if err != nil {
		return fmt.Errorf("IDトークン取得エラー: %v", err)
	}

	slog.Debug("上場銘柄情報取得開始", "date", date)

	// 上場銘柄情報を取得（pagination_key対応で全データ取得）
	listedInfo, err := s.client.ListedClient.GetListedInfo(idToken, date)
	if err != nil {
		return fmt.Errorf("上場銘柄情報取得エラー: %v", err)
	}

	if len(listedInfo) == 0 {
		slog.Warn("取得したデータがありません", "date", date)
		return nil
	}

	slog.Debug("上場銘柄情報取得完了", "count", len(listedInfo))

	// データベースに保存
	if err := s.repository.SaveListedInfo(listedInfo); err != nil {
		return fmt.Errorf("データベース保存エラー: %v", err)
	}

	slog.Info("上場銘柄情報更新完了", "count", len(listedInfo))
	return nil
}
