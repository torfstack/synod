package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/torfstack/synod/backend/config"
	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/http"
	"github.com/torfstack/synod/backend/maintenance"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run(ctx context.Context) error {
	cfg, err := config.ParseFile("config.yaml")
	if err != nil {
		return fmt.Errorf("could not parse config at './config.yaml': %v", err)
	}

	err = db.Migrate(ctx, cfg.DB.ConnectionString())
	if err != nil {
		return fmt.Errorf("could not migrate database: %v", err)
	}

	database, err := db.NewDatabase(ctx, cfg.DB.ConnectionString())
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	domainService := domain.NewDomainService(ctx, database)
	go maintenance.Run(ctx, maintenance.Job{
		Name:     "expired threshold unlocks",
		Interval: time.Hour,
		Run: func(ctx context.Context) error {
			_, err := database.DeleteExpiredUnlockRequests(ctx)
			return err
		},
	})
	server := http.NewServer(*cfg, domainService)

	return server.Start(ctx)
}
