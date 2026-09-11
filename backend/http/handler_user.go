package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) LookUpUser(c echo.Context) error {
	ctx := c.Request().Context()
	session, ok := getSession(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	searchString := c.QueryParam("find")
	if len(searchString) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no search string provided"})
	}
	recipients, err := s.domainService.SearchShareRecipients(ctx, session.UserID, searchString)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, recipients)
}
