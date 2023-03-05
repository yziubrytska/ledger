package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
	"ledger/internal/common"
	"ledger/internal/handlers"
)

type App struct {
	container *dig.Container
}

func NewApp() App {
	return App{container: dig.New()}
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	api := handlers.NewPublicAPI()
	e.POST("/users/:uid/add", api.AddMoney)
	e.POST("/users/:uid/balance", api.Balance)
	e.POST("/users/:uid/history", api.History)

	e.Logger.Fatal(e.Start(":8080")) //:TODO: Add graceful shutdown
}

func build(c *dig.Container) error {
	if err := c.Provide(common.NewConfig); err != nil {
		return err
	}

	return nil

}
