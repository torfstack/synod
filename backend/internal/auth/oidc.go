package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/torfstack/kayvault/internal/config"
	"github.com/torfstack/kayvault/internal/convert/fromdb"
	"github.com/torfstack/kayvault/internal/db"
	"github.com/torfstack/kayvault/internal/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
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

	userParams, err := insertUserParamsFromClaims(claims)
	if err != nil {
		return nil, err
	}

	d, t := o.database.WithTx(ctx)
	defer t.Rollback(ctx)
	exists, err := d.DoesUserExist(ctx, userParams.Subject)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = d.InsertUser(ctx, userParams)
		if err != nil {
			return nil, err
		}
	}
	dbUser, err := d.SelectUserByName(ctx, userParams.Subject)
	if err != nil {
		return nil, err
	}
	t.Commit(ctx)

	user := fromdb.User(dbUser)
	return &user, nil
}

func insertUserParamsFromClaims(claims jwt.MapClaims) (sqlc.InsertUserParams, error) {
	subject, ok := claims["sub"]
	if !ok {
		return sqlc.InsertUserParams{}, errors.New("error parsing subject as string")
	}
	email, ok := claims["email"]
	if !ok {
		email = ""
	}
	fullName, ok := claims["name"]
	if !ok {
		fullName = ""
	}
	return sqlc.InsertUserParams{
		Subject:  subject.(string),
		Email:    email.(string),
		FullName: fullName.(string),
	}, nil
}
