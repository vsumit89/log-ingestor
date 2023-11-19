package postgres

import (
	"fmt"

	"logswift/internal/app/config"
	"logswift/internal/models"
	"logswift/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type PGDatabase struct {
	Client *gorm.DB
}

func NewPostgresAdapter() *PGDatabase {
	return &PGDatabase{}
}

func (p *PGDatabase) GetClient() interface{} {
	return p.Client
}

func (p *PGDatabase) Connect(cfg config.DBConfig) error {
	logger := logger.GetInstance()
	logger.Info("connecting to database")
	// dsn (Data source name) is the connection string for the database
	dsn := "host=" + cfg.Host + " user=" + cfg.Username + " password=" + cfg.Password + " dbname= " + cfg.DBName + " port=" + fmt.Sprint(cfg.Port) + " sslmode=" + fmt.Sprint(cfg.SSLMode)
	logger.Info("connecting to database", "dsn", dsn)
	var err error
	p.Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return err
	}
	logger.Info("successfully connected to database")
	return nil
}

func (p *PGDatabase) Migrate() error {
	logger := logger.GetInstance()
	logger.Info("migrating database")
	err := p.Client.AutoMigrate(&models.LogEntry{})
	if err != nil {
		return err
	}
	logger.Info("database migrated successfully")
	return nil
}
