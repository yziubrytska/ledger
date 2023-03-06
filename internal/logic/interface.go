package logic

import (
	"github.com/google/uuid"
	"ledger/internal/models"
)

type Service interface {
	AddMoney(transaction models.Transactions)
	Balance(userID uuid.UUID) (int64, error)
	History(userID uuid.UUID) ([]*models.Transactions, error)
}
