package adapters

import (
	"github.com/google/uuid"
	"ledger/internal/models"
)

type Repository interface {
	AddMoney(account models.Accounts) error
	GetAccount(userID uuid.UUID) (*models.Accounts, error)
}
