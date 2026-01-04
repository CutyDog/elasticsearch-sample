package repository

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetArticleByID(ctx context.Context, id int64) (*model.Article, error)
	ListArticles(ctx context.Context, page, pageSize int) ([]*model.Article, error)

	CreateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	DeleteArticle(ctx context.Context, id int64) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetArticleByID(ctx context.Context, id int64) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) ListArticles(ctx context.Context, page, pageSize int) ([]*model.Article, error) {
	var articles []*model.Article
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) (*model.Article, error) {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error) {
	if err := r.db.WithContext(ctx).Save(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) DeleteArticle(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Article{}, id).Error
}
