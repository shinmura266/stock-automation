package query

import (
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "データベース操作",
	Long:  "データベースのテーブル内容を表示する機能を提供します",
}

func init() {
	QueryCmd.AddCommand(showCmd)
	QueryCmd.AddCommand(listCmd)
}
