package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/config"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/database"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/server"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize config and database
	config := config.GetConfig()
	db := database.NewPostgresDatabase(config)

	// Defer database close to ensure it shuts down at the end
	defer func() {
		if err := db.CloseDb(database.DbInstance.Db); err != nil {
			logrus.Error("Error closing SQL DB: ", err)
		}
		logrus.Print("Database has been closed!")
	}()

	//migrate the tables
	if err := db.AutoMigrateTables(); err != nil {
		logrus.Error("Error migrating tables: ", err)
	}

	// Start the server
	srv := server.NewEchoServer(config, db)
	go func() {
		if err := srv.Start(); err != nil {
			logrus.Fatal("Error starting server: ", err)
		}
	}()

	// Listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a shutdown signal
	<-quit
	logrus.Println("Shutdown signal received. Shutting down gracefully...")

	// Create a context with a timeout for the server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown: ", err)
	}

	logrus.Println("Server exited gracefully")
}
