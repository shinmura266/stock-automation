package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// DownCmd は1ステップのロールバックを実行するコマンドです
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "1ステップのロールバックを実行",
	Long:  "直前のマイグレーションを1つロールバックします",
	Run: func(cmd *cobra.Command, args []string) {
		executeDown()
		CloseMigrator()
	},
}

// executeDown は1ステップのロールバックを実行します
func executeDown() {
	migrator := GetMigrator()
	fmt.Println("1ステップのロールバックを実行中...")
	if err := migrator.Steps(-1); err != nil {
		slog.Error("1ステップロールバックエラー", "error", err)
		os.Exit(1)
	}
	fmt.Println("1ステップのロールバックが完了しました")
	showVersion()
}
