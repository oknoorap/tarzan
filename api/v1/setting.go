package v1

import (
	"net/http"
	"github.com/labstack/echo"
)

// Update setting by given key
func SettingSet (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Get page by given key
func SettingGet (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Show all pages
func SettingList (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}