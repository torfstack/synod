package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/torfstack/synod/backend/domain"
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

func (s *Server) GetThresholdParticipants(c echo.Context) error {
	session, secretID, ok := thresholdParticipantRequest(c)
	if !ok {
		return c.NoContent(http.StatusBadRequest)
	}
	participants, err := s.domainService.GetThresholdParticipants(
		c.Request().Context(), secretID, session.UserID,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, participants)
}

func (s *Server) AddThresholdParticipant(c echo.Context) error {
	session, secretID, ok := thresholdParticipantRequest(c)
	if !ok {
		return c.NoContent(http.StatusBadRequest)
	}
	var input models.ThresholdParticipantInput
	if err := c.Bind(&input); err != nil || input.SharingID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.AddThresholdParticipant(
		c.Request().Context(), secretID, session.UserID, input, session.Cipher,
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (s *Server) RemoveThresholdParticipant(c echo.Context) error {
	session, secretID, ok := thresholdParticipantRequest(c)
	if !ok || c.Param("sharingId") == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.RemoveThresholdParticipant(
		c.Request().Context(), secretID, session.UserID, c.Param("sharingId"), session.Cipher,
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) SetThresholdParticipantRole(c echo.Context) error {
	session, secretID, ok := thresholdParticipantRequest(c)
	if !ok || c.Param("sharingId") == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	var input models.ThresholdParticipantInput
	if err := c.Bind(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	input.SharingID = c.Param("sharingId")
	if err := s.domainService.SetThresholdParticipantRole(
		c.Request().Context(), secretID, session.UserID, input, session.Cipher,
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) SetThresholdParticipants(c echo.Context) error {
	session, secretID, ok := thresholdParticipantRequest(c)
	if !ok {
		return c.NoContent(http.StatusBadRequest)
	}
	var input models.ThresholdParticipantsInput
	if err := c.Bind(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := s.domainService.SetThresholdParticipants(
		c.Request().Context(), secretID, session.UserID, input, session.Cipher,
	); err != nil {
		if errors.Is(err, domain.ErrThresholdParticipantsChanged) {
			return c.String(http.StatusConflict, err.Error())
		}
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func thresholdParticipantRequest(c echo.Context) (*domain.Session, int64, bool) {
	session, ok := getSession(c)
	if !ok {
		return nil, 0, false
	}
	secretID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	return session, secretID, err == nil
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
		if errors.Is(err, domain.ErrUnlockAlreadyActive) {
			return c.NoContent(http.StatusConflict)
		}
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
