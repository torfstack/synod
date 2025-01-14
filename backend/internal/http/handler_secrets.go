package http

import (
	"github.com/labstack/echo/v4"
    "main/internal/auth"
	"main/internal/convert/fromdb"
	"main/internal/convert/todb"
	"main/internal/models"
)

func (s *Server)GetSecrets(c echo.Context) error {
    session, ok := c.Get("session").(auth.Session)
    if !ok {
        return echo.NewHTTPError(401, "Unauthorized")
    }   

    conn, err := s.database.Connect(c.Request().Context())
    if err != nil {
        return err
    }
    dbSecrets, err := conn.Queries().SelectSecrets(c.Request().Context(), int32(session.UserID))
    if err != nil {
        return err
    }

    secrets := fromdb.Secrets(dbSecrets)
    return c.JSON(200, secrets)
}

func (s *Server)PostSecret(c echo.Context) error {
    session, ok := c.Get("session").(auth.Session)
    if !ok {
        return echo.NewHTTPError(401, "Unauthorized")
    }   

    conn, err := s.database.Connect(c.Request().Context())
    if err != nil {
        return err
    }

    var input models.Secret
    err = c.Bind(&input)
    if err != nil {
        return err
    }

    err = conn.Queries().InsertSecret(c.Request().Context(), todb.InsertSecretParams(input, int32(session.UserID)))
    if err != nil {
        return err
    }

    return c.NoContent(201)
}
