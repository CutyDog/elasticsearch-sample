package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"elasticsearch-sample/backend/internal/infrastructure/db"
)

func main() {
	// コマンドライン引数の解析
	var (
		migrateTo  = flag.String("to", "", "指定したマイグレーションIDまでマイグレーションを実行する")
		migrateAll = flag.Bool("all", false, "全てのマイグレーションを実行する")
		showStatus = flag.Bool("status", false, "マイグレーションの状態を表示する")
	)
	flag.Parse()

	// データベースに接続
	db.ConnectDB()

	// マイグレーション状態の表示
	if *showStatus {
		executedMigrations, err := db.GetMigrationStatus()
		if err != nil {
			log.Printf("❌ マイグレーション状態の取得に失敗しました: %v", err)
			os.Exit(1)
		}

		if len(executedMigrations) == 0 {
			fmt.Println("📋 実行されたマイグレーションはありません")
		} else {
			fmt.Printf("📋 実行済みマイグレーション (%d件):\n", len(executedMigrations))
			for i, migrationID := range executedMigrations {
				fmt.Printf("  %d. %s\n", i+1, migrationID)
			}
		}
		return
	}

	// マイグレーションの実行
	if *migrateTo != "" {
		if err := db.MigrateTo(*migrateTo); err != nil {
			log.Printf("❌ マイグレーションに失敗しました: %v", err)
			os.Exit(1)
		}
		log.Printf("✅ マイグレーション '%s' まで実行が完了しました", *migrateTo)
	} else if *migrateAll {
		if err := db.RunMigrations(); err != nil {
			log.Printf("❌ マイグレーションに失敗しました: %v", err)
			os.Exit(1)
		}
		log.Printf("✅ 全てのマイグレーションの実行が完了しました")
	} else {
		fmt.Println("❌ マイグレーション対象を指定してください")
		fmt.Println("使用例:")
		fmt.Println("  ./migrate -to 202601021300_create_users")
		fmt.Println("  ./migrate -all")
		fmt.Println("  ./migrate -status")
		os.Exit(1)
	}
}
