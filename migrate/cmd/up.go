package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// UpCmd は全てのマイグレーションを適用するコマンドです
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "全てのマイグレーションを適用",
	Long:  "データベースに全ての未適用マイグレーションを適用します",
	Run: func(cmd *cobra.Command, args []string) {
		executeUp()
		CloseMigrator()
	},
}

// executeUp は全てのマイグレーションを適用します
func executeUp() {
	migrator := GetMigrator()
	fmt.Println("全てのマイグレーションを適用中...")
	if err := migrator.Up(); err != nil {
		// "no change"エラーの場合は既に最新であることを表示
		if strings.Contains(err.Error(), "no change") {
			fmt.Println("既に最新です")
			showVersion()
			return
		}
		slog.Error("マイグレーション実行エラー", "error", err)
		os.Exit(1)
	}
	fmt.Println("マイグレーションが完了しました")
	showVersion()
}
