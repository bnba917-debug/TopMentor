package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/topmentor/backend/internal/config"
	"github.com/topmentor/backend/internal/migrate"
	"github.com/topmentor/backend/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pg, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pg.Close()

	migrationsDir := filepath.Join("migrations")
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	if err := migrate.Run(context.Background(), pg.DB(), migrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	fmt.Println("migrations complete")
}
