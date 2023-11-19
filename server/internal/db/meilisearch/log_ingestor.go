package meilisearch

import (
	"fmt"
	"logswift/internal/dtos"
	"logswift/pkg/helper"
	"math/rand"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func getRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func (ms *MeiliSearch) Create(log dtos.LogEntry) error {
	formattedLog := map[string]interface{}{
		"id":                 getRandomID(),
		"timestamp":          log.Timestamp.Unix(),
		"message":            log.Message,
		"level":              log.Level,
		"resourceId":         log.ResourceID,
		"traceId":            log.TraceID,
		"spanId":             log.SpanID,
		"commit":             log.Commit,
		"metadata_parent_id": log.Metadata.ParentResourceID,
	}

	_, err := ms.client.Index(ms.indexName).AddDocuments([]interface{}{formattedLog}, "id")
	if err != nil {
		return err
	}

	return err
}

func (ms *MeiliSearch) CreateBatch(logs []dtos.LogEntry) error {
	formattedLogs := []interface{}{}
	for _, log := range logs {
		formattedLog := map[string]interface{}{
			"id":                 getRandomID(),
			"timestamp":          log.Timestamp.Unix(),
			"message":            log.Message,
			"level":              log.Level,
			"resourceId":         log.ResourceID,
			"traceId":            log.TraceID,
			"spanId":             log.SpanID,
			"commit":             log.Commit,
			"metadata_parent_id": log.Metadata.ParentResourceID,
		}
		formattedLogs = append(formattedLogs, formattedLog)
	}

	_, err := ms.client.Index(ms.indexName).AddDocuments(formattedLogs, "id")
	if err != nil {
		return err
	}

	return err
}

func (ms *MeiliSearch) Search(query string, filter dtos.LogQueryFilters) (*dtos.LogQueryResponse, error) {
	searchParams := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(filter.Limit),
		Offset: int64((filter.Page - 1) * filter.Limit),
	}

	filters := [][]string{}
	if filter.Level != nil {
		level := fmt.Sprintf("level=%s", *filter.Level)
		filters = append(filters, []string{level})
	}

	if filter.Message != nil {
		message := fmt.Sprintf("message='%s'", *filter.Message)
		filters = append(filters, []string{message})
	}

	if filter.ResourceID != nil {
		resourceID := fmt.Sprintf("resourceId=%s", *filter.ResourceID)
		filters = append(filters, []string{resourceID})
	}

	if filter.TraceID != nil {
		traceID := fmt.Sprintf("traceId=%s", *filter.TraceID)
		filters = append(filters, []string{traceID})
	}

	if filter.SpanID != nil {
		spanID := fmt.Sprintf("spanId=%s", *filter.SpanID)
		filters = append(filters, []string{spanID})
	}

	if filter.Commit != nil {
		commit := fmt.Sprintf("commit=%s", *filter.Commit)
		filters = append(filters, []string{commit})
	}

	if filter.ParentResourceID != nil {
		parentResourceID := fmt.Sprintf("metadata_parent_id=%s", *filter.ParentResourceID)
		filters = append(filters, []string{parentResourceID})
	}

	dateFilter := helper.GetFiltersUsingTime(filter.From, filter.To)
	ms.logger.Info("Date filter", "date filter", dateFilter)
	if len(dateFilter) > 0 {
		filters = append(filters, []string{dateFilter})
	}

	searchParams.Filter = filters

	searchResult, err := ms.client.Index(ms.indexName).Search(query, searchParams)
	if err != nil {
		return nil, err
	}

	var result []dtos.LogEntry
	for _, doc := range searchResult.Hits {
		eachDoc := doc.(map[string]interface{})
		data := dtos.LogEntry{
			Level:      eachDoc["level"].(string),
			Message:    eachDoc["message"].(string),
			ResourceID: eachDoc["resourceId"].(string),
			TraceID:    eachDoc["traceId"].(string),
			SpanID:     eachDoc["spanId"].(string),
			Commit:     eachDoc["commit"].(string),
			Metadata: dtos.Metadata{
				ParentResourceID: eachDoc["metadata_parent_id"].(string),
			},
		}

		timestamp := eachDoc["timestamp"].(float64)
		// if ok {
		data.Timestamp = time.Unix(int64(timestamp), 0)
		// } else {
		// 	log.Println("Failed to parse timestamp")
		// }

		result = append(result, data)
	}
	return &dtos.LogQueryResponse{
		Logs:  result,
		Total: int(searchResult.EstimatedTotalHits),
	}, nil
}
