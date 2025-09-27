package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/synod/backend/logging"
)

func (s *Server) UnsealWithPassword(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}

	err := s.domainService.UnsealWithPassword(ctx, session, "test")
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type UnsealRequest struct {
	Password string `json:"password"`
}

func (s *Server) PostSetupPlain(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}
	err := s.domainService.SetupUserPlain(ctx, *session)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (s *Server) PostSetupPassword(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in PostSetupPassword")
		return c.NoContent(http.StatusUnauthorized)
	}
	err := s.domainService.SetupUserWithPassword(ctx, *session, "test")
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
