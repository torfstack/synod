package models

import (
	"crypto/rsa"
	"time"
)

type Secret struct {
	ID        *int64   `json:"id,omitempty"`
	Value     string   `json:"value"`
	Key       string   `json:"key"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	Owned     bool     `json:"owned"`
	Locked    bool     `json:"locked,omitempty"`
	Threshold int      `json:"threshold,omitempty"`
}

type ThresholdSecretInput struct {
	Secret     Secret   `json:"secret"`
	Threshold  int      `json:"threshold"`
	SharingIDs []string `json:"sharingIds"`
}

type ThresholdSecret struct {
	ID               int64
	OwnerID          int64
	EncryptedPayload []byte
	EncryptedShare   []byte
	Threshold        int
	Key              string
	Url              string
	Tags             []string
}

type UnlockRequest struct {
	ID            int64     `json:"id"`
	SecretID      int64     `json:"secretId"`
	RequesterID   int64     `json:"-"`
	RequesterName string    `json:"requesterName"`
	SecretName    string    `json:"secretName"`
	Threshold     int       `json:"threshold"`
	Contributions int       `json:"contributions"`
	Contributed   bool      `json:"contributed"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

type UnlockResult struct {
	Ready         bool    `json:"ready"`
	Contributions int     `json:"contributions"`
	Threshold     int     `json:"threshold"`
	Secret        *Secret `json:"secret,omitempty"`
}

type EncryptedSecret Secret

type Secrets []Secret

type AccessibleSecret struct {
	ID               int64
	OwnerID          int64
	EncryptedPayload []byte
	EncryptedDataKey []byte
	Envelope         bool
	Owned            bool
	Legacy           EncryptedSecret
}

type ShareRecipient struct {
	ID       int64  `json:"-"`
	Subject  string `json:"sharingId"`
	Email    string `json:"emailHint"`
	FullName string `json:"fullName"`
}

type User struct {
	ID       *int64 `json:"id,omitempty"`
	Subject  string `json:"subject"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type ExistingUser struct {
	User
	ID        int64  `json:"id"`
	SharingID string `json:"sharingId"`
}

type KeyType int

const (
	KeyTypeHPKE KeyType = iota + 1
	KeyTypeRsa          = KeyTypeHPKE
)

type UserKeyPair struct {
	ID          *int64
	UserID      int64
	Type        KeyType
	PasswordID  *int64
	KeyMaterial []byte
	PublicKey   []byte
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
