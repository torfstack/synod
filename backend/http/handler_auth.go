package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

func (s *Server) StartAuthentication(c echo.Context) error {
	ctx := c.Request().Context()
	provider, err := oidc.NewProvider(ctx, s.cfg.Auth.Issuer)
	if err != nil {
		logging.Errorf(ctx, "could not create oidc provider from discovery url: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(
		http.StatusFound,
		authUrl(provider.Endpoint().AuthURL, s.cfg.Auth.ClientID, s.cfg.Auth.RedirectURL),
	)
}

func (s *Server) EstablishSession(c echo.Context) error {
	ctx := c.Request().Context()
	code := c.QueryParam("code")

	clientId := s.cfg.Auth.ClientID
	redirectUrl := s.cfg.Auth.RedirectURL
	clientSecret := s.cfg.Auth.ClientSecret

	provider, err := oidc.NewProvider(ctx, s.cfg.Auth.Issuer)
	if err != nil {
		logging.Errorf(ctx, "could not create oidc provider from discovery url: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	res, err := doTokenRequest(provider.Endpoint().TokenURL, clientId, clientSecret, code, redirectUrl)
	if err != nil {
		logging.Errorf(ctx, "could not perform token request: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logging.Errorf(ctx, "could not read token response from oidc provider: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var oidcResponse OidcResponse
	err = json.Unmarshal(resBytes, &oidcResponse)
	if err != nil {
		logging.Errorf(ctx, "could not unmarshal token response from oidc provider: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := s.domainService.GetUserFromToken(ctx, oidcResponse.IdToken)
	if err != nil {
		logging.Errorf(ctx, "could not get user based on id token from oidc provider: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	session, err := s.domainService.CreateSession(ctx, user.ID)
	if err != nil {
		logging.Errorf(ctx, "could not create session: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	c.SetCookie(newSessionCookie(session.SessionID, session.ExpiresAt))
	return c.Redirect(http.StatusFound, s.cfg.Server.BaseURL)
}

type OidcResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func (s *Server) IsAuthorized(c echo.Context) error {
	ctx := c.Request().Context()
	sessionID, err := getSessionIDCookie(c)
	if err != nil {
		logging.Debugf(ctx, "no sessionId cookie found")
		return c.NoContent(http.StatusUnauthorized)
	}

	session, err := s.domainService.GetSession(sessionID)
	if errors.Is(err, domain.ErrSessionNotFound) {
		return c.JSON(http.StatusUnauthorized, models.AuthStatus{IsAuthenticated: false})
	}

	isUserSetup, err := s.domainService.IsUserSetup(ctx, *session)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.AuthStatus{
		IsAuthenticated: true,
		IsSetup:         isUserSetup,
		NeedsToUnseal:   session.PrivateKey == nil,
	})
}

func (s *Server) EndSession(c echo.Context) error {
	ctx := c.Request().Context()
	sessionID, err := getSessionIDCookie(c)
	if err != nil {
		logging.Debugf(ctx, "no sessionId cookie found")
		return c.NoContent(http.StatusOK)
	}
	_ = s.domainService.DeleteSession(sessionID)

	c.SetCookie(newEmptySessionCookie())
	return c.NoContent(http.StatusOK)
}
