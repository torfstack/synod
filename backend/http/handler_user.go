package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) LookUpUser(c echo.Context) error {
	searchString := c.QueryParam("find")
	if searchString == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no search string provided"})
	}

	return c.JSON(http.StatusOK, nil)
}
