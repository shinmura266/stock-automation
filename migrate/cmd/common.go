package cmd

import (
	"log/slog"
	"os"

	dbpkg "stock-automation/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrator *migrate.Migrate

// GetMigrator はマイグレーターインスタンスを取得します
func GetMigrator() *migrate.Migrate {
	if migrator == nil {
		migrator = initMigrator()
	}
	return migrator
}

// CloseMigrator はマイグレーターをクローズします
func CloseMigrator() {
	if migrator != nil {
		migrator.Close()
	}
}

// initMigrator はマイグレーターを初期化します
func initMigrator() *migrate.Migrate {
	// 環境変数から設定を読み込み
	migratePath := os.Getenv("MIGRATE_PATH")
	if migratePath == "" {
		slog.Error("設定読み込みエラー: 必要な環境変数(MIGRATE_PATH)が設定されていません")
		os.Exit(1)
	}

	// データベース接続を作成
	slog.Debug("データベースに接続中...")
	conn, err := dbpkg.NewConnectionFromEnv()
	if err != nil {
		slog.Error("データベース設定の作成エラー", "error", err)
		os.Exit(1)
	}
	slog.Debug("データベース接続が確立されました")

	// sql.DBインスタンスを取得
	db := conn.GetDB()

	// マイグレーターを作成
	slog.Debug("マイグレーターを作成中...")
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		slog.Error("MySQL ドライバー作成エラー", "error", err)
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://"+migratePath,
		"mysql",
		driver)
	if err != nil {
		slog.Error("マイグレーター作成エラー", "error", err)
		os.Exit(1)
	}
	slog.Debug("マイグレーター作成完了", "path", migratePath)

	return migrator
}
