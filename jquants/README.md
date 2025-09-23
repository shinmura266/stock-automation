# J-Quants CLI

J-Quantsの各種データを取得して、DBへ保存するためのCLIツールです。

## 概要

このツールは、J-Quants APIを使用して以下のデータを取得・保存する機能を提供します：

- **日次株価四本値データ** - 株価の始値、高値、安値、終値
- **上場銘柄情報** - 銘柄の基本情報
- **財務情報** - 企業の財務諸表データ

## 機能

### 認証機能
- J-Quants APIの認証を自動化
- リフレッシュトークンとIDトークンの自動管理
- トークンの有効期限チェックと自動更新
- 認証情報の安全なファイル保存（`~/.config/stock-automation/token`）

### 日次株価データ取得
- 指定日付の全銘柄株価データ取得
- 指定銘柄の株価データ取得
- 複数銘柄・複数日付の一括取得
- API制限を考慮した間隔制御
- データベースへの自動保存

### データベース連携
- PostgreSQLデータベースへの接続
- 日次株価データの保存・更新
- 重複データの適切な処理

## 使用方法

### 環境設定

1. 環境変数の設定（`.env`ファイルまたは環境変数）：
```bash
JQUANTS_EMAIL=your_email@example.com
JQUANTS_PASSWORD=your_password
```

2. データベース接続設定：
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=stock_automation
DB_USER=your_username
DB_PASSWORD=your_password
```

### コマンド実行

#### 日次株価データ取得

```bash
# 当日の全銘柄株価データを取得
./jquants daily-quotes

# 指定日付の全銘柄株価データを取得
./jquants daily-quotes --date 2024-01-15

# 指定銘柄の株価データを取得
./jquants daily-quotes --code 7203 --date 2024-01-15
```

#### ログレベル設定

```bash
# デバッグログを有効にして実行
./jquants --log debug daily-quotes
```

## プロジェクト構造

```
jquants/
├── main.go                 # メインエントリーポイント
├── go.mod                  # Go モジュール定義
├── cmd/                    # CLI コマンド定義
│   └── daily_quotes.go    # 日次株価コマンド
├── api/                    # J-Quants API クライアント
│   ├── auth.go            # 認証関連
│   ├── client.go          # メインクライアント
│   ├── daily_quotes.go    # 日次株価API
│   ├── listed.go          # 上場銘柄API
│   └── statements.go      # 財務情報API
└── service/               # ビジネスロジック
    └── daily_quotes.go    # 日次株価サービス
```

## 技術仕様

- **言語**: Go 1.24.5
- **フレームワーク**: Cobra CLI
- **データベース**: PostgreSQL
- **認証**: J-Quants API認証
- **ログ**: slog（構造化ログ）

## API制限対応

- リクエスト間隔制御（デフォルト1秒）
- トークンの有効期限管理
- エラーハンドリングとリトライ機能
- ページネーション対応

## セキュリティ

- 認証情報の暗号化保存
- 環境変数による設定管理
- トークンの自動無効化機能

## 開発・ビルド

```bash
# 依存関係のインストール
go mod tidy

# ビルド
go build -o jquants

# 実行
./jquants --help
```

## 注意事項

- J-Quants APIの利用規約を遵守してください
- API制限を考慮して適切な間隔でリクエストを実行してください
- 認証情報は安全に管理してください
- データベース接続設定を正しく行ってください
