# Stock Automation プロジェクト統合Makefile
PROJECT_NAME = stock-automation
BUILD_DIR = bin

# Goコマンド設定
GO = go
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOTEST = $(GO) test
GOGET = $(GO) get
GOMOD = $(GO) mod

# ビルドフラグ
LDFLAGS = -ldflags "-s -w"
BUILD_FLAGS = $(LDFLAGS)

# ターゲットプラットフォーム
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build build-sa build-migrate clean clean-sa clean-migrate test deps update-deps tidy help cross-build info

# デフォルトターゲット
all: clean deps build

# ヘルプ
help:
	@echo "Stock Automation プロジェクト - 利用可能なコマンド:"
	@echo "  all          - clean + deps + build (sa + migrate)"
	@echo "  build        - 全バイナリをビルド (sa + migrate)"
	@echo "  build-sa     - saコマンドをビルド"
	@echo "  build-migrate - migrateコマンドをビルド"
	@echo "  clean        - 全ビルドファイルを削除"
	@echo "  clean-sa     - saコマンドのビルドファイルを削除"
	@echo "  clean-migrate - migrateコマンドのビルドファイルを削除"
	@echo "  test         - 全プロジェクトのテストを実行"
	@echo "  deps         - 全プロジェクトの依存関係をダウンロード"
	@echo "  update-deps  - 全プロジェクトの依存関係を更新"
	@echo "  tidy         - 全プロジェクトのgo.modを整理"
	@echo "  cross-build  - 複数プラットフォーム向けビルド"
	@echo "  info         - プロジェクト情報を表示"

# ビルドディレクトリ作成
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 全バイナリビルド
build: build-sa build-migrate

# saコマンドビルド
build-sa: $(BUILD_DIR)
	@echo "saコマンドをビルド中..."
	cd sa && $(GOBUILD) $(BUILD_FLAGS) -o ../$(BUILD_DIR)/sa main.go
	@echo "saコマンドビルド完了: $(BUILD_DIR)/sa"

# migrateコマンドビルド
build-migrate: $(BUILD_DIR)
	@echo "migrateコマンドをビルド中..."
	cd migrate && $(GOBUILD) $(BUILD_FLAGS) -o ../$(BUILD_DIR)/migrate .
	@echo "migrateコマンドビルド完了: $(BUILD_DIR)/migrate"

# 全クリーンアップ
clean: clean-sa clean-migrate
	@echo "ルートビルドディレクトリをクリーンアップ中..."
	rm -rf $(BUILD_DIR)
	@echo "全クリーンアップ完了"

# saコマンドクリーンアップ
clean-sa:
	@echo "saコマンドをクリーンアップ中..."
	cd sa && $(GOCLEAN)
	rm -rf sa/bin
	@echo "saコマンドクリーンアップ完了"

# migrateコマンドクリーンアップ
clean-migrate:
	@echo "migrateコマンドをクリーンアップ中..."
	cd migrate && $(GOCLEAN)
	@echo "migrateコマンドクリーンアップ完了"

# 全テスト実行
test:
	@echo "全プロジェクトのテストを実行中..."
	cd sa && $(GOTEST) -v ./...
	cd migrate && $(GOTEST) -v ./...
	cd database && $(GOTEST) -v ./...
	cd schema && $(GOTEST) -v ./...
	@echo "全テスト完了"

# 全依存関係ダウンロード
deps:
	@echo "全プロジェクトの依存関係をダウンロード中..."
	cd sa && $(GOMOD) download
	cd migrate && $(GOMOD) download
	cd database && $(GOMOD) download
	cd schema && $(GOMOD) download
	@echo "全依存関係ダウンロード完了"

# 全依存関係更新
update-deps:
	@echo "全プロジェクトの依存関係を更新中..."
	cd sa && $(GOGET) -u ./... && $(GOMOD) tidy
	cd migrate && $(GOGET) -u ./... && $(GOMOD) tidy
	cd database && $(GOGET) -u ./... && $(GOMOD) tidy
	cd schema && $(GOGET) -u ./... && $(GOMOD) tidy
	@echo "全依存関係更新完了"

# 全go.mod整理
tidy:
	@echo "全プロジェクトのgo.modを整理中..."
	cd sa && $(GOMOD) tidy
	cd migrate && $(GOMOD) tidy
	cd database && $(GOMOD) tidy
	cd schema && $(GOMOD) tidy
	@echo "全go.mod整理完了"

# クロスプラットフォームビルド
cross-build: $(BUILD_DIR)
	@echo "複数プラットフォーム向けビルド中..."
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		SA_OUTPUT=$(BUILD_DIR)/sa-$$OS-$$ARCH; \
		MIGRATE_OUTPUT=$(BUILD_DIR)/migrate-$$OS-$$ARCH; \
		if [ $$OS = "windows" ]; then \
			SA_OUTPUT=$$SA_OUTPUT.exe; \
			MIGRATE_OUTPUT=$$MIGRATE_OUTPUT.exe; \
		fi; \
		echo "ビルド中 ($$OS/$$ARCH): $$SA_OUTPUT, $$MIGRATE_OUTPUT"; \
		cd sa && GOOS=$$OS GOARCH=$$ARCH $(GOBUILD) $(BUILD_FLAGS) -o ../$$SA_OUTPUT main.go; \
		cd ../migrate && GOOS=$$OS GOARCH=$$ARCH $(GOBUILD) $(BUILD_FLAGS) -o ../$$MIGRATE_OUTPUT main.go; \
		cd ..; \
	done
	@echo "クロスプラットフォームビルド完了"

# プロジェクト情報表示
info:
	@echo "Stock Automation プロジェクト情報:"
	@echo "  プロジェクト名: $(PROJECT_NAME)"
	@echo "  ビルドディレクトリ: $(BUILD_DIR)"
	@echo "  Goバージョン: $$(go version)"
	@echo ""
	@echo "サブプロジェクト:"
	@echo "  sa/        - メインCLIツール"
	@echo "  migrate/   - データベースマイグレーションツール"
	@echo "  database/  - データベース共通ライブラリ"
	@echo "  schema/    - スキーマ定義ライブラリ"
