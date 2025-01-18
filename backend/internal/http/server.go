package http

import (
	"context"
	"github.com/labstack/echo/v4/middleware"
	"github.com/torfstack/kayvault/internal/auth"
	"github.com/torfstack/kayvault/internal/config"
	"github.com/torfstack/kayvault/internal/db"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

type Server struct {
	database       db.Database
	sessionService auth.SessionService
	firebaseAuth   auth.Auth
	cfg            config.Config
}

func NewServer() *Server {
	postgresPw := os.Getenv("POSTGRES_PASSWORD")
	if postgresPw == "" {
		panic("POSTGRES_PASSWORD environment variable not set")
	}
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	cfg := config.Config{
		DB: config.DBConfig{
			Host:     postgresHost,
			Port:     5432,
			User:     "postgres",
			Password: postgresPw,
			DBName:   "kayvault",
		},
	}

	return &Server{
		sessionService: auth.NewSessionService(),
		cfg:            cfg,
	}
}

func (s *Server) Start() {
	e := echo.New()

	err := db.Migrate(context.Background(), s.cfg.DB.ConnectionString())
	if err != nil {
		panic(err)
	}

	s.database = db.NewDatabase(s.cfg.DB)
	s.firebaseAuth, err = auth.NewFireBaseAuth(context.Background(), s.database)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		_ = c.JSON(500, map[string]string{"error": err.Error()})
	}

	var m echo.MiddlewareFunc
	if localMode == "enabled" {
		slog.Warn("Running in local mode")
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
	authorization.POST("", s.EstablishSession)
	authorization.GET("", s.IsAuthorized)
	authorization.DELETE("", s.EndSession)

	e.Logger.Fatal(e.Start(":4000"))
}

// localMode build flag, set with -ldflags -X 'github.com/torfstack/kayvault/internal/http.localMode=enabled'
var localMode string
