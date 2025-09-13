package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kabu",
	Short: "J-Quantsデータ更新CLI",
	Long:  "J-Quantsの各種データをDBへ保存するためのCLIツール",
}

// Execute ルートコマンド実行
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
