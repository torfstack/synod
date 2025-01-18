package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) EstablishSession(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return c.NoContent(http.StatusUnauthorized)
	}

	user, err := s.firebaseAuth.GetUser(c.Request().Context(), auth)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	session, err := s.sessionService.CreateSession(user.ID)
	if err != nil {
		return err
	}

	c.SetCookie(
		&http.Cookie{
			Name:     "sessionId",
			Value:    session.SessionID,
			Expires:  session.ExpiresAt,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		},
	)

	return c.NoContent(http.StatusOK)
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
			Name:    "sessionId",
			Value:   "",
			Expires: time.UnixMilli(0),
		},
	)

	return c.NoContent(http.StatusOK)
}
