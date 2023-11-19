package domain

import (
	"encoding/json"
	"logswift/internal/db"
	"logswift/internal/dtos"
	"logswift/internal/message_queue"
	"logswift/internal/models"
	"logswift/internal/repository"
	"logswift/pkg/logger"
	"sync/atomic"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

const (
	MAX_BATCH_SIZE = 1000
	FLUSH_INTERVAL = 5 // seconds
	MAX_BUFFER     = 2000
)

type ILogIngestorService interface {
	// IngestLog ingests a log into the system
	IngestLog(logEntry dtos.LogEntry)

	// QueryLogs queries logs from the system
	QueryLogs(searchQuery string, filters dtos.LogQueryFilters) (*dtos.LogQueryResponse, error)
}

type LogIngestorService struct {
	logger        logger.ILogger
	logRepo       []repository.ILogWriterRepository
	flushInterval time.Duration
	buffer        chan dtos.LogEntry
	searchIndex   db.ISearchIndex
	mq            message_queue.IMessageQueue
}

var recordCount int64

func NewLogIngestorService(repo []repository.ILogWriterRepository, searchIndex db.ISearchIndex, mq message_queue.IMessageQueue) ILogIngestorService {
	svc := &LogIngestorService{
		logger:        logger.GetInstance(),
		logRepo:       repo,
		flushInterval: FLUSH_INTERVAL * time.Second,
		buffer:        make(chan dtos.LogEntry, MAX_BATCH_SIZE+MAX_BUFFER),
		searchIndex:   searchIndex,
		mq:            mq,
	}

	go svc.startFlushRoutine()
	return svc
}

func (svc *LogIngestorService) IngestLog(logEntry dtos.LogEntry) {
	select {
	case svc.buffer <- logEntry:
		if len(svc.buffer) >= MAX_BATCH_SIZE {
			svc.flushBuffer()
		}
	default:
		// Buffer is full, consider handling this case based on your requirements
		svc.logger.Warn("Log buffer is full, consider increasing buffer size or handling overflow")
	}

	// non blocking call to publish the log entry to the message queue
	go func() {
		body, err := json.Marshal(logEntry)
		if err != nil {
			svc.logger.Error("Error marshaling log entry", "error", err)
			return
		}

		err = svc.mq.Publish(body)
		if err != nil {
			svc.logger.Error("Error publishing log entry to message queue", "error", err)
			return
		} // // index the record in the search index

	}()
}

func (svc *LogIngestorService) startFlushRoutine() {

	ticker := time.NewTicker(svc.flushInterval)
	defer ticker.Stop()

	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			svc.logger.Info("running time based flushing buffer", "buffer size", len(svc.buffer))
			svc.flushBuffer()
			svc.logger.Info("time based flushing complete", "buffer size", len(svc.buffer))
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (svc *LogIngestorService) flushBuffer() {
	svc.logger.Info("flushing buffer", "buffer size", len(svc.buffer))
	select {
	case logEntry := <-svc.buffer:
		batch := svc.collectBatch(logEntry)

		// Process the batch
		svc.processBatch(batch)
	default:
		// Buffer is empty, nothing to flush
	}
}

func (svc *LogIngestorService) collectBatch(firstLogEntry dtos.LogEntry) []dtos.LogEntry {
	batch := []dtos.LogEntry{firstLogEntry}

	// Try to read more log entries from the buffer without blocking
	for i := 0; i < MAX_BATCH_SIZE; i++ {
		select {
		case logEntry := <-svc.buffer:
			batch = append(batch, logEntry)
		default:
		}
	}

	return batch
}

func (svc *LogIngestorService) processBatch(batch []dtos.LogEntry) {
	if len(batch) == 0 {
		return
	}

	// Perform batch processing and write to the database
	dbIndex := recordCount % int64(len(svc.logRepo))
	// Convert batch to database log entries
	dbLogEntries := make([]models.LogEntry, len(batch))
	for index, logEntry := range batch {

		metadata, err := json.Marshal(logEntry.Metadata)
		if err != nil {
			svc.logger.Error("Error marshaling metadata", "error", err)
			continue
		}

		dbLogEntries[index] = models.LogEntry{
			Level:      logEntry.Level,
			Message:    logEntry.Message,
			ResourceID: logEntry.ResourceID,
			Timestamp:  logEntry.Timestamp,
			TraceID:    logEntry.TraceID,
			SpanID:     logEntry.SpanID,
			Commit:     logEntry.Commit,
			Metadata: postgres.Jsonb{
				RawMessage: metadata,
			},
		}

	}

	// Write the batch to the database using the WriteLogBatch method
	err := svc.logRepo[dbIndex].WriteLogInBatch(dbLogEntries)
	if err != nil {
		svc.logger.Error("Error writing log batch to database", "error", err)
		return
	}

	// increasing record count by 1 so that we can use it for round robin
	atomic.AddInt64(&recordCount, 1)
}

func (svc *LogIngestorService) QueryLogs(searchQuery string, filters dtos.LogQueryFilters) (*dtos.LogQueryResponse, error) {
	return svc.searchIndex.Search(searchQuery, filters)
}
