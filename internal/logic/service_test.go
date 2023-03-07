package logic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"ledger/internal/adapters"
	"ledger/internal/adapters/mocks"
	"ledger/internal/models"
	"reflect"
	"testing"
	"time"
)

func Test_service_AddMoney(t *testing.T) {
	userID := uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896")
	logger, _ := NewLogrusEntry(logrus.DebugLevel)
	transaction := models.Transactions{
		ID:     uuid.MustParse("6d21b813-7840-48af-92ca-84513c05c62d"),
		UserID: userID,
		Date:   time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC),
		Sum:    200,
	}
	account := &models.Accounts{
		ID:      uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896"),
		Email:   "test@gmail.com",
		Balance: 500,
	}
	fullAccount := models.Accounts{
		ID:      uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896"),
		Email:   "test@gmail.com",
		Balance: 700,
		Transactions: []*models.Transactions{
			{
				ID:     uuid.MustParse("6d21b813-7840-48af-92ca-84513c05c62d"),
				UserID: userID,
				Date:   time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC),
				Sum:    200,
			},
		},
	}

	type fields struct {
		repo   adapters.Repository
		logger *logrus.Entry
	}
	type args struct {
		transaction models.Transactions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "GetAccount record not found error",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(account, gorm.ErrRecordNotFound)
					return repo
				}(),
				logger: logger,
			},
			args: args{transaction: transaction},
		},
		{
			name: "GetAccount internal server error",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(nil, errors.New("internal server error"))
					return repo
				}(),
				logger: logger,
			},
			args: args{transaction: transaction},
		},
		{
			name: "Add money transaction already executed",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(account, nil)
					repo.On("AddMoney", fullAccount).Return(&pgconn.PgError{Code: "23505"})
					return repo
				}(),
				logger: logger,
			},
			args: args{transaction: transaction},
		},
		{
			name: "Add money transaction service error",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(account, nil)
					repo.On("AddMoney", fullAccount).Return(errors.New("internal service error"))
					return repo
				}(),
				logger: logger,
			},
			args: args{transaction: transaction},
		},
		{
			name: "ok",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(account, nil)
					repo.On("AddMoney", fullAccount).Return(nil)
					return repo
				}(),
				logger: logger,
			},
			args: args{transaction: transaction},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
				timeNow: func() time.Time {
					return time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
				},
			}
			s.AddMoney(tt.args.transaction)
		})
	}
}

func Test_service_Balance(t *testing.T) {
	userID := uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896")
	logger, _ := NewLogrusEntry(logrus.DebugLevel)
	account := &models.Accounts{
		ID:      uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896"),
		Email:   "test@gmail.com",
		Balance: 500,
	}
	type fields struct {
		repo    adapters.Repository
		logger  *logrus.Entry
		timeNow func() time.Time
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "get account error",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(nil, errors.New("internal service error"))
					return repo
				}(),
				logger: logger,
				timeNow: func() time.Time {
					return time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
				},
			},
			args:    args{userID: userID},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(account, nil)
					return repo
				}(),
				logger: logger,
				timeNow: func() time.Time {
					return time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
				},
			},
			args: args{userID: userID},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repo:    tt.fields.repo,
				logger:  tt.fields.logger,
				timeNow: tt.fields.timeNow,
			}
			got, err := s.Balance(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Balance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_History(t *testing.T) {
	userID := uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896")
	logger, _ := NewLogrusEntry(logrus.DebugLevel)
	fullAccount := &models.Accounts{
		ID:      uuid.MustParse("24f594ab-2207-4b45-897d-7ad948b60896"),
		Email:   "test@gmail.com",
		Balance: 700,
		Transactions: []*models.Transactions{
			{
				ID:     uuid.MustParse("6d21b813-7840-48af-92ca-84513c05c62d"),
				UserID: userID,
				Date:   time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC),
				Sum:    200,
			},
		},
	}

	type fields struct {
		repo    adapters.Repository
		logger  *logrus.Entry
		timeNow func() time.Time
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.Transactions
		wantErr bool
	}{
		{
			name: "get account error",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(nil, errors.New("internal service error"))
					return repo
				}(),
				logger: logger,
				timeNow: func() time.Time {
					return time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
				},
			},
			args:    args{userID: userID},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				repo: func() adapters.Repository {
					repo := new(mocks.Repository)
					repo.On("GetAccount", userID).Return(fullAccount, nil)
					return repo
				}(),
				logger: logger,
				timeNow: func() time.Time {
					return time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
				},
			},
			args: args{userID: userID},
			want: fullAccount.Transactions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repo:    tt.fields.repo,
				logger:  tt.fields.logger,
				timeNow: tt.fields.timeNow,
			}
			got, err := s.History(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("History() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("History() got = %v, want %v", got, tt.want)
			}
		})
	}
}
