package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/checkrates/Fime/token"
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

// Server handles all HTTP requests and manages Database calls for Fime
type Server struct {
	store  postgres.Store
	config config.Config
	router *echo.Echo
	aws    *session.Session
	token  token.Maker
}

// NewServer returns a server for Fime
func NewServer(config config.Config, store postgres.Store) (*Server, error) {
	token, err := token.NewJWTMaker(config.Token.AccessSecret)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:  store,
		token:  token,
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := echo.New()
	router.Validator = &Validator{val: validator.New()}

	router.POST("/user/login", server.loginUser)
	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUser)
	router.GET("/user", server.listUsers)

	router.POST("/image", server.postImage)
	router.GET("/image/:id", server.getImage)
	router.DELETE("/image/:id", server.deleteImage)
	router.PATCH("/image", server.updateImage)
	router.GET("/image", server.listImage)
	router.GET("/image/user/:id", server.listUserImages)

	router.GET("/tag", server.listTags)
	router.GET("/tag/:id", server.listUserTags)

	server.router = router
}

// Start the Fime echo server
func (server *Server) Start(address string) error {
	var err error
	if server.aws, err = connectAWS(server.config.S3); err != nil {
		return err
	}
	return server.router.Start(address)
}

// errorResponse formats error to Echo
func errorResponse(err error) echo.Map {
	return echo.Map{"error": err.Error()}
}
