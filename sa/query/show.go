package query

import (
	"fmt"
	"os"
	"stock-automation/database"
	"stock-automation/schema"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var showCmd = &cobra.Command{
	Use:   "show [table_name]",
	Short: "テーブル内容を表示",
	Long:  "指定したテーブルの内容を表示します\n\n利用可能なテーブル:\n  - listed_info: 上場銘柄情報\n  - market_codes: 市場区分コード",
	Args:  cobra.ExactArgs(1),
	RunE:  showTable,
}

func init() {
	// フラグを追加
	showCmd.Flags().IntP("limit", "l", 10, "表示する行数の上限")
	showCmd.Flags().BoolP("all", "a", false, "全ての行を表示")
}

func showTable(cmd *cobra.Command, args []string) error {
	tableName := args[0]

	// サポートするテーブルを限定
	supportedTables := map[string]string{
		"listed_info":  "上場銘柄情報",
		"market_codes": "市場区分コード",
	}

	_, supported := supportedTables[tableName]
	if !supported {
		return fmt.Errorf("サポートされていないテーブルです: '%s'\n\n利用可能なテーブル:\n  - listed_info: 上場銘柄情報\n  - market_codes: 市場区分コード", tableName)
	}

	// データベース接続
	conn, err := database.NewConnectionFromEnv(false) // queryでは非verbose
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %v", err)
	}
	defer conn.Close()

	gormDB := conn.GetGormDB()

	// フラグの値を取得
	limit, _ := cmd.Flags().GetInt("limit")
	showAll, _ := cmd.Flags().GetBool("all")

	// テーブル固有の表示処理
	switch tableName {
	case "listed_info":
		return showListedInfo(gormDB, limit, showAll)
	case "market_codes":
		return showMarketCodes(gormDB, limit, showAll)
	default:
		return fmt.Errorf("未実装のテーブル: %s", tableName)
	}
}

// listed_info テーブル専用の表示関数（GORM版）
func showListedInfo(gormDB *gorm.DB, limit int, showAll bool) error {
	fmt.Printf("\n=== 上場銘柄情報 (listed_info) ===\n\n")

	var listedInfos []schema.ListedInfo

	query := gormDB.Order("code")

	if !showAll {
		query = query.Limit(limit)
	}

	if err := query.Find(&listedInfos).Error; err != nil {
		return fmt.Errorf("データ取得エラー: %v", err)
	}

	// 市場コードのマップを作成（1回のクエリで全市場コードを取得）
	var marketCodes []schema.MarketCode
	if err := gormDB.Find(&marketCodes).Error; err != nil {
		return fmt.Errorf("市場コード取得エラー: %v", err)
	}

	marketMap := make(map[string]string)
	for _, mc := range marketCodes {
		marketMap[mc.Code] = mc.Name
	}

	// ヘッダーを表示
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "コード\t企業名\t17業種\t33業種\t市場コード\t市場名\t規模区分")
	fmt.Fprintln(w, "----\t----\t----\t----\t----\t----\t----")

	for _, info := range listedInfos {
		marketName := marketMap[info.MarketCode]
		if marketName == "" {
			marketName = "不明"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			info.Code, info.CompanyName, info.Sector17CodeName,
			info.Sector33CodeName, info.MarketCode, marketName, info.ScaleCategory)
	}

	w.Flush()

	if len(listedInfos) == 0 {
		fmt.Println("データが見つかりませんでした")
	} else {
		fmt.Printf("\n表示行数: %d\n", len(listedInfos))
	}

	return nil
}

// market_codes テーブル専用の表示関数（GORM版）
func showMarketCodes(gormDB *gorm.DB, limit int, showAll bool) error {
	fmt.Printf("\n=== 市場区分コード (market_codes) ===\n\n")

	var marketCodes []schema.MarketCode

	query := gormDB.Order("code")

	if !showAll {
		query = query.Limit(limit)
	}

	if err := query.Find(&marketCodes).Error; err != nil {
		return fmt.Errorf("データ取得エラー: %v", err)
	}

	// ヘッダーを表示
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "コード\t市場名")
	fmt.Fprintln(w, "----\t----")

	for _, code := range marketCodes {
		fmt.Fprintf(w, "%s\t%s\n", code.Code, code.Name)
	}

	w.Flush()

	if len(marketCodes) == 0 {
		fmt.Println("データが見つかりませんでした")
	} else {
		fmt.Printf("\n表示行数: %d\n", len(marketCodes))
	}

	return nil
}
