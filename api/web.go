package main

import (
	"net/http"
	"github.com/labstack/echo"
)

func Dashboard (c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", "")
}