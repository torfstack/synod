package auth

import (
	"context"
	"errors"
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

func NewOidcAuth(database db.Database, cfg config.Config) (Auth, error) {
	return &oidcAuth{
		database: database,
		cfg:      cfg,
	}, nil
}

func (o *oidcAuth) GetUser(ctx context.Context, idToken string) (*models.User, error) {
	// TODO: Implement proper JWT validation
	res, _ := jwt.Parse(
		idToken, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		},
	)

	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("error parsing claims as jwt.MapClaims")
	}

	subject := claims["sub"].(string)

	d, t := o.database.WithTx(ctx)
	defer t.Rollback(ctx)
	exists, err := d.DoesUserExist(ctx, subject)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = d.InsertUser(ctx, subject)
		if err != nil {
			return nil, err
		}
	}
	dbUser, err := d.SelectUserByName(ctx, subject)
	if err != nil {
		return nil, err
	}
	t.Commit(ctx)

	user := fromdb.User(dbUser)
	return &user, nil
}
