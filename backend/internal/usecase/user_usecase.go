package usecase

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/domain/repository"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByUID(ctx context.Context, userUID string) (*model.User, error)

	SeedUsers() ([]model.User, error)
}

type userUsecase struct {
	dbRepo repository.UserRepository
}

func NewUserUsecase(dbRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		dbRepo: dbRepo,
	}
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := u.dbRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetUserByUID(ctx context.Context, userUID string) (*model.User, error) {
	user, err := u.dbRepo.GetUserByUID(ctx, userUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) SeedUsers() ([]model.User, error) {
	seedUsers := []model.User{
		{UID: "admin123"},
		{UID: "admin456"},
	}

	var users []model.User
	for _, user := range seedUsers {
		// 既に存在する場合はスキップ
		existUser, _ := u.dbRepo.GetUserByUID(context.Background(), user.UID)
		if existUser != nil {
			users = append(users, *existUser)
			continue
		}

		// 新規作成
		if err := u.dbRepo.CreateUser(context.Background(), &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
