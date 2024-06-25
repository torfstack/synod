package http

import (
	"github.com/labstack/echo/v4"
	"main/internal/auth"
	"main/internal/convert/fromdb"
	"main/internal/convert/todb"
	"main/internal/db"
	"main/internal/models"
)

func GetSecrets(db db.Database, auth auth.Auth) func(echo.Context) error {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		user, err := auth.GetUser(c.Request().Context(), authHeader)
		if err != nil {
			c.Error(err)
			return err
		}

		conn, err := db.Connect(c.Request().Context())
		if err != nil {
			return err
		}
		dbSecrets, err := conn.Queries().SelectSecrets(c.Request().Context(), user.ID)
		if err != nil {
			return err
		}

		secrets := fromdb.Secrets(dbSecrets)
		return c.JSON(200, secrets)
	}
}

func PostSecret(db db.Database, auth auth.Auth) func(echo.Context) error {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		user, err := auth.GetUser(c.Request().Context(), authHeader)
		if err != nil {
			return err
		}

		conn, err := db.Connect(c.Request().Context())
		if err != nil {
			return err
		}

		var input models.Secret
		c.Bind(&input)

		err = conn.Queries().InsertSecret(c.Request().Context(), todb.InsertSecretParams(input, user.ID))
		if err != nil {
			return err
		}

		return c.NoContent(201)
	}
}
