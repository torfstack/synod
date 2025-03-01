package backend

import (
	"context"
	"fmt"
	"github.com/torfstack/kayvault/backend/config"
	"github.com/torfstack/kayvault/backend/db"
	"github.com/torfstack/kayvault/backend/domain"
	"github.com/torfstack/kayvault/backend/http"
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

	database := db.NewDatabase(cfg.DB.ConnectionString())
	domainService := domain.NewDomainService(database)
	server := http.NewServer(*cfg, domainService)

	return server.Start()
}
