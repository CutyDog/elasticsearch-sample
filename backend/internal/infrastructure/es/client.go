package es

import (
	"context"
	"fmt"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v9"
)

type Client struct {
	Typed *elasticsearch.TypedClient
}

// NewClient は環境変数から設定を読み込み、ESクライアントを初期化します
func NewClient() (*Client, error) {
	// 接続設定
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ES_HOST")},
		CloudID:   os.Getenv("ES_CLOUD_ID"),
		APIKey:    os.Getenv("ES_API_KEY"),
	}

	// v9 の TypedClient を作成
	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("❌ Elasticsearch typed client作成失敗: %w", err)
	}

	// 疎通確認 (Ping)
	// v9 では Ping() もメソッドチェーンで呼び出せます
	_, err = typedClient.Info().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("❌ Elasticsearch接続失敗: %w", err)
	}

	log.Println("✅ Elasticsearch接続が成功しました")

	return &Client{Typed: typedClient}, nil
}
