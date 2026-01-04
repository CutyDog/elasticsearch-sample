package repository

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByUID(ctx context.Context, uid string) (*model.User, error)

	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUID(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
