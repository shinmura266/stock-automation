package services

import (
	"kabu-analysis/config"
	"kabu-analysis/database"
)

// DatabaseService データベース接続を管理するサービス
type DatabaseService struct {
	connection *database.Connection
}

// NewDatabaseService 新しいデータベースサービスを作成
func NewDatabaseService(cfg *config.DatabaseConfig) (*DatabaseService, error) {
	dbConfig := database.NewConfig(
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
	)

	conn, err := database.NewConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	return &DatabaseService{
		connection: conn,
	}, nil
}

// GetConnection データベース接続を取得
func (s *DatabaseService) GetConnection() *database.Connection {
	return s.connection
}

// Close データベース接続を閉じる
func (s *DatabaseService) Close() error {
	if s.connection != nil {
		return s.connection.Close()
	}
	return nil
}
