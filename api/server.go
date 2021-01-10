package api

import (
	"github.com/checkrates/Fime/db/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

// Single instance request validator
type Validator struct {
	val *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.val.Struct(i)
}

// Server handles all HTTP requests for Fime
type Server struct {
	store  *postgres.Store
	router *echo.Echo
}

// NewServer returns a server for Fime
func NewServer(store *postgres.Store) *Server {
	server := &Server{store: store}
	router := echo.New()
	router.Validator = &Validator{val: validator.New()}

	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUser)
	router.GET("users", server.listUsers)

	server.router = router
	return server
}

// Start the Fime echo server
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

// errorResponse formats error to Echo
func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
