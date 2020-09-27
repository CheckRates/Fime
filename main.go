package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", handleHome)

	e.Start(":8080")
}

func handleHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}
