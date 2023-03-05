package adapters

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ledger/internal/common"
)

type repository struct {
	db *gorm.DB
}

func NewDB(cfg *common.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseCredentials), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
