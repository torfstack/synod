package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/logging"
)

func (s *Server) UnsealWithPassword(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in UnsealWithPassword")
		return c.NoContent(http.StatusUnauthorized)
	}

	var input UnsealRequest
	err := c.Bind(&input)
	if err != nil {
		logging.Errorf(ctx, "input could not be parsed to UnsealRequest: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	password := crypto.Password(input.Password)
	defer password.Zero()
	err = s.domainService.UnsealWithPassword(ctx, session, password)
	if err != nil {
		logging.Errorf(ctx, "could not unseal vault: %v", err)
		if errors.Is(err, domain.ErrInvalidPassword) {
			return c.NoContent(http.StatusForbidden)
		}
		return c.NoContent(http.StatusUnauthorized)
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
		logging.Errorf(ctx, "no session found in PostSetupPlain")
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

	setupPassword := crypto.Password(input.Password)
	defer setupPassword.Zero()
	err = s.domainService.SetupUserWithPassword(ctx, *session, setupPassword)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
