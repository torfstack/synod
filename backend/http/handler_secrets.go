package http

import (
	"net/http"

	"github.com/torfstack/kayvault/backend/logging"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/kayvault/backend/models"
)

func (s *Server) GetSecrets(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}

	secrets, err := s.domainService.GetSecrets(ctx, session.UserID)
	if err != nil {
		logging.Errorf(ctx, "could not retrieve secrets from DB: %v", err)
		return err
	}

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

	err = s.domainService.UpsertSecret(ctx, input, session.UserID)
	if err != nil {
		logging.Errorf(ctx, "could not insert/update secret: %v", err)
		return err
	}

	return c.NoContent(http.StatusCreated)
}
