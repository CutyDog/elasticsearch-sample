package graph

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
	ES *elasticsearch.TypedClient
}
