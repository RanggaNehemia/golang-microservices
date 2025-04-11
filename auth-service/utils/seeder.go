package utils

import (
	"context"
	"log"

	oauthModels "github.com/go-oauth2/oauth2/v4/models"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

func SeedOAuthClients(ctx context.Context, conn *pgx.Conn) {
	adapter := pgx4adapter.NewConn(conn)
	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		log.Fatalf("Failed to create client store for seeding: %v", err)
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
			log.Printf("OAuth2 client '%s' already exists, skipping.", c.ID)
			continue
		}

		err = clientStore.Create(c)
		if err != nil {
			log.Fatalf("Failed to seed client '%s': %v", c.ID, err)
		}

		log.Printf("Seeded client '%s'", c.ID)
	}
}
