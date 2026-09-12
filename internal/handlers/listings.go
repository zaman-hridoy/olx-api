package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
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
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db: db,
		logger: logger,
	}
}

func (lh ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//pg_sleep(20)
	rows, err := lh.db.QueryContext(ctx,
	`
		SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
	`)

	if err != nil {
		lh.logger.Error("listing query error", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CrteatedAt); err != nil {
			lh.logger.Error("listing row scan", "error", err)
			// log.Printf("listing rows.scan: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("rows err", "error", err)
		// log.Printf("rows.err: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	lh.logger.Info("listings fetched", "total", len(listings))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listings)
}



func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	lh.logger.Debug("debug log", "listing_id", id)
	lh.logger.Info("Starting query", "listing_id", id)
	lh.logger.Warn("warn log", "listing_id", id)

	results, err := lh.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`,id)
	if err != nil {
		// log.Fatalf("[ERROR - DELETE LISTING]: %v", err)
		
		lh.logger.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	count, err := results.RowsAffected()
	if err != nil {
		log.Fatalf("[ERROR - DELETE LISTING ROWS AFFTECTED]: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}


	lh.logger.Info("Record deleted", "listing_id", id)
	// w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Listing id: %s, count: %d\n", id, count)
}