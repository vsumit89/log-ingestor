package dtos

import "time"

// LogEntry is used as a DTO for the log entry
// it is used for ingesting data into the system
type LogEntry struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ResourceID string    `json:"resourceId"`
	Timestamp  time.Time `json:"timestamp"`
	TraceID    string    `json:"traceId"`
	SpanID     string    `json:"spanId"`
	Commit     string    `json:"commit"`
	Metadata   Metadata  `json:"metadata"`
}

// Metadata is used to store metadata about the log entry
// FUTURE: we can make it a map[string]interface{} to allow
// for more flexibility
type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}

// LogQueryResponse is used as a DTO for the log query response
type LogQueryResponse struct {
	Logs  []LogEntry `json:"logs"`
	Total int        `json:"total"`
}

// LogQueryFilters is used as a DTO for the log query filters
// it is used for querying data from the storage
type LogQueryFilters struct {
	Level            *string
	Message          *string
	ResourceID       *string
	From             *time.Time
	To               *time.Time
	TraceID          *string
	SpanID           *string
	Commit           *string
	ParentResourceID *string
	Limit            int
	Page             int
}
