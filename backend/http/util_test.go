package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/domain"
)

// ---- Cookie helpers ----

func TestNewEmptySessionCookie(t *testing.T) {
	c := newEmptySessionCookie()
	assert.Equal(t, SessionCookieName, c.Name)
	assert.Equal(t, "", c.Value)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
	assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
	assert.Equal(t, time.UnixMilli(0), c.Expires)
}

func TestNewSessionCookie(t *testing.T) {
	expiry := time.Now().Add(time.Hour).Truncate(time.Second)
	c := newSessionCookie("abc-123", expiry)
	assert.Equal(t, SessionCookieName, c.Name)
	assert.Equal(t, "abc-123", c.Value)
	assert.Equal(t, expiry, c.Expires)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
}

func TestNewEmptyPKCECookie(t *testing.T) {
	c := newEmptyPKCECookie()
	assert.Equal(t, PKCECookieName, c.Name)
	assert.Equal(t, "", c.Value)
	assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
}

func TestNewPKCECookie(t *testing.T) {
	expiry := time.Now().Add(10 * time.Minute).Truncate(time.Second)
	c := newPKCECookie("verifier-xyz", expiry)
	assert.Equal(t, PKCECookieName, c.Name)
	assert.Equal(t, "verifier-xyz", c.Value)
	assert.Equal(t, expiry, c.Expires)
	assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
}

// ---- authUrl ----

func TestAuthUrl_ContainsAllParams(t *testing.T) {
	u := authUrl("https://auth.example.com/authorize", "client-id", "https://app/callback", "challenge-abc")
	assert.Contains(t, u, "https://auth.example.com/authorize")
	assert.Contains(t, u, "client_id=client-id")
	assert.Contains(t, u, "response_type=code")
	assert.Contains(t, u, "redirect_uri=https%3A%2F%2Fapp%2Fcallback")
	assert.Contains(t, u, "code_challenge=challenge-abc")
	assert.Contains(t, u, "code_challenge_method=S256")
}

// ---- generatePKCE ----

func TestGeneratePKCE_ProducesNonEmptyValues(t *testing.T) {
	verifier, challenge, err := generatePKCE()
	require.NoError(t, err)
	assert.NotEmpty(t, verifier)
	assert.NotEmpty(t, challenge)
	assert.NotEqual(t, verifier, challenge)
}

func TestGeneratePKCE_ProducesUniquePairs(t *testing.T) {
	v1, c1, err := generatePKCE()
	require.NoError(t, err)
	v2, c2, err := generatePKCE()
	require.NoError(t, err)

	assert.NotEqual(t, v1, v2)
	assert.NotEqual(t, c1, c2)
}

// ---- getSessionIDCookie ----

func TestGetSessionIDCookie_Present_ReturnsValue(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "my-session"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	val, err := getSessionIDCookie(c)
	require.NoError(t, err)
	assert.Equal(t, "my-session", val)
}

func TestGetSessionIDCookie_Missing_ReturnsError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, err := getSessionIDCookie(c)
	require.Error(t, err)
}

// ---- setSession / getSession ----

func TestSetAndGetSession_RoundTrip(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	sess := &domain.Session{SessionID: "s1", UserID: 7}
	setSession(c, sess)

	got, ok := getSession(c)
	require.True(t, ok)
	assert.Equal(t, sess, got)
}

func TestGetSession_NotSet_ReturnsFalse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	got, ok := getSession(c)
	assert.False(t, ok)
	assert.Nil(t, got)
}

// ---- doTokenRequest builds the right HTTP request ----

func TestDoTokenRequest_SendsPostWithBasicAuth(t *testing.T) {
	// We use a test server to inspect the request received.
	var gotMethod, gotContentType, gotAuthHeader, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		gotAuthHeader = r.Header.Get("Authorization")
		buf := new(strings.Builder)
		b := make([]byte, 4096)
		for {
			n, err := r.Body.Read(b)
			if n > 0 {
				buf.Write(b[:n])
			}
			if err != nil {
				break
			}
		}
		gotBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := doTokenRequest(ts.URL, "my-client", "my-secret", "auth-code", "https://redirect", "code-verifier")
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "application/x-www-form-urlencoded", gotContentType)
	assert.True(t, strings.HasPrefix(gotAuthHeader, "Basic "), "should use Basic auth")
	assert.Contains(t, gotBody, "code=auth-code")
	assert.Contains(t, gotBody, "grant_type=authorization_code")
	assert.Contains(t, gotBody, "redirect_uri=")
	assert.Contains(t, gotBody, "code_verifier=code-verifier")
}
