package rest

import (
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/pkg/context"
	"github.com/checkrates/Fime/pkg/http"
	"github.com/checkrates/Fime/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Interface for starting and stopping a RESTful API
type Server interface {
	Close() error
	Start(address string) error
}

// Validator struct for custom echo server validation
type Validator struct {
	val *validator.Validate
}

// Validate the request body struct
func (v *Validator) Validate(i interface{}) error {
	return v.val.Struct(i)
}

// Server handles all HTTP requests and manages Database calls for Fime
type restServer struct {
	config   config.Config
	router   *echo.Echo
	database *sqlx.DB

	ports struct {
		user http.UserPort
		tag  http.TagPort
		post http.PostPort
	}
}

// NewServer returns a server for Fime
func NewServer(config config.Config) (Server, error) {
	e := echo.New()

	conn, err := sqlx.Open("postgres", config.ConnString)
	if err != nil {
		return nil, err
	}

	s := restServer{
		config:   config,
		router:   e,
		database: conn,
	}

	// Configure HTTP ports
	s.ports.user, err = context.NewUserPort(conn, config)
	if err != nil {
		return nil, err
	}
	s.ports.post = context.NewPostPort(
		s.database,
		config.S3.Region,
		config.S3.Bucket,
		config.S3.Access,
		config.S3.Secret,
	)
	s.ports.tag = context.NewTagPort(s.database)

	// Start services
	token, err := context.NewTokenUsecase(config.Token.AccessSecret)
	if err != nil {
		return nil, err
	}

	s.setupRouter(token)
	return s, nil
}

func (s restServer) setupRouter(token service.TokenMaker) {
	s.router.Use(middleware.RemoveTrailingSlash())
	s.router.Validator = &Validator{val: validator.New()}

	s.router.POST("/user/login", s.ports.user.Login)
	s.router.POST("/user", s.ports.user.Register)

	authRoutes := s.router.Group("/*", http.AuthMiddleware(token))

	authRoutes.GET("/user/:id", s.ports.user.FindById)
	authRoutes.GET("/user", s.ports.user.GetMultiple)

	authRoutes.POST("/image", s.ports.post.Create)
	authRoutes.GET("/image/:id", s.ports.post.FindById)
	authRoutes.DELETE("/image/:id", s.ports.post.Delete)
	authRoutes.PATCH("/image", s.ports.post.Update)
	authRoutes.GET("/image", s.ports.post.GetMultiple)
	authRoutes.GET("/image/user/:id", s.ports.post.GetByUserId)

	authRoutes.GET("/tag", s.ports.tag.GetMultiple)
	authRoutes.GET("/tag/:id", s.ports.tag.GetUserTags)
}

func (s restServer) Start(address string) error {
	return s.router.Start(address)
}

func (s restServer) Close() error {
	return s.router.Close()
}
