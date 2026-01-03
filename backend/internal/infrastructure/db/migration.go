package db

import (
	"elasticsearch-sample/backend/internal/domain/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// getMigrations マイグレーションの定義を取得する
func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202601021300_create_users",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&model.User{})
			},
		},
		{
			ID: "202601021310_create_articles",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&model.Article{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&model.Article{})
			},
		},
	}
}
