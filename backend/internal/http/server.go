package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"main/internal/auth"
	"main/internal/config"
	"main/internal/db"
	"os"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	e := echo.New()

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
	err := db.Migrate(context.Background(), cfg.DB.ConnectionString())
	if err != nil {
		panic(err)
	}

	database := db.NewDatabase(cfg.DB)

	authentication, err := auth.NewFireBaseAuth(context.Background(), database)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.CORS())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		println(err.Error())
		_ = c.JSON(500, map[string]string{"error": err.Error()})
	}

	e.GET("/secret", GetSecrets(database, authentication))
	e.POST("/secret", PostSecret(database, authentication))

	e.Logger.Fatal(e.Start(":4000"))
}
