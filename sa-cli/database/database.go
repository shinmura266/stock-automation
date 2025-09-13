package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Config データベース設定
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewConfig 新しいデータベース設定を作成
func NewConfig(host string, port int, user, password, database string) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

// Connection データベース接続
type Connection struct {
	db *sql.DB
}

// NewConnection 新しいデータベース接続を作成
func NewConnection(config *Config) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %v", err)
	}

	log.Println("データベースに接続しました")

	return &Connection{db: db}, nil
}

// Close データベース接続を閉じる
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// CreateTables 必要なテーブルを作成（マイグレーションを使用）
func (c *Connection) CreateTables() error {
	// マイグレーターを作成
	migrator, err := NewMigrator(c)
	if err != nil {
		return fmt.Errorf("マイグレーター作成エラー: %v", err)
	}
	defer migrator.Close()

	// マイグレーションを実行
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("マイグレーション実行エラー: %v", err)
	}

	return nil
}

// GetDB データベースインスタンスを取得
func (c *Connection) GetDB() *sql.DB {
	return c.db
}
