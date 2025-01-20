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
	provider, err := oidc.NewProvider(
		c.Request().Context(),
		"https://accounts.google.com",
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	clientId := "301903522799-iibaqptt020lktncq9vh1f1n152sa2q0.apps.googleusercontent.com"
	e := provider.Endpoint()
	authUrl := fmt.Sprintf(
		"%s?client_id=%s&response_type=code&scope=openid+email+profile&redirect_uri=http://localhost:4000/api/callback",
		e.AuthURL,
		clientId,
	)
	return c.Redirect(http.StatusFound, authUrl)
}

func (s *Server) EstablishSession(c echo.Context) error {
	code := c.QueryParam("code")

	clientId := "301903522799-iibaqptt020lktncq9vh1f1n152sa2q0.apps.googleusercontent.com"
	clientSecret := "XXX"

	provider, err := oidc.NewProvider(
		c.Request().Context(),
		"https://accounts.google.com",
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	e := provider.Endpoint()
	values := make(url.Values)
	values.Add("code", code)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", "http://localhost:4000/api/callback")
	r := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", e.TokenURL, r)
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
	_ = json.Unmarshal(resBytes, &oidcResponse)
	fmt.Println(string(resBytes))

	fmt.Printf("Access token: %s, ID token: %s\n", oidcResponse.AccessToken, oidcResponse.IdToken)
	return c.Redirect(http.StatusFound, "http://localhost:5173")
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
			Value:    "",
			Expires:  time.UnixMilli(0),
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		},
	)

	return c.NoContent(http.StatusOK)
}
