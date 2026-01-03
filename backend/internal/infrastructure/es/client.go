package es

import (
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v9"
)

var Elasticsearch *elasticsearch.TypedClient

func ConnectElasticsearch() {
	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ES_HOST")},
		CloudID:   os.Getenv("ES_CLOUD_ID"),
		APIKey:    os.Getenv("ES_API_KEY"),
	})

	if err != nil {
		panic("❌ Elasticsearchへの接続に失敗しました: " + err.Error())
	}

	log.Printf("✅ Elasticsearchへの接続が成功しました")

	Elasticsearch = typedClient
}
