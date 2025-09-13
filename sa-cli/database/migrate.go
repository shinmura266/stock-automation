package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator マイグレーション管理
type Migrator struct {
	migrate *migrate.Migrate
}

// NewMigrator 新しいマイグレーターを作成
func NewMigrator(conn *Connection) (*Migrator, error) {
	// MySQLドライバーを取得
	driver, err := mysql.WithInstance(conn.GetDB(), &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("MySQLドライバー作成エラー: %v", err)
	}

	// マイグレーションファイルのパスを取得
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("作業ディレクトリ取得エラー: %v", err)
	}

	migrationsPath := filepath.Join(workDir, "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("マイグレーションディレクトリが見つかりません: %s", migrationsPath)
	}

	// マイグレーターを作成
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"mysql", driver)
	if err != nil {
		return nil, fmt.Errorf("マイグレーター作成エラー: %v", err)
	}

	return &Migrator{migrate: m}, nil
}

// Up マイグレーションを実行
func (m *Migrator) Up() error {
	log.Println("マイグレーションを実行中...")

	err := m.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("マイグレーション実行エラー: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("マイグレーションは既に最新です")
	} else {
		log.Println("マイグレーションが完了しました")
	}

	return nil
}

// Down マイグレーションをロールバック
func (m *Migrator) Down() error {
	log.Println("マイグレーションをロールバック中...")

	err := m.migrate.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("マイグレーションロールバックエラー: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("ロールバックするマイグレーションがありません")
	} else {
		log.Println("マイグレーションのロールバックが完了しました")
	}

	return nil
}

// DownOne 直前のバージョンまで1つだけロールバック
func (m *Migrator) DownOne() error {
	log.Println("直前のバージョンへ1つだけロールバック中...")

	if err := m.migrate.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("ロールバックするマイグレーションがありません")
			return nil
		}
		return fmt.Errorf("1ステップロールバックエラー: %v", err)
	}

	log.Println("1ステップのロールバックが完了しました")
	return nil
}

// Goto 指定したバージョンへ移動（Up/Downを実行して到達）
func (m *Migrator) Goto(version uint) error {
	log.Printf("バージョン %d へ移動中...", version)

	if err := m.migrate.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("指定バージョンへ変更はありません")
			return nil
		}
		return fmt.Errorf("指定バージョンへの移動エラー: %v", err)
	}

	log.Printf("バージョン %d へ移動が完了しました", version)
	return nil
}

// Force マイグレーションを強制実行
func (m *Migrator) Force(version int) error {
	log.Printf("マイグレーションをバージョン %d に強制設定中...", version)

	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("マイグレーション強制設定エラー: %v", err)
	}

	log.Printf("マイグレーションがバージョン %d に設定されました", version)
	return nil
}

// Version 現在のマイグレーションバージョンを取得
func (m *Migrator) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		return 0, false, fmt.Errorf("マイグレーションバージョン取得エラー: %v", err)
	}

	return version, dirty, nil
}

// Close マイグレーターを閉じる
func (m *Migrator) Close() error {
	if m.migrate != nil {
		if _, err := m.migrate.Close(); err != nil {
			return err
		}
	}
	return nil
}
