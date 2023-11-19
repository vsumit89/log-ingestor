package handlers

import (
	"logswift/internal/domain"

	"github.com/go-chi/chi"
)

// logHandler holds the dependencies for the log ingestor handler
type logHandler struct {
	logSvc domain.ILogIngestorService
	router *chi.Mux
}

// RegisterLogHandler will register the log ingestor routes
// with the handlers
func RegisterLogHandler(
	logSvc domain.ILogIngestorService,
	router *chi.Mux,
) {
	handler := &logHandler{
		logSvc: logSvc,
		router: router,
	}

	router.Post("/", handler.ingest)
	router.Get("/query", handler.query)
}
