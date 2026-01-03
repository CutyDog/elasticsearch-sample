package es

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
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

func (c *Client) CreateIndex(name string, req *create.Request) error {
	req.Settings = buildGlobalSettings()

	exists, _ := c.Typed.Indices.Exists(name).Do(context.Background())
	if exists {
		return nil
	}

	_, err := c.Typed.Indices.Create(name).Request(req).Do(context.Background())
	return err
}
