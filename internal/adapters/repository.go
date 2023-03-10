package adapters

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ledger/internal/common"
	"ledger/internal/models"
)

type repository struct {
	db *gorm.DB
}

func NewDB(cfg *common.Config) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseCredentials), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error while opening a connection to the db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error while getting an sql instance")
	}

	return db, sqlDB, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r repository) AddMoney(account models.Accounts) error {
	k := models.Transactions{}
	r.db.Where("id = ?", account.Transactions[0].ID).FirstOrCreate(&k)
	if err := r.db.Model(&models.Transactions{}).Create(&account.Transactions).Error; err != nil {
		return err
	}
	err := r.db.Model(&models.Accounts{ID: account.ID}).Save(&account).Error
	if err != nil {
		return errors.Wrap(err, "error while updating a balance")
	}
	return nil
}

func (r repository) GetAccount(userID uuid.UUID) (*models.Accounts, error) {
	var result models.Accounts
	err := r.db.Model(models.Accounts{ID: userID}).Preload("Transactions", func(db *gorm.DB) *gorm.DB {
		return db.Order("transactions.date ASC")
	}).First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
