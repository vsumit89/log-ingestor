package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type LogEntry struct {
	ID         int64          `json:"id"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	ResourceID string         `json:"resourceId"`
	Timestamp  time.Time      `json:"timestamp"`
	TraceID    string         `json:"traceId"`
	SpanID     string         `json:"spanId"`
	Commit     string         `json:"commit"`
	Metadata   postgres.Jsonb `json:"metadata"`
}
