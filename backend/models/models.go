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

type KeyType int

const (
	KeyTypeRsa KeyType = iota + 1
)

type UserKeyPair struct {
	ID         *int64
	UserID     int64
	Type       KeyType
	PasswordID *int64
	Public     []byte
	Private    []byte
}

type KeyPair struct {
	Public  rsa.PublicKey
	Private rsa.PrivateKey
}

type HashedPassword struct {
	ID         *int64
	Hash       []byte
	Salt       []byte
	Iterations int64
}

type AuthStatus struct {
	IsAuthenticated bool `json:"isAuthenticated"`
	IsSetup         bool `json:"isSetup"`
	NeedsToUnseal   bool `json:"needsToUnseal"`
}
