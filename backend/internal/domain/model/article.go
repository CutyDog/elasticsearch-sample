package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index:idx_user_on_articles"`
	Title   string `gorm:"not null;size:255"`
	Content string
	Status  string `gorm:"not null;default:draft"` // draft, published, archived

	Author User `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}
