package es

import (
	"context"
	"elasticsearch-sample/backend/internal/domain/model"
	"elasticsearch-sample/backend/internal/domain/repository"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type articleSearchRepo struct {
	client *Client
}

func NewArticleSearchRepository(client *Client) repository.ArticleSearchRepository {
	return &articleSearchRepo{client: client}
}

// 記事固有のインデックス名
const ArticleIndexName = "articles"

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
func (r *articleSearchRepo) CreateIndex() error {
	return r.client.CreateIndex(ArticleIndexName, &create.Request{
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
}

// Index: データの保存
func (r *articleSearchRepo) Index(article *model.Article) error {
	document := articleDocument{
		Id:        int(article.ID),
		Title:     article.Title,
		Content:   article.Content,
		Status:    article.Status,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: article.UpdatedAt.String(),
	}

	_, err := r.client.Typed.
		Index(ArticleIndexName).
		Id(strconv.Itoa(int(article.ID))).
		Request(document).
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
