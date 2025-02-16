package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) LookUpUser(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
