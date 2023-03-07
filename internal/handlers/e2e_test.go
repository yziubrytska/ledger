package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/adrianbrad/psqltest"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ledger/internal/adapters"
	"ledger/internal/logic"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	moneyRequest = AddMoneyRequest{
		TransactionID: uuid.MustParse("5e6c281e-f737-4950-8ad5-e6e995ba11e2"),
		Money:         160,
	}

	balanceResponse = BalanceResponse{Balance: 140}
)

func TestPublicAPI(t *testing.T) {
	t.Run("Balance", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		e.Validator = &CustomValidator{Validator: validator.New()}
		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			db := psqltest.NewTransactionTestingDB(t)
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})

			req := httptest.NewRequest(http.MethodPost, "/users/:uid/balance", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:uid/balance")
			c.SetParamNames("uid")
			c.SetParamValues("24f594ab-2207-4b45-897d-7ad948b60896")

			h := PublicAPI{logic.NewService(adapters.NewRepository(gormDB), nil), nil}

			expected, _ := json.Marshal(balanceResponse)
			expected = append(expected, '\n')
			// Assertions
			if assert.NoError(t, h.Balance(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, string(expected), rec.Body.String())
			}
		})
	})

	t.Run("AddMoney", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		e.Validator = &CustomValidator{Validator: validator.New()}
		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			db := psqltest.NewTransactionTestingDB(t)
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})

			moneyBody, _ := json.Marshal(moneyRequest)
			req := httptest.NewRequest(http.MethodPost, "/users/:uid/add", bytes.NewBuffer(moneyBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:uid/add")
			c.SetParamNames("uid")
			c.SetParamValues("24f594ab-2207-4b45-897d-7ad948b60896")
			h := PublicAPI{logic.NewService(adapters.NewRepository(gormDB), nil), nil}

			// Assertions
			if assert.NoError(t, h.AddMoney(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
			time.Sleep(1 * time.Second)
			if assert.NoError(t, h.History(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	})
}
