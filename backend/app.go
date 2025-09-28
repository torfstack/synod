package backend

import (
	"context"
	"fmt"

	"github.com/torfstack/synod/backend/config"
	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/http"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run() error {
	cfg, err := config.ParseFile("config.yaml")
	if err != nil {
		return fmt.Errorf("could not parse config at './config.yaml': %v", err)
	}

	err = db.Migrate(context.Background(), cfg.DB.ConnectionString())
	if err != nil {
		return fmt.Errorf("could not migrate database: %v", err)
	}

	database, err := db.NewDatabase(context.Background(), cfg.DB.ConnectionString())
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	domainService := domain.NewDomainService(database)
	server := http.NewServer(*cfg, domainService)

	return server.Start()
}
