package models

import (
	"crypto/rsa"
)

type Secret struct {
	ID    *int64   `json:"id,omitempty"`
	Value string   `json:"value"`
	Key   string   `json:"key"`
	Url   string   `json:"url"`
	Tags  []string `json:"tags"`
}

type EncryptedSecret Secret

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

type UserKeyPair struct {
	ID     *int64
	UserID int64
	KeyPair
}

type KeyPair struct {
	Public  rsa.PublicKey
	Private rsa.PrivateKey
}
