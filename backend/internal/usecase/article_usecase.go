package usecase

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/domain/repository"
)

type ArticleUsecase struct {
	dbRepo     repository.ArticleRepository
	searchRepo repository.ArticleSearchRepository
}

func NewArticleUsecase(
	dbRepo repository.ArticleRepository,
	searchRepo repository.ArticleSearchRepository,
) *ArticleUsecase {
	return &ArticleUsecase{
		dbRepo:     dbRepo,
		searchRepo: searchRepo,
	}
}

func (u *ArticleUsecase) CreateArticle(
	ctx context.Context,
	title,
	content string,
	status string,
	userID uint,
) (*model.Article, error) {
	// DBに記事作成
	article := &model.Article{
		Title:   title,
		Content: content,
		Status:  status,
		UserID:  userID,
	}
	createdArticle, err := u.dbRepo.CreateArticle(ctx, article)
	if err != nil {
		return nil, err
	}

	// 検索エンジンにインデックス登録
	err = u.searchRepo.Index(createdArticle)
	if err != nil {
		return nil, err
	}

	return createdArticle, nil
}
