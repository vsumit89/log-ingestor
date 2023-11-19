package app

import (
	"context"
	"fmt"
	"logswift/internal/app/config"
	"logswift/internal/db"
	"logswift/internal/repository"

	"logswift/internal/domain"
	httpx "logswift/internal/http"
	"logswift/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var AppCfg *config.AppConfig

type App struct {
	Router *chi.Mux
	Server *http.Server
	DBs    []db.IDatabase
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() error {
	log := logger.GetInstance()

	// logWriteRepo is a slice of ILogWriterRepository
	// this is because we can have multiple databases to write to
	// we write to multiple dbs to ensure inserts are fast
	var logWriteRepo []repository.ILogWriterRepository

	// setting up and connecting to the database as well as setting up the migrations and repositories
	for index, dbCfg := range AppCfg.DB {
		dbService := db.NewDBService()
		log.Info("connecting to the database", "db", dbCfg, "index", index)
		err := dbService.Connect(dbCfg)
		if err != nil {
			return err
		}
		err = dbService.Migrate()
		if err != nil {
			return err
		}

		a.DBs = append(a.DBs, dbService)

		logWriteRepo = append(logWriteRepo, repository.NewLogWriterRepository(dbService))
	}

	// searchIndexSvc is the service that will be used to search logs in the database
	// it will help in full text search
	log.Info("connecting to the search index", "search", AppCfg.Search)
	searchIndexSvc := db.NewSearchIndex()
	err := searchIndexSvc.Connect(AppCfg.Search)
	if err != nil {
		return err
	}

	// logIngestorSvc is the service that will be used to write logs to the database
	// we pass in the logWriteRepo to the service as it is a dependency
	logIngestorSvc := domain.NewLogIngestorService(logWriteRepo, searchIndexSvc)

	// setting up the http router
	// using go-chi as the router
	a.Router = chi.NewRouter()

	// setting up cors
	a.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// setting up NilPointerMiddleware to handle nil pointer errors
	// a.Router.Use(httpx.NilPointerMiddleware)
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(log.GetHTTPMiddleWare())

	httpx.RegisterRoutes(logIngestorSvc, a.Router)

	log.Info("starting http server", "port", AppCfg.Server.Port)
	a.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", AppCfg.Server.Port),
		Handler: a.Router,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error starting the server"+err.Error(), "error", err.Error())
		}
	}()

	// Block until we receive our signal
	<-shutdown

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	a.Server.Shutdown(ctx)
	log.Info("shutting down the application")
	return nil
}
