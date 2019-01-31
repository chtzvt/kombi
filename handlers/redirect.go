package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

// Handler
func LinkServer(c echo.Context) error {
	return c.Redirect(http.StatusFound, "Hello, World!")
}

