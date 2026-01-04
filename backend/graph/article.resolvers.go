package graph

import (
	"context"
	"fmt"
	"strconv"

	model "elasticsearch-sample/backend/graph/model"
	entity "elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/usecase"
)

func (r *Resolver) Article() ArticleResolver { return &articleResolver{r} }

type articleResolver struct{ *Resolver }

// ========================
// 汎用関数
// ========================
func ToModelArticle(article entity.Article) *model.Article {
	return &model.Article{
		ID:      fmt.Sprintf("%d", article.ID),
		Title:   article.Title,
		Content: &article.Content,
		Status:  model.ArticleStatus(article.Status),
	}
}

// =======================
// Resolver
// ========================
func (r *articleResolver) Author(ctx context.Context, obj *model.Article) (*model.User, error) {
	userID, err := strconv.ParseUint(obj.UserID, 10, 32)
	if err != nil {
		return nil, err
	}

	user, err := r.UserUsecase.GetUserByID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	return ToModelUser(*user), nil
}

// ========================
// Mutation
// ========================

func (r *mutationResolver) CreateArticle(ctx context.Context, input model.CreateArticleInput) (*model.Article, error) {
	// ログインユーザーIDの取得
	userUID, ok := GetUserUIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}
	// ログインユーザー情報の取得
	user, err := r.UserUsecase.GetUserByUID(ctx, userUID)
	if err != nil {
		return nil, err
	}

	// Usecaseの呼び出し
	createInput := usecase.CreateArticleInput{
		Title:   input.Title,
		Content: *input.Content,
		Status:  "draft",
		UserID:  user.ID,
	}
	createdArticle, err := r.ArticleUsecase.CreateArticle(ctx, createInput)
	if err != nil {
		return nil, err
	}

	return ToModelArticle(*createdArticle), nil
}

func (r *mutationResolver) PublishArticle(ctx context.Context, input model.PublishArticleInput) (*model.Article, error) {
	// ログインユーザーIDの取得
	_, ok := GetUserUIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	// Usecaseの呼び出し
	articleID, err := strconv.ParseUint(input.ID, 10, 32)
	if err != nil {
		return nil, err
	}
	status := "published"
	updatedArticle, err := r.ArticleUsecase.UpdateArticle(ctx, usecase.UpdateArticleInput{
		ArticleID: uint(articleID),
		Status:    &status,
	})
	if err != nil {
		return nil, err
	}

	return ToModelArticle(*updatedArticle), nil
}

func (r *mutationResolver) ArchiveArticle(ctx context.Context, input model.ArchiveArticleInput) (*model.Article, error) {
	// ログインユーザーIDの取得
	_, ok := GetUserUIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	// Usecaseの呼び出し
	articleID, err := strconv.ParseUint(input.ID, 10, 32)
	if err != nil {
		return nil, err
	}
	status := "archived"
	updatedArticle, err := r.ArticleUsecase.UpdateArticle(ctx, usecase.UpdateArticleInput{
		ArticleID: uint(articleID),
		Status:    &status,
	})
	if err != nil {
		return nil, err
	}

	return ToModelArticle(*updatedArticle), nil
}

// ========================
// Query
// ========================

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	// Usecaseの呼び出し
	articles := []*model.Article{}
	var page int = 1
	var pageSize int = 100
	for {
		fetchedArticles, err := r.ArticleUsecase.ListArticles(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		if len(fetchedArticles) == 0 {
			break
		}
		for _, article := range fetchedArticles {
			articles = append(articles, ToModelArticle(*article))
		}
		page++
	}

	return articles, nil
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	// Usecaseの呼び出し
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}
	article, err := r.ArticleUsecase.GetArticleByID(ctx, uint(articleID))
	if err != nil {
		return nil, err
	}

	return ToModelArticle(*article), nil
}
