package es

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/updatealiases"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// 共通の日本語解析設定
func buildGlobalSettings() *types.IndexSettings {
	return &types.IndexSettings{
		Analysis: &types.IndexSettingsAnalysis{
			Analyzer: map[string]types.Analyzer{
				"ja_analyzer": types.CustomAnalyzer{
					Type:      "custom",
					Tokenizer: "kuromoji_tokenizer",
					Filter:    []string{"kuromoji_baseform", "lowercase", "icu_normalizer"},
				},
			},
		},
	}
}

// CreateIndex: インデックスを作成
func (c *Client) CreateIndex(name string, req *create.Request) error {
	req.Settings = buildGlobalSettings()

	exists, _ := c.Typed.Indices.Exists(name).Do(context.Background())
	if exists {
		return nil
	}

	_, err := c.Typed.Indices.Create(name).Request(req).Do(context.Background())
	return err
}

// UpdateAlias: エイリアスの張り替えをアトミックに実行
func (c *Client) UpdateAlias(aliasName string, newIndexName string) error {
	// // 1. 現在そのエイリアスが付与されているインデックスを特定する
	res, err := c.Typed.Indices.GetAlias().Name(aliasName).Do(context.Background())
	if err != nil {
		// まだエイリアスが存在しない場合は新規作成
		_, err = c.Typed.Indices.PutAlias(newIndexName, aliasName).Do(context.Background())
		return err
	}

	// 2. エイリアスの張り替えアクションを準備
	actions := []types.IndicesAction{}

	// 3. 既存のエイリアスを削除するアクションを追加
	for oldIndex := range res {
		actions = append(actions, types.IndicesAction{
			Remove: &types.RemoveAction{
				Index: &oldIndex,
				Alias: &aliasName,
			},
		})
	}

	// 4. 新しいインデックスにエイリアスを追加するアクションを追加
	actions = append(actions, types.IndicesAction{
		Add: &types.AddAction{
			Index: &newIndexName,
			Alias: &aliasName,
		},
	})

	// 5. アクションを実行してエイリアスを更新
	_, err = c.Typed.Indices.UpdateAliases().Request(&updatealiases.Request{Actions: actions}).Do(context.Background())
	return err
}

// DeleteIndex: インデックスを削除
func (c *Client) DeleteIndex(name string) error {
	_, err := c.Typed.Indices.Delete(name).Do(context.Background())
	return err
}
