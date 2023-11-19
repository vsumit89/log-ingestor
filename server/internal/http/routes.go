package http

import (
	"logswift/internal/domain"
	"logswift/internal/http/handlers"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterRoutes(
	logSvc domain.ILogIngestorService,
	router *chi.Mux,
) {

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	handlers.RegisterLogHandler(logSvc, router)
}
