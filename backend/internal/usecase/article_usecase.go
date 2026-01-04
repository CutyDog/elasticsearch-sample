package usecase

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/domain/repository"
)

type CreateArticleInput struct {
	Title   string
	Content string
	Status  string
	UserID  uint
}

type UpdateArticleInput struct {
	ArticleID uint
	Title     *string
	Content   *string
	Status    *string
}

type ArticleUsecase interface {
	GetArticleByID(ctx context.Context, articleID uint) (*model.Article, error)
	ListArticles(ctx context.Context, page int, pageSize int) ([]*model.Article, error)
	SearchArticles(ctx context.Context, keyword string) ([]*model.Article, error)

	CreateArticle(ctx context.Context, input CreateArticleInput) (*model.Article, error)
	UpdateArticle(ctx context.Context, input UpdateArticleInput) (*model.Article, error)
	DeleteArticle(ctx context.Context, articleID uint) error

	ReindexSearchEngine() error

	SeedArticles(userID uint) ([]model.Article, error)
}

type articleUsecase struct {
	dbRepo     repository.ArticleRepository
	searchRepo repository.ArticleSearchRepository
}

func NewArticleUsecase(dbRepo repository.ArticleRepository, searchRepo repository.ArticleSearchRepository) ArticleUsecase {
	return &articleUsecase{
		dbRepo:     dbRepo,
		searchRepo: searchRepo,
	}
}

// GetArticleByID: IDで記事取得
func (u *articleUsecase) GetArticleByID(ctx context.Context, articleID uint) (*model.Article, error) {
	article, err := u.dbRepo.GetArticleByID(ctx, int64(articleID))
	if err != nil {
		return nil, err
	}
	return article, nil
}

// ListArticles: 記事一覧取得
func (u *articleUsecase) ListArticles(ctx context.Context, page int, pageSize int) ([]*model.Article, error) {
	articles, err := u.dbRepo.ListArticles(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// SearchArticles: キーワードで記事検索
func (u *articleUsecase) SearchArticles(ctx context.Context, keyword string) ([]*model.Article, error) {
	articles, err := u.searchRepo.SimpleSearch(keyword)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// CreateArticle: 記事作成
func (u *articleUsecase) CreateArticle(ctx context.Context, input CreateArticleInput) (*model.Article, error) {
	// DBに記事作成
	article := &model.Article{
		Title:   input.Title,
		Content: input.Content,
		Status:  input.Status,
		UserID:  input.UserID,
	}
	createdArticle, err := u.dbRepo.CreateArticle(ctx, article)
	if err != nil {
		return nil, err
	}

	// 検索エンジンにインデックス登録
	err = u.searchRepo.Index(nil, createdArticle)
	if err != nil {
		return nil, err
	}

	return createdArticle, nil
}

// UpdateArticle: 記事更新
func (u *articleUsecase) UpdateArticle(ctx context.Context, input UpdateArticleInput) (*model.Article, error) {
	article, err := u.dbRepo.GetArticleByID(ctx, int64(input.ArticleID))
	if err != nil {
		return nil, err
	}

	// DBで記事更新
	hasChanged := false
	if input.Title != nil {
		article.Title = *input.Title
		hasChanged = true
	}
	if input.Content != nil {
		article.Content = *input.Content
		hasChanged = true
	}
	if input.Status != nil {
		article.Status = *input.Status
		hasChanged = true
	}

	if !hasChanged {
		return article, nil
	}
	updatedArticle, err := u.dbRepo.UpdateArticle(ctx, article)
	if err != nil {
		return nil, err
	}

	// 検索エンジンで記事更新
	err = u.searchRepo.Index(nil, updatedArticle)
	if err != nil {
		return nil, err
	}

	return updatedArticle, nil
}

// DeleteArticle: 記事削除
func (u *articleUsecase) DeleteArticle(ctx context.Context, articleID uint) error {
	id := int64(articleID)

	// DBから記事削除
	err := u.dbRepo.DeleteArticle(ctx, id)
	if err != nil {
		return err
	}

	// 検索エンジンから記事削除
	err = u.searchRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// ReindexSearchEngine: 検索エンジンのインデックス再構築
func (u *articleUsecase) ReindexSearchEngine() error {
	// 新しいインデックスを作成
	newIndexName, err := u.searchRepo.CreateIndex()
	if err != nil {
		return err
	}

	// 既存のDB記事をすべて取得
	var page int = 1
	var pageSize int = 100
	for {
		batch, err := u.dbRepo.ListArticles(context.Background(), page, pageSize)
		if err != nil {
			return err
		}
		if len(batch) == 0 {
			break
		}

		// 指定した新しいインデックスに対して一括登録
		if err := u.searchRepo.BulkIndex(&newIndexName, batch); err != nil {
			return err
		}
		page++
	}

	// エイリアスを新しいインデックスに切り替え
	err = u.searchRepo.SwitchAlias(newIndexName)
	if err != nil {
		return err
	}

	return nil
}

func (u *articleUsecase) SeedArticles(userID uint) ([]model.Article, error) {
	var seedArticles []model.Article = []model.Article{
		{Title: "First Article", Content: "This is the content of the first article.", Status: "published", UserID: userID},
	}

	var articles []model.Article
	for _, article := range seedArticles {
		createdArticle, err := u.dbRepo.CreateArticle(context.Background(), &article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, *createdArticle)
	}
	return articles, nil
}
