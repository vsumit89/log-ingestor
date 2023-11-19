package repository

import (
	"logswift/internal/db"
	"logswift/internal/db/postgres/logIngestor"
	"logswift/internal/dtos"
	"logswift/internal/models"
)

type ILogWriterRepository interface {
	// WriteLog writes a log to the database
	WriteLog(logEntry models.LogEntry) error

	WriteLogInBatch(logEntries []models.LogEntry) error
}

type ILogReadRepository interface {
	// QueryLogs queries logs from the database
	QueryLogs(searchQuery string, filters dtos.LogQueryFilters) ([]models.LogEntry, error)
}

func NewLogWriterRepository(db db.IDatabase) ILogWriterRepository {
	return logIngestor.NewLogIngestor(db)
}
