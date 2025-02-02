package http

import (
	"context"
	"github.com/labstack/echo/v4/middleware"
	"github.com/torfstack/kayvault/internal/auth"
	"github.com/torfstack/kayvault/internal/config"
	"github.com/torfstack/kayvault/internal/db"
	"github.com/torfstack/kayvault/internal/logging"

	"github.com/labstack/echo/v4"
)

type Server struct {
	database       db.Database
	sessionService auth.SessionService
	oidcAuth       auth.Auth
	cfg            config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		sessionService: auth.NewSessionService(),
		cfg:            cfg,
	}
}

func (s *Server) Start() {
	e := echo.New()

	err := db.Migrate(context.Background(), s.cfg.DB.ConnectionString())
	if err != nil {
		logging.Fatalf(context.Background(), "could not migrate database: %v", err)
		return
	}

	s.database = db.NewDatabase(s.cfg.DB.ConnectionString())
	s.oidcAuth, err = auth.NewOidcAuth(s.database, s.cfg)
	if err != nil {
		logging.Fatalf(context.Background(), "could not create oidc auth: %v", err)
		return
	}

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
					AllowOrigins:     []string{"http://localhost:5173"},
					AllowCredentials: true,
				},
			),
		)
		m = s.LocalDevelopmentSession
	} else {
		m = s.SessionCheck
	}

	secrets := e.Group("/secrets", m)
	secrets.GET("", s.GetSecrets)
	secrets.POST("", s.PostSecret)

	authorization := e.Group("/auth")
	authorization.GET("/start", s.StartAuthentication)
	authorization.GET("/callback", s.EstablishSession)
	authorization.GET("", s.IsAuthorized)
	authorization.DELETE("", s.EndSession)

	e.Logger.Fatal(e.Start(":4000"))
}

// localMode build flag, set with -ldflags -X 'github.com/torfstack/kayvault/internal/http.localMode=enabled'
var localMode string
