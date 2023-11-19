package db

import (
	"logswift/internal/app/config"
	"logswift/internal/db/postgres"
)

type DBOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  bool
}
type IDatabase interface {
	Connect(cfg config.DBConfig) error
	GetClient() interface{}
	Migrate() error
}

func NewDBService() IDatabase {
	return postgres.NewPostgresAdapter()
}
