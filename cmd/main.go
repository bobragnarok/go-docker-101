package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/api/ping", ping)

	go startServer(e)

	waitForGracefulShutdown(e)
}

func ping(c echo.Context) error {
	var ping string

	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println("File reading error", err)
	}
	defer file.Close()

	data := make([]byte, 1024)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("File reading error", err)
		}
		ping = string(data[:n])
	}

	if len(ping) == 0 {
		return c.String(http.StatusOK, "OK")
	} else {
		return c.String(http.StatusOK, ping)
	}
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
