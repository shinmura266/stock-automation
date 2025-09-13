package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config アプリケーション全体の設定
type Config struct {
	Database *DatabaseConfig
	JQuants  *JQuantsConfig
}

// DatabaseConfig データベース接続設定
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// JQuantsConfig J-Quants API設定
type JQuantsConfig struct {
	Email    string
	Password string
	BaseURL  string
}

// LoadFromEnv 環境変数から設定を読み込む
func LoadFromEnv() (*Config, error) {
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	jquantsConfig, err := loadJQuantsConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Database: dbConfig,
		JQuants:  jquantsConfig,
	}, nil
}

// loadDatabaseConfig データベース設定を環境変数から読み込む
func loadDatabaseConfig() (*DatabaseConfig, error) {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	if host == "" || portStr == "" || user == "" || password == "" || database == "" {
		return nil, fmt.Errorf("データベース設定の環境変数を設定してください (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("DB_PORTの値が無効です: %v", err)
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}, nil
}

// loadJQuantsConfig J-Quants設定を環境変数から読み込む
func loadJQuantsConfig() (*JQuantsConfig, error) {
	email := os.Getenv("JQUANTS_EMAIL")
	password := os.Getenv("JQUANTS_PASSWORD")

	if email == "" || password == "" {
		return nil, fmt.Errorf("JQUANTS_EMAIL と JQUANTS_PASSWORD の環境変数を設定してください")
	}

	return &JQuantsConfig{
		Email:    email,
		Password: password,
		BaseURL:  "https://api.jquants.com/v1",
	}, nil
}
