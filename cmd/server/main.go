package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sweng-task/internal/config"
	"sweng-task/internal/handler"
	"sweng-task/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	log := logger.Sugar()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Infow("Configuration loaded",
		"environment", cfg.App.Environment,
		"log_level", cfg.App.LogLevel,
		"server_port", cfg.Server.Port,
	)

	// Initialize services
	lineItemService := service.NewLineItemService(log)
	KafkaConfig := config.KafkaConfigLoad()
	pubSub := service.NewPubSub(log, cfg)
	errPubSub := pubSub.Connect(KafkaConfig)
	if errPubSub != nil {
		panic(fmt.Sprintf("Failed to connect to Kafka: %v", errPubSub))
	}

	generator := service.NewDataGenerator(log, lineItemService)
	generator.GenerateLineItems()
	runTimeDBService := service.NewRunTimeDB(log)
	advertisementService := service.NewAdService(log, runTimeDBService, lineItemService)
	dataProcessorService := service.NewDataProcessorService(log, runTimeDBService, lineItemService)
	onload := service.NewOnloadService(log, dataProcessorService)
	onload.Start()

	// Note: AdService implementation is left for the candidate

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Ad Bidding Service",
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Timeout,
	})

	// Register middleware
	//app.Use(recover.New())
	//app.Use(recover.New(recover.Config{
	//	EnableStackTrace: true,
	//}))
	app.Use(fiberlogger.New())
	app.Use(cors.New())

	// Register routes
	app.Get("/health", handler.HealthCheck)

	api := app.Group("/api/v1")

	// Line Item endpoints
	lineItemHandler := handler.NewLineItemHandler(lineItemService, log)
	api.Post("/lineitems", lineItemHandler.Create)
	api.Get("/lineitems", lineItemHandler.GetAll)
	api.Get("/lineitems/:id", lineItemHandler.GetByID)

	adHandler := handler.NewAdHandler(log, advertisementService)
	api.Get("/ads", adHandler.GetWinningAds)

	trackingHandler := handler.NewTrackingHandler(log, pubSub)
	api.Post("/tracking", trackingHandler.TrackEvent)

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Infof("Starting server on %s", address)
		if err := app.Listen(address); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Info("Server gracefully stopped")
}
