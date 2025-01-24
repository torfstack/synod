package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

func (s *Server) StartAuthentication(c echo.Context) error {
	provider, err := oidc.NewProvider(c.Request().Context(), s.cfg.Auth.Issuer)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(
		http.StatusFound,
		authUrl(provider.Endpoint().AuthURL, s.cfg.Auth.ClientID, s.cfg.Auth.RedirectURL),
	)
}

func (s *Server) EstablishSession(c echo.Context) error {
	code := c.QueryParam("code")

	clientId := s.cfg.Auth.ClientID
	redirectUrl := s.cfg.Auth.RedirectURL
	clientSecret := s.cfg.Auth.ClientSecret

	provider, err := oidc.NewProvider(
		c.Request().Context(),
		s.cfg.Auth.Issuer,
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res, err := doTokenRequest(provider.Endpoint().TokenURL, clientId, clientSecret, code, redirectUrl)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	var oidcResponse OidcResponse
	err = json.Unmarshal(resBytes, &oidcResponse)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := s.oidcAuth.GetUser(c.Request().Context(), oidcResponse.IdToken)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	session, err := s.sessionService.CreateSession(user.ID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.SetCookie(newSessionCookie(session.SessionID, session.ExpiresAt))
	return c.Redirect(http.StatusFound, s.cfg.Auth.BaseURL)
}

type OidcResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func (s *Server) IsAuthorized(c echo.Context) error {
	sessionID, err := getSessionIDCookie(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	_, err = s.sessionService.GetSession(sessionID)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) EndSession(c echo.Context) error {
	sessionID, err := getSessionIDCookie(c)
	if err != nil {
		return c.NoContent(http.StatusOK)
	}
	_ = s.sessionService.DeleteSession(sessionID)

	c.SetCookie(newEmptySessionCookie())
	return c.NoContent(http.StatusOK)
}
