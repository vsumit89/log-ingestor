package meilisearch

import (
	"fmt"
	"logswift/internal/app/config"
	"logswift/pkg/logger"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearch struct {
	logger    logger.ILogger
	client    *meilisearch.Client
	indexName string
}

func NewMeiliSearch() *MeiliSearch {
	return &MeiliSearch{
		logger: logger.GetInstance(),
	}
}

func (ms *MeiliSearch) Connect(cfg *config.SearchConfig) error {
	ms.logger.Info("Connecting to MeiliSearch", "host", cfg.Host, "port", cfg.Port, "name", cfg.Name)
	ms.client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		APIKey: cfg.API_KEY,
	})

	filterableAttr := []string{"level", "resourceId", "traceId", "spanId", "commit", "metadata_parent_id", "timestamp", "message"}
	_, err := ms.client.Index(cfg.Name).UpdateFilterableAttributes(&filterableAttr)
	if err != nil {
		return err
	}

	ms.indexName = cfg.Name

	ms.logger.Info("Connected to MeiliSearch", "host", cfg.Host, "port", cfg.Port, "name", cfg.Name)
	return nil
}

func (ms *MeiliSearch) Disconnect() error {
	return nil
}
