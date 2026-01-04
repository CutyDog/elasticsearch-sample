package es

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/domain/repository"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

const (
	ArticleIndexName = "articles" // 記事インデックス名(エイリアスとして付与される)
)

type articleSearchRepo struct {
	client *Client
}

func NewArticleSearchRepository(client *Client) repository.ArticleSearchRepository {
	return &articleSearchRepo{client: client}
}

// articleDocument: ESに保存する記事ドキュメントの構造体
type articleDocument struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ConvertToModel: ESのレスポンスをドメインモデルに変換
func ConvertToModel(data articleDocument) *model.Article {
	article := &model.Article{}

	article.ID = uint(data.Id)
	article.Title = data.Title
	article.Content = data.Content
	article.Status = data.Status
	article.CreatedAt, _ = time.Parse(time.RFC3339, data.CreatedAt)
	article.UpdatedAt, _ = time.Parse(time.RFC3339, data.UpdatedAt)

	return article
}

// Search: 検索クエリ実行
func Search(es *Client, req *search.Request) ([]*model.Article, error) {
	res, err := es.Typed.
		Search().
		Index(ArticleIndexName).
		Request(req).
		Do(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	articles := []*model.Article{}
	for _, hit := range res.Hits.Hits {
		var doc articleDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			continue
		}
		article := ConvertToModel(doc)
		articles = append(articles, article)
	}

	return articles, nil
}

// CreateIndex: インデックス作成
func (r *articleSearchRepo) CreateIndex() (string, error) {
	newIndexName := fmt.Sprintf("article_%s", time.Now().Format("200601021504"))

	err := r.client.CreateIndex(newIndexName, &create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				"id":         types.NewKeywordProperty(),
				"title":      types.NewTextProperty(),
				"content":    types.NewTextProperty(),
				"status":     types.NewKeywordProperty(),
				"created_at": types.NewDateProperty(),
				"updated_at": types.NewDateProperty(),
			},
		},
	})
	if err != nil {
		log.Println("❌ インデックス作成失敗:", err)
	}
	return newIndexName, err
}

// DeleteIndex: インデックス削除
func (r *articleSearchRepo) DeleteIndex(indexName string) error {
	err := r.client.DeleteIndex(indexName)
	if err != nil {
		log.Println("❌ インデックス削除失敗:", err)
	}
	return err
}

// SwitchAlias: エイリアス切り替え
func (r *articleSearchRepo) SwitchAlias(newIndexName string) error {
	err := r.client.UpdateAlias(ArticleIndexName, newIndexName)
	if err != nil {
		log.Println("❌ エイリアス更新失敗:", err)
	}
	return err
}

// Index: データの保存
func (r *articleSearchRepo) Index(indexName *string, article *model.Article) error {
	// インデックス名の指定がなければ、エイリアスが付与されてるインデックスを使用
	index := ArticleIndexName
	if indexName != nil {
		index = *indexName
	}

	document := articleDocument{
		Id:        int(article.ID),
		Title:     article.Title,
		Content:   article.Content,
		Status:    article.Status,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: article.UpdatedAt.String(),
	}

	_, err := r.client.Typed.
		Index(index).
		Id(strconv.Itoa(int(article.ID))).
		Request(document).
		Do(context.Background())

	return err
}

// BulkIndex: データの一括保存
func (r *articleSearchRepo) BulkIndex(indexName *string, articles []*model.Article) error {
	// インデックス名の指定がなければ、エイリアスが付与されてるインデックスを使用
	index := ArticleIndexName
	if indexName != nil {
		index = *indexName
	}

	bulkRequest := r.client.Typed.Bulk()

	for _, article := range articles {
		document := articleDocument{
			Id:        int(article.ID),
			Title:     article.Title,
			Content:   article.Content,
			Status:    article.Status,
			CreatedAt: article.CreatedAt.Format(time.RFC3339),
			UpdatedAt: article.UpdatedAt.Format(time.RFC3339),
		}

		id := strconv.Itoa(int(article.ID))
		err := bulkRequest.IndexOp(
			types.IndexOperation{
				Id_:    &id,
				Index_: &index,
			},
			document,
		)
		if err != nil {
			return fmt.Errorf("failed to add operation: %w", err)
		}
	}

	res, err := bulkRequest.Do(context.Background())
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}

	if res.Errors {
		errorDetails := ""
		for _, item := range res.Items {
			for op, result := range item {
				if result.Error != nil {
					errorReason := ""
					if result.Error.Reason != nil {
						errorReason = *result.Error.Reason
					}
					errorDetails += fmt.Sprintf("Operation: %s, Status: %d, Error: %s\n", op, result.Status, errorReason)
				}
			}
		}
		return fmt.Errorf("bulk request had item errors:\n%s", errorDetails)
	}

	return nil
}

// Delete: ドキュメント削除
func (r *articleSearchRepo) Delete(id int64) error {
	articleId := strconv.Itoa(int(id))
	_, err := r.client.Typed.
		Delete(ArticleIndexName, articleId).
		Do(context.Background())

	return err
}

// SimpleSearch: キーワード検索
func (r *articleSearchRepo) SimpleSearch(keyword string) ([]*model.Article, error) {
	req := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						// 公開ステータスの絞り込み
						Term: map[string]types.TermQuery{
							"status": {Value: "published"},
						},
					},
					{
						// キーワード検索
						MultiMatch: &types.MultiMatchQuery{
							Query:  keyword,
							Fields: []string{"title", "content"},
						},
					},
				},
			},
		},
	}

	return Search(r.client, req)
}
