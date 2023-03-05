package handlers

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"ledger/internal/logic"
	"net/http"
)

type AddMoneyRequest struct {
	Money string `json:"money" validate:"required"`
}

type BalanceResponse struct {
	Balance string `json:"balance" validate:"required"`
}

type HistoryResponse struct {
	History []struct {
		Date        string `json:"date"`
		Transaction string `json:"transaction"`
	} `json:"history"`
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

type publicAPI struct {
	service logic.Service //TODO: Add logger
}

func NewPublicAPI() LedgerApi {
	return publicAPI{}
}

func (p publicAPI) AddMoney(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m := new(AddMoneyRequest)
	if err = c.Bind(m); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(m); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = p.service.AddMoney(userID, m.Money)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusOK)
}

func (p publicAPI) Balance(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	balance, err := p.service.Balance(userID)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, BalanceResponse{Balance: balance})
}

func (p publicAPI) History(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

}
