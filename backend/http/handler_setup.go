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

	var input UnsealRequest
	err := c.Bind(&input)
	if err != nil {
		logging.Errorf(ctx, "input could not be parsed to UnsealRequest: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.domainService.UnsealWithPassword(ctx, session, input.Password)
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

type SetupPasswordRequest struct {
	Password string `json:"password"`
}

func (s *Server) PostSetupPassword(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in PostSetupPassword")
		return c.NoContent(http.StatusUnauthorized)
	}

	var input SetupPasswordRequest
	err := c.Bind(&input)
	if err != nil {
		logging.Errorf(ctx, "input could not be parsed to SetupPasswordRequest: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.domainService.SetupUserWithPassword(ctx, *session, input.Password)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
