package http

import (
	"context"
	"fmt"
	"time"

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

	loggerMiddleware := middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogStatus:  true,
			LogURI:     true,
			LogMethod:  true,
			LogLatency: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				fmt.Printf(
					"[%s] method=%s, uri=%s, status=%d, latency=%s\n",
					v.StartTime.Format(time.RFC3339), v.Method, v.URI, v.Status, v.Latency.String(),
				)
				return nil
			},
		},
	)

	api := e.Group("/api")
	secrets := api.Group("/secrets", m, loggerMiddleware)
	secrets.GET("", s.GetSecrets)
	secrets.POST("", s.PostSecret)

	authorization := api.Group("/auth", loggerMiddleware)
	authorization.GET("/start", s.StartAuthentication)
	authorization.GET("/callback", s.EstablishSession)
	authorization.GET("", s.IsAuthorized)
	authorization.DELETE("", s.EndSession)

	setup := api.Group("/setup", m, loggerMiddleware)
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
