package handlers

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger/internal/logic"
	"ledger/internal/models"
	"net/http"
)

type AddMoneyRequest struct {
	TransactionID uuid.UUID `json:"transactionID" validate:"required"`
	Money         int64     `json:"money" validate:"required"`
}

type BalanceResponse struct {
	Balance int64 `json:"balance" validate:"required"`
}

type HistoryResponse struct {
	History []*models.Transactions `json:"history"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type PublicAPI struct {
	service logic.Service
	logger  *logrus.Entry
}

func NewPublicAPI(s logic.Service, l *logrus.Entry) *PublicAPI {
	return &PublicAPI{
		service: s,
		logger:  l,
	}
}

func (p PublicAPI) AddMoney(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		p.logger.WithError(err).Error("error while parsing the user id param")

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m := new(AddMoneyRequest)
	if err = c.Bind(m); err != nil {
		p.logger.WithError(err).Error("error binding the request")

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(m); err != nil {
		p.logger.WithError(err).Error("error validating the struct")

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	go p.service.AddMoney(models.Transactions{
		ID:     m.TransactionID,
		UserID: userID,
		Sum:    m.Money,
	})

	return echo.NewHTTPError(http.StatusOK)
}

func (p PublicAPI) Balance(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		p.logger.WithError(err).Error("error parsing the user id param")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	balance, err := p.service.Balance(userID)
	if err != nil {
		p.logger.WithError(err).Error("error getting the balance by user id %s", userID.String())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, BalanceResponse{Balance: balance})
}

func (p PublicAPI) History(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		p.logger.WithError(err).Error("error parsing the user id param")

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transactions, err := p.service.History(userID)
	if err != nil {
		p.logger.WithError(err).Errorf("error getting the history by user id %s", userID.String())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, HistoryResponse{History: transactions})
}
