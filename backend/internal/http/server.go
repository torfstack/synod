package http

import (
	"context"
	"main/internal/auth"
	"main/internal/config"
	"main/internal/db"
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

	/*	e.Use(
			middleware.CORSWithConfig(
				middleware.CORSConfig{
					AllowOrigins:     []string{"http://127.0.0.1:5173"},
					AllowCredentials: true,
				},
			),
		)
	*/
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		_ = c.JSON(500, map[string]string{"error": err.Error()})
	}

	secrets := e.Group("/secrets", s.SessionCheck)
	secrets.GET("", s.GetSecrets)
	secrets.POST("", s.PostSecret)

	e.POST("/auth", s.Auth)

	e.Logger.Fatal(e.Start(":4000"))
}
