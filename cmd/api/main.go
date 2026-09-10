package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/zaman-hridoy/olx-api/internal/config"
)


func main() {

	cfg := config.MustLoad()

	fmt.Println("Starting olx server...", cfg.Env)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "All OK"}`))
	})


	srv := http.Server{
		Addr: ":" + os.Getenv("PORT"),
		Handler: mux,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server faild: %v", err)
	}


}