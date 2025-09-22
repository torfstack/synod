package domain

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/torfstack/synod/backend/models"
)

var _ UserService = &service{}

func (s *service) DoesUserExist(ctx context.Context, username string) (bool, error) {
	return s.database.DoesUserExist(ctx, username)
}

func (s *service) InsertUser(ctx context.Context, user models.User) error {
	return s.database.InsertUser(ctx, user)
}

func (s *service) GetUserFromToken(ctx context.Context, idToken string) (models.ExistingUser, error) {
	// TODO: Implement proper JWT validation
	res, _ := jwt.Parse(
		idToken, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		},
	)

	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok {
		return models.ExistingUser{}, errors.New("error parsing claims as jwt.MapClaims")
	}

	userParams, err := userFromClaims(claims)
	if err != nil {
		return models.ExistingUser{}, err
	}

	d, t := s.database.WithTx(ctx)
	defer t.Rollback(ctx)
	exists, err := d.DoesUserExist(ctx, userParams.Subject)
	if err != nil {
		return models.ExistingUser{}, err
	}
	if !exists {
		err = d.InsertUser(ctx, userParams)
		if err != nil {
			return models.ExistingUser{}, err
		}
	}
	user, err := d.SelectUserByName(ctx, userParams.Subject)
	if err != nil {
		return models.ExistingUser{}, err
	}
	t.Commit(ctx)

	return models.NewExistingUser(user)
}

func userFromClaims(claims jwt.MapClaims) (models.User, error) {
	subject, ok := claims["sub"]
	if !ok {
		return models.User{}, errors.New("error parsing subject as string")
	}
	email, ok := claims["email"]
	if !ok {
		email = ""
	}
	fullName, ok := claims["name"]
	if !ok {
		fullName = ""
	}
	return models.User{
		Subject:  subject.(string),
		Email:    email.(string),
		FullName: fullName.(string),
	}, nil
}
