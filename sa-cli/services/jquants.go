package services

import (
	"kabu-analysis/config"
	"kabu-analysis/jquants/client"
)

// JQuantsService J-Quants API接続を管理するサービス
type JQuantsService struct {
	client *client.Client
	config *config.JQuantsConfig
}

// NewJQuantsService 新しいJ-Quantsサービスを作成
func NewJQuantsService(cfg *config.JQuantsConfig) *JQuantsService {
	c := client.NewClient(cfg.BaseURL)

	return &JQuantsService{
		client: c,
		config: cfg,
	}
}

// Authenticate J-Quants APIに認証
func (s *JQuantsService) Authenticate() error {
	return s.client.Authenticate(s.config.Email, s.config.Password)
}

// GetClient 認証済みのクライアントを取得
func (s *JQuantsService) GetClient() *client.Client {
	return s.client
}

// GetIDToken IDトークンを取得
func (s *JQuantsService) GetIDToken() string {
	return s.client.GetIDToken()
}
