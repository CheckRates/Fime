package api

import (
	"github.com/checkrates/Fime/db/postgres"
	"github.com/labstack/echo"
)

// Server handles all HTTP requests for Fime
type Server struct {
	store  *postgres.Store
	router *echo.Echo
}

// NewServer returns a server for Fime
func NewServer(store *postgres.Store) *Server {
	server := &Server{store: store}
	router := echo.New()

	router.POST("/user", server.createUser)

	server.router = router
	return server
}

// Start the Fime echo server
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
