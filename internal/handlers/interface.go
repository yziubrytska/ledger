package handlers

import "github.com/labstack/echo/v4"

type LedgerApi interface {
	AddMoney(c echo.Context) error
	Balance(c echo.Context) error
	History(c echo.Context) error
}
