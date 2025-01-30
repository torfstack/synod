package http

import (
	"github.com/torfstack/kayvault/internal/logging"
	"net/http"
	"time"

	"github.com/torfstack/kayvault/internal/auth"

	"github.com/labstack/echo/v4"
)

func (s *Server) SessionCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cookie, err := getSessionIDCookie(c)
		if err != nil {
			logging.Debugf(ctx, "No sessionId cookie found")
			return c.NoContent(http.StatusUnauthorized)
		}

		session, err := s.sessionService.GetSession(cookie)
		if err != nil {
			logging.Debugf(ctx, "Could not get session: %v", err)
			c.SetCookie(newEmptySessionCookie())
			return c.NoContent(http.StatusUnauthorized)
		}

		setSession(c, session)
		logging.WithLogAttributeUserId(c.Request().Context(), int(session.UserID))
		return next(c)
	}
}

func (s *Server) LocalDevelopmentSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		setSession(
			c, &auth.Session{
				SessionID: "local-development",
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		)
		return next(c)
	}
}
