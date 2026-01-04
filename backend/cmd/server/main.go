package main

import (
	"elasticsearch-sample/backend/graph"
	"elasticsearch-sample/backend/internal/domain/repository"
	"elasticsearch-sample/backend/internal/infrastructure/db"
	"elasticsearch-sample/backend/internal/infrastructure/es"
	"elasticsearch-sample/backend/internal/usecase"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	// Infrastructure初期化
	db.ConnectDB()
	esClient, err := es.NewClient()
	if err != nil {
		log.Fatalf("Elasticsearch接続失敗 (設定やプラグインを確認してください): %v", err)
	}

	// Repository初期化
	articleDBRepo := repository.NewArticleRepository(db.DB)
	articleSearchRepo := es.NewArticleSearchRepository(esClient)

	// Usecase初期化
	articleUsecase := usecase.NewArticleUsecase(articleDBRepo, articleSearchRepo)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		ArticleUsecase: &articleUsecase,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", c.Handler(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
