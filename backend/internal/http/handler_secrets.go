package http

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"main/internal/auth"
	"main/internal/convert/fromdb"
	"main/internal/convert/todb"
	"main/internal/models"
	"net/http"
)

func (s *Server) GetSecrets(c echo.Context) error {
	session, ok := c.Get("sessionId").(*auth.Session)
	if !ok {
		slog.Debug("No session found")
		return c.NoContent(http.StatusUnauthorized)
	}

	conn, err := s.database.Connect(c.Request().Context())
	if err != nil {
		return err
	}
	dbSecrets, err := conn.Queries().SelectSecrets(c.Request().Context(), session.UserID)
	if err != nil {
		return err
	}

	secrets := fromdb.Secrets(dbSecrets)
	return c.JSON(http.StatusOK, secrets)
}

func (s *Server) PostSecret(c echo.Context) error {
	session, ok := c.Get("sessionId").(*auth.Session)
	if !ok {
		slog.Debug("No session found")
		return c.NoContent(http.StatusUnauthorized)
	}

	conn, err := s.database.Connect(c.Request().Context())
	if err != nil {
		return err
	}

	var input models.Secret
	err = c.Bind(&input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = conn.Queries().InsertSecret(c.Request().Context(), todb.InsertSecretParams(input, session.UserID))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
