package adapters

import (
	"github.com/google/uuid"
	"ledger/internal/models"
)

type Repository interface {
	//GetTransactionByID(id uuid.UUID) (models.Transactions, error)
	AddMoney(account models.Accounts) error
	GetAccount(userID uuid.UUID) (*models.Accounts, error)
}
