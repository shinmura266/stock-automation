package helper

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Debug("警告: .envファイルの読み込みに失敗しました", "error", err)
		slog.Debug("システム環境変数を使用します")
	}
}
