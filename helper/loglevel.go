package helper

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// SetupLoggerWithLevel は指定されたログレベルでロガーを設定します
func SetupLoggerWithLevel(specifiedLevel string) {
	var logLevel string

	// --verboseフラグが指定されている場合はそれを使用、なければ環境変数
	if specifiedLevel != "" {
		logLevel = strings.ToLower(specifiedLevel)
	} else {
		logLevel = strings.ToLower(os.Getenv("LOG_LEVEL"))
		if logLevel == "" {
			logLevel = "info"
		}
	}

	// ログレベルを設定
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
		slog.Warn("未知のログレベル、INFOレベルに設定しました", "specified_level", logLevel)
	}

	// プログラムルートディレクトリを取得
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	projectRoot := filepath.Dir(filename)

	// ハンドラーを作成してグローバルロガーに設定
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // ファイル名と行数を表示
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// sourceファイルパスを相対パスに変更
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					if relPath, err := filepath.Rel(projectRoot, source.File); err == nil {
						source.File = relPath
					}
				}
			}
			return a
		},
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Debug("ログレベルを設定しました", "level", strings.ToUpper(logLevel))
}
