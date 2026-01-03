.PHONY: migrate-all migrate-to migrate-status rollback-last rollback-to gqlgen-generate

# マイグレーションコマンドのビルド
migrate-build:
	go build -o migrate cmd/migrate/main.go

# 全てのマイグレーションを実行
migrate-all: migrate-build
	./migrate -all

# 指定したマイグレーションまで実行
migrate-to:
	@read -p "マイグレーション対象のマイグレーションIDを入力してください: " migration_id; \
	./migrate -to $$migration_id

# マイグレーション状態を確認
migrate-status: migrate-build
	./migrate -status

# ロールバックコマンドのビルド
rollback-build:
	go build -o rollback cmd/rollback/main.go

# 最後のマイグレーションをロールバック
rollback-last: rollback-build
	./rollback -last

# 指定したマイグレーションまでロールバック
rollback-to:
	@read -p "ロールバック対象のマイグレーションIDを入力してください: " migration_id; \
	./rollback -to $$migration_id

# GQLスキーマ生成
gqlgen-generate:
	go run github.com/99designs/gqlgen generate