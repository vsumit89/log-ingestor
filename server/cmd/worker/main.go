package main

import (
	"encoding/json"
	"logswift/internal/app"
	"logswift/internal/db"
	"logswift/internal/dtos"
	"logswift/internal/message_queue"
	"logswift/pkg/logger"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	MAX_BATCH_SIZE = 1000
	FLUSH_INTERVAL = 5 // seconds
	MAX_BUFFER     = 2000
)

type messageHandler struct {
	searchIndex   db.ISearchIndex
	flushInterval time.Duration
	buffer        chan dtos.LogEntry
	logger        logger.ILogger
}

func NewMessageHandler(searchIndex db.ISearchIndex) *messageHandler {
	mhHandler := &messageHandler{
		flushInterval: FLUSH_INTERVAL * time.Second,
		buffer:        make(chan dtos.LogEntry, MAX_BATCH_SIZE+MAX_BUFFER),
		searchIndex:   searchIndex,
		logger:        logger.GetInstance(),
	}

	go mhHandler.startFlushRoutine()
	return mhHandler
}

func (mh *messageHandler) HandleMessage(message []byte) {
	mh.logger.Info("message received", "message", string(message))

	var logEntry dtos.LogEntry
	err := json.Unmarshal(message, &logEntry)
	if err != nil {
		mh.logger.Error("error unmarshalling message", "error", err.Error())
		return
	}

	select {
	case mh.buffer <- logEntry:
		if len(mh.buffer) >= MAX_BATCH_SIZE {
			mh.flushBuffer()
		}
	default:
		// Buffer is full, consider handling this case based on your requirements
		mh.logger.Warn("Log buffer is full, consider increasing buffer size or handling overflow")
	}
}

func main() {
	log := logger.GetInstance()

	log.Info("starting the application")

	log.Info("reading config file")

	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Error("error reading config file", "error", err.Error())
		return
	}

	err = yaml.Unmarshal(file, &app.AppCfg)
	if err != nil {
		log.Error("error unmarshalling config file", "error", err.Error())
		return
	}

	log.Info("config file read successfully")

	mq := message_queue.NewMessageQueue()

	err = mq.Connect(*app.AppCfg.MQCfg)
	if err != nil {
		log.Error("error connecting to message queue", "error", err.Error())
		return
	}

	err = mq.DeclareQueue()
	if err != nil {
		log.Error("error declaring queue", "error", err.Error())
		return
	}

	searchIndex := db.NewSearchIndex()

	err = searchIndex.Connect(app.AppCfg.Search)
	if err != nil {
		log.Error("error connecting to search index", "error", err.Error())
		return
	}

	mh := NewMessageHandler(searchIndex)
	mq.Consume(mh)
}

func (mh *messageHandler) startFlushRoutine() {
	ticker := time.NewTicker(mh.flushInterval)
	defer ticker.Stop()

	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			mh.logger.Info("running time based flushing buffer", "buffer size", len(mh.buffer))
			mh.flushBuffer()
			mh.logger.Info("time based flushing complete", "buffer size", len(mh.buffer))
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (mh *messageHandler) flushBuffer() {
	mh.logger.Info("flushing buffer", "buffer size", len(mh.buffer))
	select {
	case logEntry := <-mh.buffer:
		batch := mh.collectBatch(logEntry)

		// Process the batch
		mh.processBatch(batch)
	default:
		// Buffer is empty, nothing to flush
	}
}

func (mh *messageHandler) collectBatch(firstLogEntry dtos.LogEntry) []dtos.LogEntry {
	batch := []dtos.LogEntry{firstLogEntry}

	// Try to read more log entries from the buffer without blocking
	for i := 0; i < MAX_BATCH_SIZE; i++ {
		select {
		case logEntry := <-mh.buffer:
			batch = append(batch, logEntry)
		default:
		}
	}

	return batch
}

func (mh *messageHandler) processBatch(batch []dtos.LogEntry) {
	mh.logger.Info("processing batch", "batch size", len(batch))
	if len(batch) == 0 {
		return
	}

	retry := 3
	// Perform batch processing and write to the database
	for i := 0; i < retry; i++ {
		err := mh.searchIndex.CreateBatch(batch)
		if err != nil {
			mh.logger.Error("error creating batch", "error", err.Error())
			mh.logger.Info("retrying", "retry", i+1)
			continue
		}
		break
	}

	mh.logger.Info("batch processing complete", "batch size", len(batch))
}
