package graph

import (
	"elasticsearch-sample/backend/internal/usecase"
)

type Resolver struct {
	ArticleUsecase *usecase.ArticleUsecase
}
