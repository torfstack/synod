package http

import (
	"github.com/labstack/echo/v4"
	"github.com/torfstack/kayvault/internal/auth"
	"github.com/torfstack/kayvault/internal/convert/fromdb"
	"github.com/torfstack/kayvault/internal/convert/todb"
	"github.com/torfstack/kayvault/internal/models"
	"log/slog"
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

	if input.ID != 0 {
		err = conn.Queries().UpdateSecret(c.Request().Context(), todb.UpdateSecretParams(input, session.UserID))
	} else {
		err = conn.Queries().InsertSecret(c.Request().Context(), todb.InsertSecretParams(input, session.UserID))
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
