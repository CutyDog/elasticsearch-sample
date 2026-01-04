package repository

import (
	"elasticsearch-sample/backend/internal/domain/model"
)

type ArticleSearchRepository interface {
	// 記事インデックスを作成する
	CreateIndex() (string, error)
	// 記事インデックスを削除する
	DeleteIndex(indexName string) error
	// エイリアスを新しいインデックスに切り替える
	SwitchAlias(newIndexName string) error
	// 記事ドキュメントを保存する(作成・更新)
	Index(indexName *string, article *model.Article) error
	// 記事ドキュメントを一括保存する
	BulkIndex(indexName *string, articles []*model.Article) error
	// 記事ドキュメントを削除する
	Delete(id int64) error
	// キーワードで記事を探す
	SimpleSearch(keyword string) ([]*model.Article, error)
}
