package utils

import (
	"context"

	oauthModels "github.com/go-oauth2/oauth2/v4/models"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
	"go.uber.org/zap"
)

func SeedOAuthClients(ctx context.Context, conn *pgx.Conn) {
	adapter := pgx4adapter.NewConn(conn)
	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		Logger.Fatal("Failed to create client store for seeding", zap.Error(err))
	}

	clients := []*oauthModels.Client{
		{
			ID:     "data-service",
			Secret: "dataservicesecret",
			Domain: "",
			UserID: "",
		},
		{
			ID:     "trade-service",
			Secret: "tradeservicesecret",
			Domain: "",
			UserID: "",
		},
		{
			ID:     "webclient",
			Secret: "webclientsecret",
			Domain: "",
			UserID: "",
		},
	}

	for _, c := range clients {
		_, err := clientStore.GetByID(ctx, c.ID)
		if err == nil {
			Logger.Warn("OAuth2 client already exists, skipping.", zap.String("Client", c.ID))
			continue
		}

		err = clientStore.Create(c)
		if err != nil {
			Logger.Fatal("Failed to seed client", zap.Error(err))
		}

		Logger.Info("Seeded client", zap.String("Client", c.ID))
	}
}
