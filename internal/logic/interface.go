package logic

import (
	"github.com/google/uuid"
	"ledger/internal/models"
)

type Service interface {
	AddMoney(userID uuid.UUID, balance string) error
	Balance(userID uuid.UUID) (string, error)
	History(userID uuid.UUID) ([]models.Transaction, error)
}
