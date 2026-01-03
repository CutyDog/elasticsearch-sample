package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"elasticsearch-sample/backend/internal/db"
)

func main() {
	// コマンドライン引数の解析
	var (
		rollbackTo   = flag.String("to", "", "指定したマイグレーションIDまでロールバックする")
		rollbackLast = flag.Bool("last", false, "最後のマイグレーションをロールバックする")
		showStatus   = flag.Bool("status", false, "マイグレーションの状態を表示する")
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

	// ロールバックの実行
	if *rollbackTo != "" {
		if err := db.RollbackTo(*rollbackTo); err != nil {
			log.Printf("❌ ロールバックに失敗しました: %v", err)
			os.Exit(1)
		}
		log.Printf("✅ マイグレーション '%s' までロールバックが完了しました", *rollbackTo)
	} else if *rollbackLast {
		if err := db.RollbackLast(); err != nil {
			log.Printf("❌ ロールバックに失敗しました: %v", err)
			os.Exit(1)
		}
		log.Printf("✅ 最後のマイグレーションのロールバックが完了しました")
	} else {
		fmt.Println("❌ ロールバック対象を指定してください")
		fmt.Println("使用例:")
		fmt.Println("  ./rollback -to 202601021300_create_users")
		fmt.Println("  ./rollback -last")
		fmt.Println("  ./rollback -status")
		os.Exit(1)
	}
}
