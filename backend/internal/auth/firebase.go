package auth

import (
	"context"
	"main/internal/convert/fromdb"
	"main/internal/db"
	"main/internal/models"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type fireBaseAuth struct {
	database db.Database
	auth     *auth.Client
}

func NewFireBaseAuth(ctx context.Context, database db.Database) (Auth, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("kayvault.json"))
	if err != nil {
		return nil, err
	}
	a, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &fireBaseAuth{
		database: database,
		auth:     a,
	}, nil
}

func (f *fireBaseAuth) GetUser(ctx context.Context, token string) (*models.User, error) {
	token = strings.TrimPrefix(token, "Bearer")
	res, err := f.auth.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	conn, err := f.database.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	exists, err := conn.Queries().DoesUserExist(ctx, res.Subject)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = conn.Queries().InsertUser(ctx, res.Subject)
		if err != nil {
			return nil, err
		}
	}

	dbUser, err := conn.Queries().SelectUserByName(ctx, res.Subject)
	if err != nil {
		return nil, err
	}

	user := fromdb.User(dbUser)
	return &user, nil
}
