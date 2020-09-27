package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	echo := echo.New()

	// Set up static pages
	echo.Static("/", "./views")

	// Set up API group of handles
	api := echo.Group("/api/v1")
	{
		api.GET("/", handleHome)
		api.GET("/images", handleImages)
	}

	// Start server
	echo.Start(":8080")
}

func handleHome(c echo.Context) error {
	return c.String(http.StatusOK, "69 Lmao")
}

func handleImages(c echo.Context) error {
	return c.String(http.StatusOK, "This should show all images ")
}
