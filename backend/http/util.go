package http

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	PKCECookieName     = "pkce_verifier"
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

func generatePKCE() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return
}

func newPKCECookie(verifier string, expiry time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     PKCECookieName,
		Path:     "/api/auth",
		Value:    verifier,
		Expires:  expiry,
		SameSite: http.SameSiteLaxMode, // Lax required: cookie must arrive on the OIDC redirect
		HttpOnly: true,
		Secure:   true,
	}
}

func newEmptyPKCECookie() *http.Cookie {
	return &http.Cookie{
		Name:     PKCECookieName,
		Path:     "/api/auth",
		Value:    "",
		Expires:  time.UnixMilli(0),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
}

func authUrl(authBaseURL, clientID, redirectURL, codeChallenge string) string {
	return fmt.Sprintf(
		"%s?client_id=%s&response_type=code&scope=openid+email+profile&redirect_uri=%s&code_challenge=%s&code_challenge_method=S256",
		authBaseURL,
		clientID,
		redirectURL,
		codeChallenge,
	)
}

func doTokenRequest(tokenBaseURL, clientID, clientSecret, authCode, redirectURL, codeVerifier string) (
	*http.Response, error,
) {
	values := make(url.Values)
	values.Add("code", authCode)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", redirectURL)
	values.Add("code_verifier", codeVerifier)
	r := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", tokenBaseURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)
	return http.DefaultClient.Do(req)
}
