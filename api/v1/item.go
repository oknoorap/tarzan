package v1

import (
	"net/http"
	"github.com/labstack/echo"
)

// Get item by given id
func ItemGet (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

// Show all items
func ItemList (c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
