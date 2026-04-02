package main

import (
	"context"
	"log"

	"github.com/superbkibbles/ecommerce/internal/adapters/repository/mongodb"
	"github.com/superbkibbles/ecommerce/internal/config"
	"github.com/superbkibbles/ecommerce/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := mongodb.NewMongoDBConnection(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}

	log.Println("Repairing category BSON keys and category references on products / home_sections…")
	if err := utils.RepairCategoryBSONKeysAndRefs(context.Background(), db); err != nil {
		log.Fatalf("repair failed: %v", err)
	}
	log.Println("Repair finished successfully.")
}
