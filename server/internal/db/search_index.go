package db

import (
	"logswift/internal/app/config"
	"logswift/internal/db/meilisearch"
	"logswift/internal/dtos"
)

type ISearchIndex interface {
	Connect(cfg *config.SearchConfig) error
	Disconnect() error
	Search(query string, filter dtos.LogQueryFilters) (*dtos.LogQueryResponse, error)
	Create(log dtos.LogEntry) error
	CreateBatch(logs []dtos.LogEntry) error
}

func NewSearchIndex() ISearchIndex {
	return meilisearch.NewMeiliSearch()
}
