# J-Quants API 上場銘柄取得プログラム

このプログラムは、J-Quants APIを使用して上場銘柄一覧を取得し、MySQLデータベースに保存するGoプログラムです。

## 機能

- J-Quants APIへの認証（リフレッシュトークン・IDトークン取得）
- 上場銘柄一覧の取得
- MySQLデータベースへの保存
- golang-migrateを使用したデータベーススキーマ管理
- 取得したデータのJSONファイルへのバックアップ保存
- 取得した銘柄情報の表示

## 必要な環境

- Go 1.21以上
- MySQL 5.7以上
- J-Quants APIアカウント

## セットアップ

1. 依存関係をインストール
```bash
go mod tidy
```

2. 環境変数を設定
```bash
cp env.example .env
# .envファイルを編集して実際の認証情報を設定
```

または、システム環境変数として設定：
```bash
export JQUANTS_EMAIL="your_email@example.com"
export JQUANTS_PASSWORD="your_password"
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="kabu_user"
export DB_PASSWORD="your_db_password"
export DB_NAME="kabu_analysis"
```

## データベースマイグレーション

### 初回セットアップ
```bash
# マイグレーションを実行してテーブルを作成
go run cmd/migrate/main.go -command=up
```

### マイグレーション管理
```bash
# マイグレーションを実行
go run cmd/migrate/main.go -command=up

# マイグレーションをロールバック
go run cmd/migrate/main.go -command=down

# 現在のバージョンを確認
go run cmd/migrate/main.go -command=version

# 特定のバージョンに強制設定
go run cmd/migrate/main.go -command=force -version=1
```

## 使用方法

```bash
# メインプログラムを実行
go run main.go

# または、シェルスクリプトを使用
./run.sh
```

## 出力

- コンソールに取得した銘柄数と最初の5件の情報を表示
- MySQLデータベースに全データを保存
- タイムスタンプ付きのJSONファイルにバックアップ保存

## データベーススキーマ

### listed_info テーブル
上場銘柄の基本情報を格納

### market_codes テーブル
市場区分コードと名称のマスターデータ

### sector17_codes テーブル
17業種コードと名称のマスターデータ

### sector33_codes テーブル
33業種コードと名称のマスターデータ

## API仕様

このプログラムは以下のJ-Quants APIエンドポイントを使用します：

- `/token/auth_user` - リフレッシュトークン取得
- `/token/auth_refresh` - IDトークン取得  
- `/listed/info` - 上場銘柄一覧取得

詳細は[J-Quants API仕様書](https://jpx.gitbook.io/j-quants-ja/api-reference)を参照してください。

## ディレクトリ構造

```
app/
├── main.go                    # メインプログラム
├── cmd/
│   └── migrate/
│       └── main.go           # マイグレーション管理ツール
├── migrations/                # マイグレーションファイル
│   ├── 000001_create_listed_info_table.up.sql
│   ├── 000001_create_listed_info_table.down.sql
│   ├── 000002_create_market_codes_table.up.sql
│   ├── 000002_create_market_codes_table.down.sql
│   ├── 000003_create_sector_codes_table.up.sql
│   └── 000003_create_sector_codes_table.down.sql
├── jquants/
│   ├── types.go              # データ型定義
│   ├── config.go             # 設定管理
│   ├── utils.go              # ユーティリティ関数
│   ├── client/               # APIクライアント
│   │   ├── client.go         # 統合クライアント
│   │   ├── auth.go           # 認証関連API
│   │   └── listed.go         # 上場銘柄関連API
│   └── database/             # データベース関連
│       ├── database.go       # データベース接続管理
│       ├── migrate.go        # マイグレーション管理
│       └── repository.go     # データ保存・取得処理
└── README.md                 # このファイル
```

## 注意事項

- API利用制限にご注意ください
- 認証情報は適切に管理してください
- 取得したデータの利用はJ-Quantsの利用規約に従ってください
- データベースのバックアップを定期的に取得してください
