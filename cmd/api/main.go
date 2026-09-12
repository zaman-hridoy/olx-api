package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/zaman-hridoy/olx-api/internal/config"
	"github.com/zaman-hridoy/olx-api/internal/db"
	"github.com/zaman-hridoy/olx-api/internal/handlers"
	"github.com/zaman-hridoy/olx-api/internal/middleware"
)


func main() {

	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("DATABASE not Connected: %v", err)
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	fmt.Println("Database connected")
	fmt.Println("Starting olx server...")

	lh := handlers.NewListingHandler(db, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handlers.Health)
	mux.HandleFunc("GET /api/listings", lh.GetListings)
	mux.HandleFunc("DELETE /api/listings/{id}", lh.DeleteListing)


	handler := middleware.RequestId(mux)

	srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: handler,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server faild: %v", err)
	}


}