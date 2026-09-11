package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zaman-hridoy/olx-api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down | force>")
	}

	cfg := config.MustLoad()


	 m, err := migrate.New(
        "file://migrations",
        cfg.DatabaseUrl)

	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("m.Up: %v\n", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("m.Down: %v\n", err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}

		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("m.Force: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s\n", os.Args[1])
	}
	fmt.Println("Running Migration...")
}