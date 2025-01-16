package http

import (
	"log/slog"
	"net/http"
	"time"

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
			c.SetCookie(
				&http.Cookie{
					Name:    "sessionId",
					Value:   "",
					Expires: time.UnixMilli(0),
				},
			)
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("sessionId", session)
		return next(c)
	}
}
