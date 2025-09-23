package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// VersionCmd は現在のマイグレーションバージョンを表示するコマンドです
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "現在のバージョン表示",
	Long:  "現在のデータベースマイグレーションバージョンを表示します",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
		CloseMigrator()
	},
}

// ShowVersion は現在のマイグレーションバージョンを表示します
func showVersion() {
	migrator := GetMigrator()
	slog.Debug("現在のマイグレーションバージョンを取得中...")
	version, dirty, err := migrator.Version()
	if err != nil {
		slog.Error("マイグレーションバージョン取得エラー", "error", err)
		os.Exit(1)
	}
	fmt.Printf("現在のマイグレーションバージョン version=%d dirty=%t\n", version, dirty)
}
