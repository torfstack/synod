package models

import (
	"crypto/ed25519"
	"errors"
)

type Secret struct {
	ID    *int64   `json:"id,omitempty"`
	Value string   `json:"value"`
	Key   string   `json:"key"`
	Url   string   `json:"url"`
	Tags  []string `json:"tags"`
}

type Secrets []Secret

type User struct {
	ID       *int64 `json:"id,omitempty"`
	Subject  string `json:"subject"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type ExistingUser struct {
	User
	ID int64 `json:"id"`
}

func NewExistingUser(user User) (ExistingUser, error) {
	if user.ID == nil {
		return ExistingUser{}, errors.New("user ID cannot be nil")
	}
	return ExistingUser{
		User: user,
		ID:   *user.ID,
	}, nil
}

type UserKeyPair struct {
	ID     *int64
	UserID int64
	KeyPair
}

type KeyPair struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}
