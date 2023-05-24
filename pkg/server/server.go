package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strconv"
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

	return e.Start(":" + strconv.Itoa(port))
}
