package query

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "利用可能なテーブル一覧を表示",
	Long:  "データベース内の利用可能なテーブル一覧を表示します",
	RunE:  listTables,
}

func listTables(cmd *cobra.Command, args []string) error {
	fmt.Printf("\n=== サポートしているテーブル一覧 ===\n\n")

	tables := []struct {
		name        string
		description string
	}{
		{"listed_info", "上場銘柄情報"},
		{"market_codes", "市場区分コード"},
	}

	for i, table := range tables {
		fmt.Printf("%d. %-15s - %s\n", i+1, table.name, table.description)
	}

	fmt.Printf("\n総テーブル数: %d\n", len(tables))
	fmt.Println("\n使用例:")
	fmt.Printf("  kabu db show listed_info     # 上場銘柄情報を表示\n")
	fmt.Printf("  kabu db show market_codes    # 市場区分コードを表示\n")
	fmt.Printf("  kabu db show listed_info -l 20  # 最大20行表示\n")
	fmt.Printf("  kabu db show listed_info -a     # 全行表示\n")

	return nil
}
