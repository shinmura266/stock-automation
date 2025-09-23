package main

import (
	"log/slog"
	"os"

	"migrate/cmd"
	"stock-automation/helper"

	"github.com/spf13/cobra"
)

var (
	logLevel string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("コマンド実行エラー", "error", err)
		os.Exit(1)
	}
}

// rootCmd はルートコマンドを定義します
var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "データベースマイグレーションツール",
	Long:  "golang-migrateを使用したデータベースマイグレーションツール",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// ログレベルを設定
		helper.SetupLoggerWithLevel(logLevel)

		// .envファイルから環境変数を読み込み
		helper.LoadDotEnv()
	},
}

func init() {
	// グローバルフラグを追加
	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "", "ログレベル (debug, info, warn, error)")

	// サブコマンドを追加
	rootCmd.AddCommand(cmd.UpCmd)
	rootCmd.AddCommand(cmd.DownCmd)
	rootCmd.AddCommand(cmd.GotoCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
}
