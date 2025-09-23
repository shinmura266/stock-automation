package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// GotoCmd は指定されたバージョンへ移動するコマンドです
var GotoCmd = &cobra.Command{
	Use:   "goto <version>",
	Short: "指定バージョンへ移動",
	Long:  "データベースを指定されたマイグレーションバージョンに移動します",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil || version < 0 {
			slog.Error("エラー: version は0以上の整数で指定してください")
			os.Exit(1)
		}
		executeGotoWithVersion(uint(version))
		CloseMigrator()
	},
}

// executeGotoWithVersion は指定されたバージョンへ移動します
func executeGotoWithVersion(version uint) {
	migrator := GetMigrator()
	fmt.Printf("バージョンへ移動中... version=%d\n", version)
	if err := migrator.Migrate(version); err != nil {
		slog.Error("指定バージョンへの移動エラー", "error", err, "version", version)
		os.Exit(1)
	}
	fmt.Printf("バージョンへ移動が完了しました version=%d\n", version)
	showVersion()
}
