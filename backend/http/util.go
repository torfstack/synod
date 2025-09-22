package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/torfstack/synod/backend/domain"

	"github.com/labstack/echo/v4"
)

const (
	SessionCookieName  = "sessionId"
	SessionContextName = "session"
)

func newEmptySessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		Value:    "",
		Expires:  time.UnixMilli(0),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

func newSessionCookie(sessionID string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		Value:    sessionID,
		Expires:  expiresAt,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

func getSessionIDCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func setSession(c echo.Context, session *domain.Session) {
	c.Set(SessionContextName, session)
}

func getSession(c echo.Context) (*domain.Session, bool) {
	session := c.Get(SessionContextName)
	if session == nil {
		return nil, false
	}
	return session.(*domain.Session), true
}

func authUrl(authBaseURL, clientID, redirectURL string) string {
	return fmt.Sprintf(
		"%s?client_id=%s&response_type=code&scope=openid+email+profile&redirect_uri=%s",
		authBaseURL,
		clientID,
		redirectURL,
	)
}

func doTokenRequest(tokenBaseURL, clientID, clientSecret, authCode, redirectURL string) (*http.Response, error) {
	values := make(url.Values)
	values.Add("code", authCode)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", redirectURL)
	r := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", tokenBaseURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)
	return http.DefaultClient.Do(req)
}
