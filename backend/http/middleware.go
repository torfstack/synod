package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/logging"

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

		session, err := s.domainService.GetSession(cookie)
		if err != nil {
			logging.Debugf(ctx, "Could not get session: %v", err)
			c.SetCookie(newEmptySessionCookie())
			return c.NoContent(http.StatusUnauthorized)
		}

		setSession(c, session)
		logging.WithLogAttributeUserId(ctx, int(session.UserID))
		return next(c)
	}
}

func (s *Server) LocalDevelopmentSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		setSession(
			c, &domain.Session{
				SessionID: "local-development",
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		)
		return next(c)
	}
}

func (s *Server) RequestAndResponseLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request() != nil {
			req := c.Request()
			logging.Debugf(
				req.Context(),
				fmt.Sprintf("[%s] %s", req.Method, req.RequestURI),
			)
		}
		if c.Response() != nil {
			res := c.Response()
			req := c.Request()
			logging.Debugf(
				req.Context(),
				fmt.Sprintf("[%s] %s returned %d", req.Method, req.RequestURI, res.Status),
			)
		}
		return next(c)
	}
}
