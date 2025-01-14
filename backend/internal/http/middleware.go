package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) SessionCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie("token")
		if err != nil || cookie.Expires.Before(time.Now()) {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		session, err := s.sessionService.GetSession(cookie.Value)
		if err != nil {
			c.SetCookie(&http.Cookie{
				Name:    "token",
				Value:   "",
				Expires: time.UnixMilli(0),
			})
			c.Response().WriteHeader(401)
			return nil
		}

		c.Set("session", session)
		return next(c)
	}
}
