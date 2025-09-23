# Stock Automation

J-Quants APIを使用した株式市場データの自動取得・管理システムです。

## 概要

このプロジェクトは、J-Quants APIから株式市場データを取得し、MySQLデータベースに保存するためのGo言語ベースのCLIツール群です。複数のモジュールで構成され、データの取得から管理まで一貫したワークフローを提供します。

## プロジェクト構成

### モジュール構成

- **`sa/`** - メインCLIツール（Stock Analysis）
- **`jquants/`** - J-Quants APIクライアント
- **`migrate/`** - データベースマイグレーションツール
- **`database/`** - データベース共通ライブラリ
- **`schema/`** - API型定義ライブラリ
- **`helper/`** - 共通ユーティリティライブラリ

### ディレクトリ構造

```
stock-automation/
├── sa/                    # メインCLIツール
│   ├── main.go
│   ├── query/            # クエリサブコマンド
│   └── go.mod
├── jquants/              # J-Quants APIクライアント
│   ├── main.go
│   ├── api/              # API実装
│   ├── service/          # ビジネスロジック
│   ├── cmd/              # CLIコマンド
│   └── go.mod
├── migrate/              # データベースマイグレーション
│   ├── main.go
│   ├── cmd/              # マイグレーションコマンド
│   ├── sql/              # SQLマイグレーションファイル
│   └── go.mod
├── database/             # データベースライブラリ
│   ├── connection.go     # DB接続管理
│   ├── daily_quotes.go   # 日次四本値リポジトリ
│   ├── listed_info.go    # 上場銘柄情報リポジトリ
│   ├── statements.go     # 財務情報リポジトリ
│   └── go.mod
├── schema/               # 型定義
│   ├── jquants.go        # J-Quants API型定義
│   └── go.mod
├── helper/               # 共通ユーティリティ
│   ├── env.go            # 環境変数管理
│   ├── loglevel.go       # ログレベル設定
│   └── go.mod
├── bin/                  # ビルド成果物
├── Makefile              # ビルド・管理スクリプト
├── go.work               # Go Workspace設定
└── env.example           # 環境変数テンプレート
```

## 機能

### 1. J-Quants APIクライアント (`jquants/`)

- **認証管理**: リフレッシュトークン・IDトークンの取得・管理
- **データ取得**: 
  - 日次四本値データ (`daily_quotes`)
  - 上場銘柄情報 (`listed_info`)
  - 財務情報 (`financial_statements`)

### 2. データベース管理 (`database/`)

- **接続管理**: MySQL接続の確立・管理
- **リポジトリパターン**: ジェネリック対応のデータアクセス層
- **トランザクション管理**: 自動ロールバック機能付き

### 3. マイグレーション (`migrate/`)

- **golang-migrate**: データベーススキーマのバージョン管理
- **コマンド**: up, down, goto, version

### 4. メインCLI (`sa/`)

- **クエリ機能**: データベースからのデータ検索・表示
- **統合管理**: 各モジュールの統合インターフェース

## セットアップ

### 1. 環境変数の設定

```bash
cp env.example .env
```

`.env`ファイルを編集して、以下の設定を行います：

```env
# J-Quants API設定
JQUANTS_EMAIL=your_email@example.com
JQUANTS_PASSWORD=your_password

# データベース設定
DB_HOST=localhost
DB_PORT=3306
DB_USER=kabu_user
DB_PASSWORD=your_db_password
DB_NAME=kabu_analysis
```

### 2. 依存関係のインストール

```bash
make deps
```

### 3. ビルド

```bash
make build
```

## 使用方法

### データベースマイグレーション

```bash
# マイグレーション実行
./bin/migrate up

# マイグレーション確認
./bin/migrate version
```

### J-Quantsデータ取得

```bash
# 日次四本値データ取得
./bin/jquants daily-quotes --date 2024-01-01
```

### データクエリ

```bash
# データ一覧表示
./bin/sa query list

# 特定データ表示
./bin/sa query show --code 7203
```

## 利用可能なコマンド

### Makefileコマンド

- `make all` - 全ビルド（clean + deps + build）
- `make build` - 全バイナリビルド
- `make build-sa` - saコマンドビルド
- `make build-migrate` - migrateコマンドビルド
- `make clean` - 全ビルドファイル削除
- `make test` - 全プロジェクトテスト実行
- `make deps` - 依存関係ダウンロード
- `make update-deps` - 依存関係更新
- `make cross-build` - 複数プラットフォーム向けビルド

## 技術スタック

- **言語**: Go 1.24.5
- **データベース**: MySQL
- **ORM**: GORM
- **CLI**: Cobra
- **マイグレーション**: golang-migrate
- **環境変数**: godotenv
- **ログ**: slog

## データモデル

### 主要テーブル

- **`daily_quotes`** - 日次四本値データ
- **`listed_info`** - 上場銘柄情報
- **`financial_statements`** - 財務情報
- **`market_codes`** - 市場区分コード
- **`sector17_codes`** - 17業種コード
- **`sector33_codes`** - 33業種コード

## 開発

### Go Workspace

このプロジェクトはGo Workspaceを使用して複数モジュールを管理しています：

```bash
go work use ./sa ./jquants ./migrate ./database ./schema ./helper
```

### テスト実行

```bash
make test
```

### 依存関係更新

```bash
make update-deps
make tidy
```

## ライセンス

このプロジェクトのライセンス情報については、各ファイルのヘッダーを確認してください。

## 貢献

プルリクエストやイシューの報告は歓迎します。貢献する前に、既存のコードスタイルとテストカバレッジを確認してください。
