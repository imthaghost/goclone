package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Serve ...
func Serve(projectPath string) error {
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

	return e.Start(":5000")
}
