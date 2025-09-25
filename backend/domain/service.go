package domain

import (
	"github.com/torfstack/synod/backend/db"
)

type service struct {
	database db.Database
	sessions sessionStore
}

var _ Service = (*service)(nil)

func NewDomainService(db db.Database) Service {
	return &service{
		database: db,
		sessions: make(sessionStore),
	}
}
