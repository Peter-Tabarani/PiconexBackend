package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Peter-Tabarani/PiconexBackend/internal"
	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
)

func main() {
	dsn := "piconex:pjaplmTabs7!@tcp(178.156.189.138:3306)/piconexdb"
	db, err := utils.Connect(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	router := internal.NewRouter(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		log.Println("✅ Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")
}
