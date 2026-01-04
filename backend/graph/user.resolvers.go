package graph

import (
	"context"
	"fmt"

	model "elasticsearch-sample/backend/graph/model"
	entity "elasticsearch-sample/backend/internal/domain/model"
)

func (r *Resolver) User() ArticleResolver { return &articleResolver{r} }

type userResolver struct{ *Resolver }

// ========================
// 汎用関数
// ========================
func ToModelUser(user entity.User) *model.User {
	return &model.User{
		ID:  fmt.Sprintf("%d", user.ID),
		UID: user.UID,
	}
}

// ========================
// Resolver
// ========================

// ========================
// Mutation
// ========================

// ========================
// Query
// ========================

func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	// ログインユーザーIDの取得
	userUID, ok := GetUserUIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	// Usecaseの呼び出し
	user, err := r.UserUsecase.GetUserByUID(ctx, userUID)
	if err != nil {
		return nil, err
	}

	return ToModelUser(*user), nil
}
