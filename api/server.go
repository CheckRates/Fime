package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

// Validator struct for custom echo server validation
type Validator struct {
	val *validator.Validate
}

// Validate the request body struct
func (v *Validator) Validate(i interface{}) error {
	return v.val.Struct(i)
}

// Server handles all HTTP requests for Fime
type Server struct {
	store  *postgres.Store
	router *echo.Echo
	aws    *session.Session
}

// NewServer returns a server for Fime
func NewServer(store *postgres.Store) *Server {
	server := &Server{store: store}
	router := echo.New()
	router.Validator = &Validator{val: validator.New()}

	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUser)
	router.GET("users", server.listUsers)

	router.POST("/image", server.postImage)
	router.GET("/image/:id", server.getImage)
	router.DELETE("/image/:id", server.deleteImage)
	router.PATCH("image", server.updateImage)
	router.GET("images", server.listImage)

	router.GET("/tags", server.listTags)

	server.router = router
	return server
}

// Start the Fime echo server
func (server *Server) Start(address string) error {
	var err error
	if server.aws, err = connectAWS(); err != nil {
		return err
	}
	return server.router.Start(address)
}

// errorResponse formats error to Echo
func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
