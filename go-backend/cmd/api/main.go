package main

import (
	"fmt"
	"log"

	"github.com/nextolympservice/go-backend/config"
	"github.com/nextolympservice/go-backend/internal/database"
	"github.com/nextolympservice/go-backend/internal/router"
	"github.com/nextolympservice/go-backend/internal/utils"
)

func main() {
	// Config yuklash
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Custom validatorlarni ro'yxatga olish
	utils.SetupValidator()

	// Database ulanish
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto migration
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Router sozlash
	r := router.Setup(cfg, db)

	// Server ishga tushirish
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("Server starting on %s (env: %s)", addr, cfg.App.Env)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
