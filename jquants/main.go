package main

import (
	"log/slog"
	"os"
	"stock-automation/helper"
	"stock-automation/jquants/cmd"

	"github.com/spf13/cobra"
)

var (
	verbose bool
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
		// .envファイルから環境変数を読み込み
		helper.LoadDotEnv()

		// ログレベルを設定
		var logLevel string
		if verbose {
			logLevel = "debug"
		}
		// verboseがfalseの場合は空文字列を渡し、helper.SetupLoggerWithLevelで環境変数を処理
		helper.SetupLoggerWithLevel(logLevel)
	},
}

func init() {
	// グローバルフラグを追加
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "詳細ログを出力")

	// サブコマンドを追加
	rootCmd.AddCommand(cmd.DailyQuotesCmd)
	rootCmd.AddCommand(cmd.StatementsCmd)
	rootCmd.AddCommand(cmd.ListedInfoCmd)
}
