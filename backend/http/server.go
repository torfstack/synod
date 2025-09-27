package http

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	"github.com/torfstack/synod/backend/config"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/logging"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg           config.Config
	domainService domain.Service
}

func NewServer(cfg config.Config, domainService domain.Service) *Server {
	return &Server{
		cfg:           cfg,
		domainService: domainService,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		_ = c.JSON(500, map[string]string{"error": err.Error()})
	}

	var m echo.MiddlewareFunc
	if localMode == "enabled" {
		logging.Warnf(context.Background(), "Running in local mode")
		e.Use(
			middleware.CORSWithConfig(
				middleware.CORSConfig{
					AllowOrigins:     []string{s.cfg.Server.BaseURL},
					AllowCredentials: true,
				},
			),
		)
		m = s.LocalDevelopmentSession
	} else {
		m = s.SessionCheck
	}

	api := e.Group("/api")
	secrets := api.Group("/secrets", m)
	secrets.GET("", s.GetSecrets)
	secrets.POST("", s.PostSecret)

	authorization := api.Group("/auth")
	authorization.GET("/start", s.StartAuthentication)
	authorization.GET("/callback", s.EstablishSession)
	authorization.GET("", s.IsAuthorized)
	authorization.DELETE("", s.EndSession)

	setup := api.Group("/setup", m)
	setup.POST("/plain", s.PostSetupPlain)
	setup.POST("/password", s.PostSetupPassword)
	setup.POST("/unseal", s.UnsealWithPassword)

	users := api.Group("/users")
	users.GET("/lookup", s.LookUpUser)

	e.Static("/", "static")
	e.File("/", "static/index.html")

	return e.Start(fmt.Sprintf(":%d", s.cfg.Server.Port))
}

// localMode build flag, set with -ldflags "-X github.com/torfstack/synod/internal/http.localMode=enabled"
var localMode string
