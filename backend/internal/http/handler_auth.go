package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) Auth(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	user, err := s.firebaseAuth.GetUser(c.Request().Context(), auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	session, err := s.sessionService.CreateSession(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
	})

	return nil
}
