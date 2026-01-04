package main

import (
	"context"
	"log"
	"time"

	"elasticsearch-sample/backend/internal/domain/repository"
	"elasticsearch-sample/backend/internal/infrastructure/db"
	"elasticsearch-sample/backend/internal/infrastructure/es"
	"elasticsearch-sample/backend/internal/usecase"
)

func main() {
	// タイムアウト付きのコンテキストを作成（10分間）
	// 大量データの移行を想定し、少し長めに設定します
	_, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Infrastructure初期化
	db.ConnectDB()
	esClient, err := es.NewClient()
	if err != nil {
		log.Fatalf("❌ Elasticsearchへの接続に失敗しました: %v", err)
	}

	// Repository初期化
	articleDBRepo := repository.NewArticleRepository(db.DB)
	articleSearchRepo := es.NewArticleSearchRepository(esClient)

	// Usecase初期化
	articleUsecase := usecase.NewArticleUsecase(articleDBRepo, articleSearchRepo)

	log.Println("🚀 検索エンジンの再構築を開始します...")

	// 再構築処理の実行
	if err := articleUsecase.ReindexSearchEngine(); err != nil {
		log.Fatalf("❌ 再構築中にエラーが発生しました: %v", err)
	}

	log.Println("✅ 検索エンジンの再構築が正常に完了しました。")
}
