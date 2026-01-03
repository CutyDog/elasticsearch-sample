package db

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var dsn string

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	} else {
		// 個別の環境変数を使用
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("❌ データベース接続に失敗しました: " + err.Error())
	}

	DB = db
	log.Printf("✅ データベース接続が成功しました")
}

// RollbackTo 指定したマイグレーションIDまでロールバックする
func RollbackTo(migrationID string) error {
	if DB == nil {
		return fmt.Errorf("データベースが接続されていません")
	}

	// マイグレーションIDの存在確認
	migrations := getMigrations()
	migrationExists := false
	for _, migration := range migrations {
		if migration.ID == migrationID {
			migrationExists = true
			break
		}
	}
	if !migrationExists {
		return fmt.Errorf("マイグレーションID '%s' が見つかりません", migrationID)
	}

	m := gormigrate.New(DB, gormigrate.DefaultOptions, migrations)

	log.Printf("🔄 マイグレーション '%s' までロールバックを開始します...", migrationID)

	if err := m.RollbackTo(migrationID); err != nil {
		log.Printf("❌ ロールバックに失敗しました: %v", err)
		return fmt.Errorf("ロールバックに失敗しました: %w", err)
	}

	log.Printf("✅ マイグレーション '%s' までロールバックしました", migrationID)
	return nil
}

// RollbackLast 最後のマイグレーションをロールバックする
func RollbackLast() error {
	if DB == nil {
		return fmt.Errorf("データベースが接続されていません")
	}

	migrations := getMigrations()
	if len(migrations) == 0 {
		return fmt.Errorf("マイグレーションが定義されていません")
	}

	m := gormigrate.New(DB, gormigrate.DefaultOptions, migrations)

	// 最後に実行されたマイグレーションを確認
	lastMigration, err := getLastRunMigration()
	if err != nil {
		return fmt.Errorf("最後のマイグレーションの確認に失敗しました: %w", err)
	}
	if lastMigration == "" {
		return fmt.Errorf("実行されたマイグレーションがありません")
	}

	log.Printf("🔄 最後のマイグレーション '%s' をロールバックします...", lastMigration)

	if err := m.RollbackLast(); err != nil {
		log.Printf("❌ 最後のマイグレーションのロールバックに失敗しました: %v", err)
		return fmt.Errorf("最後のマイグレーションのロールバックに失敗しました: %w", err)
	}

	log.Printf("✅ 最後のマイグレーション '%s' をロールバックしました", lastMigration)
	return nil
}

// getLastRunMigration 最後に実行されたマイグレーションIDを取得する
func getLastRunMigration() (string, error) {
	if DB == nil {
		return "", fmt.Errorf("データベースが接続されていません")
	}

	var migration struct {
		ID string `gorm:"column:id"`
	}

	err := DB.Table("migrations").Order("id DESC").First(&migration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil // マイグレーションが実行されていない
		}
		return "", err
	}

	return migration.ID, nil
}

// GetMigrationStatus マイグレーションの状態を取得する
func GetMigrationStatus() ([]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("データベースが接続されていません")
	}

	var migrations []struct {
		ID string `gorm:"column:id"`
	}

	err := DB.Table("migrations").Order("id ASC").Find(&migrations).Error
	if err != nil {
		return nil, fmt.Errorf("マイグレーション状態の取得に失敗しました: %w", err)
	}

	var executedMigrations []string
	for _, migration := range migrations {
		executedMigrations = append(executedMigrations, migration.ID)
	}

	return executedMigrations, nil
}

// MigrateTo 指定したマイグレーションIDまでマイグレーションを実行する
func MigrateTo(migrationID string) error {
	if DB == nil {
		return fmt.Errorf("データベースが接続されていません")
	}

	// マイグレーションIDの存在確認
	migrations := getMigrations()
	migrationExists := false
	for _, migration := range migrations {
		if migration.ID == migrationID {
			migrationExists = true
			break
		}
	}
	if !migrationExists {
		return fmt.Errorf("マイグレーションID '%s' が見つかりません", migrationID)
	}

	m := gormigrate.New(DB, gormigrate.DefaultOptions, migrations)

	log.Printf("🔄 マイグレーション '%s' まで実行を開始します...", migrationID)

	if err := m.MigrateTo(migrationID); err != nil {
		log.Printf("❌ マイグレーションに失敗しました: %v", err)
		return fmt.Errorf("マイグレーションに失敗しました: %w", err)
	}

	log.Printf("✅ マイグレーション '%s' まで実行しました", migrationID)
	return nil
}

// RunMigrations マイグレーションを実行する（サーバー起動時には呼ばれない）
func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("データベースが接続されていません")
	}

	migrations := getMigrations()
	if len(migrations) == 0 {
		return fmt.Errorf("マイグレーションが定義されていません")
	}

	m := gormigrate.New(DB, gormigrate.DefaultOptions, migrations)

	log.Printf("🔄 マイグレーションを開始します...")
	log.Printf("📋 実行予定のマイグレーション数: %d", len(migrations))

	if err := m.Migrate(); err != nil {
		log.Printf("❌ マイグレーションに失敗しました: %v", err)
		return fmt.Errorf("マイグレーションに失敗しました: %w", err)
	}

	log.Printf("✅ マイグレーションが正常に実行されました")
	return nil
}
