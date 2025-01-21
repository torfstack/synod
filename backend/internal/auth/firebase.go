package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/torfstack/kayvault/internal/config"
	"github.com/torfstack/kayvault/internal/convert/fromdb"
	"github.com/torfstack/kayvault/internal/db"
	"github.com/torfstack/kayvault/internal/models"
)

type oidcAuth struct {
	database db.Database
	cfg      config.Config
}

func NewOidcAuth(ctx context.Context, database db.Database, cfg config.Config) (Auth, error) {
	return &oidcAuth{
		database: database,
		cfg:      cfg,
	}, nil
}

func (o *oidcAuth) GetUser(ctx context.Context, idToken string) (*models.User, error) {
	res, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	conn, err := o.database.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Error parsing claims as jwt.MapClaims")
		return nil, err
	}

	subject := claims["sub"].(string)
	exists, err := conn.Queries().DoesUserExist(ctx, subject)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = conn.Queries().InsertUser(ctx, subject)
		if err != nil {
			return nil, err
		}
	}

	dbUser, err := conn.Queries().SelectUserByName(ctx, subject)
	if err != nil {
		return nil, err
	}

	user := fromdb.User(dbUser)
	return &user, nil
}
