package http

import (
	"net/http"
	"strconv"

	"github.com/torfstack/synod/backend/logging"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/synod/backend/models"
)

func (s *Server) GetSecrets(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		logging.Errorf(ctx, "no session found in GetSecrets")
		return c.NoContent(http.StatusUnauthorized)
	}

	secrets, err := s.domainService.GetSecrets(ctx, session.UserID, session.Cipher)
	if err != nil {
		logging.Errorf(ctx, "could not retrieve secrets from DB: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, secrets)
}

type shareSecretRequest struct {
	SharingID string `json:"sharingId"`
}

func (s *Server) GetSecretRecipients(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	secretID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	recipients, err := s.domainService.GetSecretRecipients(ctx, secretID, session.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, recipients)
}

func (s *Server) ShareSecret(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	secretID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	var input shareSecretRequest
	if err := c.Bind(&input); err != nil || input.SharingID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.ShareSecret(ctx, secretID, session.UserID, input.SharingID, session.Cipher); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (s *Server) RevokeSecretAccess(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	secretID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.RevokeSecretAccess(ctx, secretID, session.UserID, c.Param("sharingId")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
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

	_, err = s.domainService.UpsertSecret(ctx, input, session.UserID, session.Cipher)
	if err != nil {
		logging.Errorf(ctx, "could not insert/update secret: %v", err)
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) PostThresholdSecret(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	var input models.ThresholdSecretInput
	if err := c.Bind(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id, err := s.domainService.CreateThresholdSecret(c.Request().Context(), input, session.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) StartUnlock(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	secretID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id, err := s.domainService.StartUnlock(c.Request().Context(), secretID, session.UserID, session.Cipher)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) GetUnlockRequests(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	requests, err := s.domainService.GetUnlockRequests(c.Request().Context(), session.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, requests)
}

func (s *Server) ContributeToUnlock(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.ContributeToUnlock(
		c.Request().Context(),
		requestID,
		session.UserID,
		session.Cipher,
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetUnlockResult(c echo.Context) error {
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	result, err := s.domainService.GetUnlockResult(c.Request().Context(), requestID, session.UserID, session.Cipher)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
