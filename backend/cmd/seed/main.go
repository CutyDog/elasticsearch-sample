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
	// タイムアウト付きのコンテキストを作成（5分間）
	// 大量データの投入を想定し、少し長めに設定します
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Infrastructure初期化
	db.ConnectDB()
	esClient, err := es.NewClient()
	if err != nil {
		log.Fatalf("❌ Elasticsearchへの接続に失敗しました: %v", err)
	}

	// Repository初期化
	userDBRepo := repository.NewUserRepository(db.DB)
	articleDBRepo := repository.NewArticleRepository(db.DB)
	articleSearchRepo := es.NewArticleSearchRepository(esClient)

	// Usecase初期化
	userUsecase := usecase.NewUserUsecase(userDBRepo)
	articleUsecase := usecase.NewArticleUsecase(articleDBRepo, articleSearchRepo)

	log.Println("🚀 シードデータの投入を開始します...")

	// ユーザーデータのシード投入
	users, err := userUsecase.SeedUsers()
	if err != nil {
		log.Fatalf("❌ ユーザーデータのシード投入中にエラーが発生しました: %v", err)
	}
	log.Printf("✅ %d件のユーザーデータを投入しました。", len(users))

	// 記事データのシード投入（最初のユーザーに紐づけ）
	if len(users) == 0 {
		log.Fatalf("❌ シード投入用のユーザーが存在しません。")
	}
	articles, err := articleUsecase.SeedArticles(users[0].ID)
	if err != nil {
		log.Fatalf("❌ 記事データのシード投入中にエラーが発生しました: %v", err)
	}
	log.Printf("✅ %d件の記事データを投入しました。", len(articles))

	// 記事データを検索エンジンにインデックス
	if err := articleUsecase.ReindexSearchEngine(); err != nil {
		log.Fatalf("❌ 記事データの検索エンジンへのインデックス中にエラーが発生しました: %v", err)
	}

	log.Println("✅ シードデータの投入が正常に完了しました。")
}
