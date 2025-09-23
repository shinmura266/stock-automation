# Database Library

株式市場データ管理用の共通データベースライブラリです。

## 機能

- MySQL データベース接続管理
- ジェネリック対応リポジトリクラス：
  - 財務情報（FinancialStatements）
  - 上場銘柄情報（ListedInfo）
  - 日次四本値（DailyQuotes）

## 使用方法

### 1. インポート

```go
import (
    "stock-automation/database"
    "stock-automation/schema"  // API型定義が必要な場合
)
```

### 2. データベース接続

```go
// 設定作成 (connection.go内の機能)
config := database.NewConfig("localhost", 3306, "user", "password", "database")

// 接続作成
conn, err := database.NewConnection(config)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

### 3. リポジトリ使用（ジェネリック対応）

```go
// JQuants APIレスポンス取得
var response schema.FinancialStatementsResponse
// ... API呼び出し

// 財務情報リポジトリ
stmtRepo := database.NewStatementsRepository(conn)
count, err := stmtRepo.SaveFinancialStatements(&response)

// 上場銘柄情報リポジトリ
var listedInfo schema.ListedInfoResponse
// ... API呼び出し
listedRepo := database.NewListedInfoRepository(conn)
err = listedRepo.SaveListedInfo(&listedInfo)

// 日次四本値リポジトリ
var dailyQuotes schema.DailyQuotesResponse
// ... API呼び出し
quotesRepo := database.NewDailyQuotesRepository(conn)
err = quotesRepo.SaveDailyQuotes(&dailyQuotes)
```

## 依存関係

- github.com/go-sql-driver/mysql v1.7.1

## 注意事項

- データベース接続は使用後必ずClose()してください
- トランザクションは自動的にロールバック機能付きで管理されています
- マイグレーションは別プログラムで管理してください
