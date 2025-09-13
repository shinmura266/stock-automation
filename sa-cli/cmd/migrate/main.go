package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"kabu-analysis/config"
	"kabu-analysis/database"
	"kabu-analysis/errors"
	"kabu-analysis/services"

	"github.com/joho/godotenv"
)

func main() {
	// .envファイルから環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("システム環境変数を使用します")
	}

	// サブコマンド必須チェック
	if len(os.Args) < 2 {
		printUsageAndExit()
	}
	command := os.Args[1]

	// 設定を読み込み
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("設定読み込みエラー: %v", errors.ConfigError(err))
	}

	// データベース接続
	fmt.Println("データベースに接続中...")
	dbService, err := services.NewDatabaseService(cfg.Database)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", errors.DatabaseError(err))
	}
	defer dbService.Close()

	// マイグレーターを作成
	migrator, err := database.NewMigrator(dbService.GetConnection())
	if err != nil {
		log.Fatalf("マイグレーター作成エラー: %v", err)
	}
	defer migrator.Close()

	// サブコマンド実行
	switch command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("マイグレーション実行エラー: %v", err)
		}
		fmt.Println("マイグレーションが完了しました")

	case "down":
		if err := migrator.DownOne(); err != nil {
			log.Fatalf("1ステップロールバックエラー: %v", err)
		}
		fmt.Println("1ステップのロールバックが完了しました")

	case "goto":
		if len(os.Args) < 3 {
			fmt.Println("エラー: goto には <version> が必要です")
			printUsageAndExit()
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil || v < 0 {
			fmt.Println("エラー: version は0以上の整数で指定してください")
			os.Exit(1)
		}
		if err := migrator.Goto(uint(v)); err != nil {
			log.Fatalf("指定バージョンへの移動エラー: %v", err)
		}
		fmt.Printf("バージョン %d へ移動が完了しました\n", v)

	case "version":
		version, dirty, err := migrator.Version()
		if err != nil {
			log.Fatalf("マイグレーションバージョン取得エラー: %v", err)
		}
		fmt.Printf("現在のマイグレーションバージョン: %d (dirty: %t)\n", version, dirty)

	default:
		fmt.Printf("未知のコマンド: %s\n", command)
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	fmt.Println("使い方:")
	fmt.Println("  migrate up                   # すべて最新まで適用")
	fmt.Println("  migrate down                 # 直前へ1つだけロールバック")
	fmt.Println("  migrate goto <version>       # 指定バージョンへ移動")
	fmt.Println("  migrate version              # 現在のバージョン表示")
	os.Exit(1)
}
