package main

import (
	"log/slog"
	"os"
	"stock-automation/helper"
	"stock-automation/jquants/cmd"

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

var rootCmd = &cobra.Command{
	Use:   "jquants",
	Short: "J-Quantsデータ取得CLI",
	Long:  "J-Quantsの各種データを取得して、DBへ保存するためのCLIツール",
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
	rootCmd.AddCommand(cmd.DailyQuotesCmd)
}
