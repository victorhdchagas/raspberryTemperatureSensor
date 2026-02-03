package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/api"
	"github.com/wutachi/raspberryTemperatureSensor/internal/app"
	"github.com/wutachi/raspberryTemperatureSensor/internal/config"
	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
	"github.com/wutachi/raspberryTemperatureSensor/internal/maintenance"
	"github.com/wutachi/raspberryTemperatureSensor/internal/sensor"
	"github.com/wutachi/raspberryTemperatureSensor/pkg/web"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	database, err := db.New(cfg.DB.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	sensorReader, err := sensor.NewDHT11(cfg.Sensor.GPIO)
	if err != nil {
		log.Fatalf("Failed to initialize sensor: %v", err)
	}
	go sensorReader.Start(context.Background(), cfg.Sensor.Interval, database)

	maintenanceWorker := maintenance.NewWorker(database)
	go maintenanceWorker.Start(context.Background(), 24*time.Hour)

	appConfig := app.NewConfig()
	apiHandler := api.NewHandler(database, appConfig)

	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)
	web.RegisterStaticRoutes(mux)

	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
