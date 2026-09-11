package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zaman-hridoy/olx-api/internal/config"
	"github.com/zaman-hridoy/olx-api/internal/db"
	"github.com/zaman-hridoy/olx-api/internal/handlers"
)


func main() {

	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("DATABASE not Connected: %v", err)
	}

	fmt.Println("Database connected")
	fmt.Println("Starting olx server...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.Listings(db))


	srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server faild: %v", err)
	}


}