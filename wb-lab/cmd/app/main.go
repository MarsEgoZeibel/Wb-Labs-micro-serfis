package main

import (
	"database/sql"
	"fmt"
	"log"
	"wb-lab/internal/cache"
	"wb-lab/internal/config"
	"wb-lab/internal/handlers"
	"wb-lab/internal/nats"
	"wb-lab/internal/repository"
	"wb-lab/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	connStr := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	db, err := sql.Open("postgres",
		fmt.Sprintf(connStr,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
		),
	)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	orderCache := cache.NewOrderCache()
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, orderCache)
	orderHandler := handlers.NewOrderHandler(orderService)

	err = orderService.RestoreCache()
	if err != nil {
		log.Printf("Error restoring cache: %v", err)
	}

	natsService, err := nats.NewNATSService(
		cfg.NATS.ClusterID,
		cfg.NATS.ClientID,
		orderService,
	)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsService.Close()

	err = natsService.Subscribe(cfg.NATS.Subject)
	if err != nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}

	router := gin.Default()
	handlers.SetupRoutes(router, orderHandler)

	log.Printf("Starting server on port %s", cfg.HTTP.Port)
	if err := router.Run(":" + cfg.HTTP.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
