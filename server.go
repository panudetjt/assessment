package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/panudetjt/assessment/expense"
	"github.com/panudetjt/assessment/health"
)

func main() {
	port := os.Getenv("PORT")

	db := expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health.HealthHandler)

	eh := &expense.Handler{DB: db}

	e.POST("/expenses", eh.CreateExpensesHandler)
	e.GET("/expenses/:id", eh.GetExpenseByIdHandler)
	e.PUT("/expenses/:id", eh.UpdateExpensesHandler)

	go func() {
		e.Logger.Info("Server started at ", port)
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Logger.Info("closing the database connection")
	if err := db.Close(); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("shutting down the server")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("bye bye!")
}
