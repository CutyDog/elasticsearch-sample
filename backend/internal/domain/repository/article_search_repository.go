package repository

import (
	"elasticsearch-sample/backend/internal/domain/model"
)

type ArticleSearchRepository interface {
	// 記事インデックスを作成する
	CreateIndex() error
	// 記事を検索エンジンに保存する
	Index(article *model.Article) error
	// キーワードで記事を探す
	SimpleSearch(keyword string) ([]*model.Article, error)
}
