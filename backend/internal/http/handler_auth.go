package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

func (s *Server) StartAuthentication(c echo.Context) error {
	provider, err := oidc.NewProvider(c.Request().Context(),
		s.cfg.Auth.Issuer,
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	clientId := s.cfg.Auth.ClientID
	redirectUrl := s.cfg.Auth.RedirectURL
	authUrl := fmt.Sprintf(
		"%s?client_id=%s&response_type=code&scope=openid+email+profile&redirect_uri=%s",
		provider.Endpoint().AuthURL,
		clientId,
		redirectUrl,
	)
	return c.Redirect(http.StatusFound, authUrl)
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
	values := make(url.Values)
	values.Add("code", code)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", redirectUrl)
	r := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", provider.Endpoint().TokenURL, r)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientId, clientSecret)
	res, err := http.DefaultClient.Do(req)
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

	c.SetCookie(&http.Cookie{
		Name:     "sessionId",
		Path:     "/",
		Value:    session.SessionID,
		Expires:  session.ExpiresAt,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	})

	return c.Redirect(http.StatusFound, s.cfg.Auth.BaseURL)
}

type OidcResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func (s *Server) IsAuthorized(c echo.Context) error {
	sessionID, err := c.Cookie("sessionId")
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	_, err = s.sessionService.GetSession(sessionID.Value)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) EndSession(c echo.Context) error {
	sessionID, err := c.Cookie("sessionId")
	if err != nil {
		return c.NoContent(http.StatusOK)
	}
	_ = s.sessionService.DeleteSession(sessionID.Value)

	c.SetCookie(
		&http.Cookie{
			Name:     "sessionId",
			Path:     "/",
			Value:    "",
			Expires:  time.UnixMilli(0),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		},
	)

	return c.NoContent(http.StatusOK)
}
