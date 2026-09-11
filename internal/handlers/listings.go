package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	ID			string		`json:"id"`
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Price		string		`json:"Price"`
	City		string		`json:"city"`
	CrteatedAt	time.Time	`json:"created_at"`
}

func Listings(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get listing from db
		rows, err := db.Query(
		`
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100
		`)

		if err != nil {
			log.Printf("query: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		listings := []listing{}
		for rows.Next() {
			var l listing
			if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CrteatedAt); err != nil {
				log.Printf("listing rows.scan: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			listings = append(listings, l)
		}

		if err := rows.Err(); err != nil {
			log.Printf("rows.err: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listings)
	}
}