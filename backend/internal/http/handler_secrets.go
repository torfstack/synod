package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/kayvault/internal/convert/fromdb"
	"github.com/torfstack/kayvault/internal/convert/todb"
	"github.com/torfstack/kayvault/internal/models"
)

func (s *Server) GetSecrets(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	dbSecrets, err := s.database.SelectSecrets(c.Request().Context(), session.UserID)
	if err != nil {
		return err
	}

	secrets := fromdb.Secrets(dbSecrets)
	return c.JSON(http.StatusOK, secrets)
}

func (s *Server) PostSecret(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	var input models.Secret
	err := c.Bind(&input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if input.ID != 0 {
		err = s.database.UpdateSecret(c.Request().Context(), todb.UpdateSecretParams(input, session.UserID))
	} else {
		err = s.database.InsertSecret(c.Request().Context(), todb.InsertSecretParams(input, session.UserID))
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
