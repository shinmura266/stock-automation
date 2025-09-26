package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config データベース設定
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Connection データベース接続
type Connection struct {
	db   *sql.DB
	gorm *gorm.DB
}

// NewConnectionFromEnv 環境変数からデータベース接続を作成
func NewConnectionFromEnv() (*Connection, error) {
	config, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewConnection(config)
}

// loadDatabaseConfigFromEnv データベース設定を環境変数から読み込む
func NewConfigFromEnv() (*Config, error) {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || portStr == "" || user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("データベース設定の環境変数を設定してください (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("DB_PORTの値が無効です: %v", err)
	}

	return NewConfig(host, port, user, password, dbName), nil
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

// NewConnection 新しいデータベース接続を作成
func NewConnection(config *Config) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// カスタムロガーを作成（スローログの閾値を1秒に設定）
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             1 * time.Second, // スロークエリの閾値
			LogLevel:                  logger.Info,     // ログレベル
			IgnoreRecordNotFoundError: true,            // ErrRecordNotFound エラーを無視
			Colorful:                  false,           // カラー出力を無効化
		},
	)

	// GORM接続を作成
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// 基底のsql.DBを取得
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("sql.DB取得エラー: %v", err)
	}

	// 接続テスト
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %v", err)
	}

	slog.Debug("データベースに接続しました",
		"host", config.Host,
		"port", config.Port,
		"database", config.Database)

	return &Connection{db: sqlDB, gorm: gormDB}, nil
}

// Close データベース接続を閉じる
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetDB データベースインスタンスを取得
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

// GetGormDB GORMデータベースインスタンスを取得
func (c *Connection) GetGormDB() *gorm.DB {
	return c.gorm
}
