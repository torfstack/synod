package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/logging"
	"golang.org/x/time/rate"

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
		c.SetRequest(c.Request().WithContext(logging.WithLogAttributeUserId(ctx, int(session.UserID))))
		return next(c)
	}
}

func newUnsealRateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:  rate.Limit(5.0 / 60),
			Burst: 5,
		}),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			session, ok := getSession(c)
			if !ok {
				return "", echo.NewHTTPError(http.StatusUnauthorized)
			}
			return strconv.FormatInt(session.UserID, 10), nil
		},
	})
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
