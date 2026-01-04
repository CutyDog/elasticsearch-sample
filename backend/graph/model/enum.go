package model

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "DRAFT"
	ArticleStatusPublished ArticleStatus = "PUBLISHED"
	ArticleStatusArchived  ArticleStatus = "ARCHIVED"
)
