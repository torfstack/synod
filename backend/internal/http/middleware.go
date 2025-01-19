package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/torfstack/kayvault/internal/auth"

	"github.com/labstack/echo/v4"
)

func (s *Server) SessionCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie("sessionId")
		if err != nil {
			slog.Debug("No cookie found")
			return c.NoContent(http.StatusUnauthorized)
		}

		session, err := s.sessionService.GetSession(cookie.Value)
		if err != nil {
			slog.Debug("session not found")
			c.SetCookie(
				&http.Cookie{
					Name:     "sessionId",
					Value:    "",
					Expires:  time.UnixMilli(0),
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				},
			)
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("sessionId", session)
		return next(c)
	}
}

func (s *Server) LocalDevelopmentSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		slog.Debug("Local development session with id 1")
		c.Set(
			"sessionId", &auth.Session{
				SessionID: "local-development",
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		)
		return next(c)
	}
}
