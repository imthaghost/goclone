package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Serve ...
func Serve(projectPath string, port int) error {
	e := echo.New()

	// Log Output
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// static files
	e.Static("/", projectPath)
	e.File("/", projectPath)

	// Create a context that can be cancelled
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for the interrupt signal
	<-ctx.Done()

	// Graceful shutdown of the server
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
