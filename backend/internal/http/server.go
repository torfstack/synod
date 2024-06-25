package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"main/internal/auth"
	"main/internal/config"
	"main/internal/db"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	e := echo.New()

	cfg := config.Config{DB: config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "mysecretpassword",
		DBName:   "kayvault",
	}}
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
