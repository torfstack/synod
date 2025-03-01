package domain

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/torfstack/kayvault/backend/convert/fromdb"
	"github.com/torfstack/kayvault/backend/convert/todb"
	"github.com/torfstack/kayvault/backend/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
)

var _ UserService = &service{}

func (s *service) DoesUserExist(ctx context.Context, username string) (bool, error) {
	return s.database.DoesUserExist(ctx, username)
}

func (s *service) InsertUser(ctx context.Context, user models.User) error {
	return s.database.InsertUser(ctx, todb.InsertUserParams(user))
}

func (s *service) GetUserFromToken(ctx context.Context, idToken string) (*models.User, error) {
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

	d, t := s.database.WithTx(ctx)
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
