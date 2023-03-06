package migrations

import (
	"github.com/pkg/errors"
	"ledger/internal/adapters"
	"ledger/internal/common"
	"ledger/internal/models"
)

func Migrate(config *common.Config) error {
	db, sqlDB, err := adapters.NewDB(config)
	if err != nil {
		return err
	}

	db.AutoMigrate(&models.Accounts{}, &models.PgTransactions{})
	return sqlDB.Close()
}

func Drop(config *common.Config) error {
	db, sqlDB, err := adapters.NewDB(config)
	if err != nil {
		return err
	}

	err = db.Migrator().DropTable(&models.Accounts{}, &models.PgTransactions{})
	if err != nil {
		return errors.Wrap(err, "error while deleting user table")
	}
	return sqlDB.Close()
}
