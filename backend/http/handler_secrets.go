package http

import (
	"github.com/torfstack/kayvault/backend/logging"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/kayvault/backend/convert/fromdb"
	"github.com/torfstack/kayvault/backend/convert/todb"
	"github.com/torfstack/kayvault/backend/models"
)

func (s *Server) GetSecrets(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}

	dbSecrets, err := s.database.SelectSecrets(ctx, session.UserID)
	if err != nil {
		logging.Errorf(ctx, "could not retrieve secrets from DB: %v", err)
		return err
	}

	secrets := fromdb.Secrets(dbSecrets)
	return c.JSON(http.StatusOK, secrets)
}

func (s *Server) PostSecret(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}

	var input models.Secret
	err := c.Bind(&input)
	if err != nil {
		logging.Errorf(ctx, "input could not be parsed to models.Secret: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if input.ID != 0 {
		logging.Debugf(ctx, "updating secret with id %d", input.ID)
		err = s.database.UpdateSecret(ctx, todb.UpdateSecretParams(input, session.UserID))
	} else {
		logging.Debugf(ctx, "inserting new secret")
		err = s.database.InsertSecret(ctx, todb.InsertSecretParams(input, session.UserID))
	}
	if err != nil {
		logging.Errorf(ctx, "could not insert/update secret: %v", err)
		return err
	}

	return c.NoContent(http.StatusCreated)
}
