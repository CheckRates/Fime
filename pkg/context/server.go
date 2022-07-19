package http

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/token"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	router.Use(middleware.RemoveTrailingSlash())
	router.Validator = &Validator{val: validator.New()}

	router.POST("/user/login", server.loginUser)
	router.POST("/user", server.createUser)

	authRoutes := router.Group("/*", authMiddleware(server.token))

	authRoutes.GET("/user/:id", server.getUser)
	authRoutes.GET("/user", server.listUsers)

	authRoutes.POST("/image", server.postImage)
	authRoutes.GET("/image/:id", server.getImage)
	authRoutes.DELETE("/image/:id", server.deleteImage)
	authRoutes.PATCH("/image", server.updateImage)
	authRoutes.GET("/image", server.listImage)
	authRoutes.GET("/image/user/:id", server.listUserImages)

	authRoutes.GET("/tag", server.listTags)
	authRoutes.GET("/tag/:id", server.listUserTags)

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
