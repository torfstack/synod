package domain

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/models"
)

var _ UserService = &service{}

func (s *service) DoesUserExist(ctx context.Context, username string) (bool, error) {
	return s.database.DoesUserExist(ctx, username)
}

func (s *service) InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error) {
	createdUser, err := s.database.InsertUser(ctx, user)
	if err != nil {
		return models.ExistingUser{}, err
	}
	return createdUser, err
}

func (s *service) GetUserFromToken(ctx context.Context, idToken *oidc.IDToken) (models.ExistingUser, error) {
	var c claims
	err := idToken.Claims(&c)
	if err != nil {
		return models.ExistingUser{}, fmt.Errorf("oidc: failed to decode claims: %w", err)
	}
	userParams := models.User{
		Subject:  c.Subject,
		Email:    c.Email,
		FullName: c.Name,
	}

	var user models.ExistingUser
	err = s.database.WithTx(
		ctx, func(d db.Database) error {
			var exists bool
			exists, err = d.DoesUserExist(ctx, c.Subject)
			if err != nil {
				return err
			}
			if !exists {
				_, err = s.InsertUser(ctx, userParams)
				if err != nil {
					return err
				}
			}
			user, err = d.SelectUserByName(ctx, userParams.Subject)
			if err != nil {
				return err
			}
			return err
		},
	)
	return user, err
}

type claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}
