package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/api/ping", ping)

	go startServer(e)

	waitForGracefulShutdown(e)
}

func ping(c echo.Context) error {
	db := initDB()

	logDate := LogDate{
		Ping:        "Ping",
		CreatedDate: time.Now(),
	}

	db.Table("log_date").Create(&logDate)

	return c.String(http.StatusOK, "OK")
}

func startServer(e *echo.Echo) {
	if err := e.Start(":8080"); err != nil {
		log.Info("shutting down the server")
	}
}

func waitForGracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initDB() *gorm.DB {
	url := os.Getenv("DATABASE_URL")
	dsn := fmt.Sprintf("host=%s port=5432 dbname=postgres user=postgres password=123456 sslmode=disable", url)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db
}

type LogDate struct {
	ID          int       `json:"id" gorm:"primary_key"`
	Ping        string    `json:"ping"`
	CreatedDate time.Time `json:"created_date`
}
