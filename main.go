package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"ledger/internal/adapters"
	"ledger/internal/common"
	"ledger/internal/handlers"
	"ledger/internal/logic"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	c := dig.New()
	if err := c.Provide(common.NewConfig); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(adapters.NewDB); err != nil {
		log.Fatal(err)
	}

	if err := c.Provide(adapters.NewRepository); err != nil {
		log.Fatal(err)
	}

	if err := c.Provide(logic.NewService); err != nil {
		log.Fatal(err)
	}

	if err := c.Provide(handlers.NewPublicAPI); err != nil {
		log.Fatal(err)
	}

	if err := c.Provide(func(config *common.Config) (logrus.Level, error) {
		return logrus.ParseLevel(config.LogLevel)
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.Provide(handlers.NewLogrusEntry); err != nil {
		log.Fatal(err)
	}

	err := c.Invoke(func(api *handlers.PublicAPI) {
		e.POST("/users/:uid/add", api.AddMoney)
		e.POST("/users/:uid/balance", api.Balance)
		e.POST("/users/:uid/history", api.History)
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
