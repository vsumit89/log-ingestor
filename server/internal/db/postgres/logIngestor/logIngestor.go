package logIngestor

import (
	"logswift/internal/db"
	"logswift/internal/models"

	"gorm.io/gorm"
)

type PGLogIngestor struct {
	Client *gorm.DB
}

func NewLogIngestor(db db.IDatabase) *PGLogIngestor {
	return &PGLogIngestor{
		Client: db.GetClient().(*gorm.DB),
	}
}

// WriteLog creates a log in the database
func (p *PGLogIngestor) WriteLog(log models.LogEntry) error {
	err := p.Client.Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

// WriteLogInBatch writes logs in batch to the database
func (p *PGLogIngestor) WriteLogInBatch(logEntries []models.LogEntry) error {
	err := p.Client.Create(&logEntries).Error
	if err != nil {
		return err
	}
	return nil
}
