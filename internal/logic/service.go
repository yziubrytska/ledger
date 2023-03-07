package logic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"ledger/internal/adapters"
	"ledger/internal/models"
	"time"
)

var (
	ErrAlreadyExecuted = errors.New("transaction was already executed")
)

type service struct {
	repo    adapters.Repository
	logger  *logrus.Entry
	timeNow func() time.Time
}

func NewService(r adapters.Repository, l *logrus.Entry) Service {
	return &service{
		repo:    r,
		logger:  l,
		timeNow: time.Now,
	}
}

func (s service) AddMoney(transaction models.Transactions) {
	account, err := s.repo.GetAccount(transaction.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.WithError(err).Errorf("error record not found for user id %s", transaction.UserID.String())

		return
	}
	if err != nil {
		s.logger.WithError(err).Errorf("error getting account for user id %s", transaction.UserID.String())

		return
	}

	err = s.repo.AddMoney(models.Accounts{
		ID:      account.ID,
		Email:   account.Email,
		Balance: account.Balance + transaction.Sum,
		Transactions: []*models.Transactions{
			{
				ID:     transaction.ID,
				UserID: transaction.UserID,
				Sum:    transaction.Sum,
				Date:   s.timeNow(),
			},
		},
	})
	if dbError, ok := err.(*pgconn.PgError); ok && dbError.Code == "23505" {
		s.logger.WithError(err).Errorf("transaction already executed for user id %s", transaction.UserID.String())

		return
	}
	if err != nil {
		s.logger.WithError(err).Errorf("error while adding money for user id %s", transaction.UserID.String())
	}
}

func (s service) Balance(userID uuid.UUID) (int64, error) {
	account, err := s.repo.GetAccount(userID)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (s service) History(userID uuid.UUID) ([]*models.Transactions, error) {
	account, err := s.repo.GetAccount(userID)
	if err != nil {
		return nil, err
	}

	return account.Transactions, nil
}
