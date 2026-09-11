package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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


type ListingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (lh ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//pg_sleep(20)
	rows, err := lh.db.QueryContext(ctx,
	`
		SELECT id, title, description, price, city, created_at, 
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



func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	results, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`,id)
	if err != nil {
		log.Fatalf("[ERROR - DELETE LISTING]: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	count, err := results.RowsAffected()
	if err != nil {
		log.Fatalf("[ERROR - DELETE LISTING ROWS AFFTECTED]: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}


	// w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Listing id: %s, count: %d\n", id, count)
}