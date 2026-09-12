package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/zaman-hridoy/olx-api/internal/httpx"
	"github.com/zaman-hridoy/olx-api/internal/middleware"
)

type listing struct {
	ID			string		`json:"id"`
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Price		int64		`json:"Price"`
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
	request_id := middleware.GetRequestId(ctx)
	id := r.PathValue("id")

	// lh.logger.Debug("debug log", "listing_id", id)
	// lh.logger.Info("Starting query", "listing_id", id, "request_id", request_id)
	// lh.logger.Warn("warn log", "listing_id", id)

	results, err := lh.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`,id)
	if err != nil {
		// log.Fatalf("[ERROR - DELETE LISTING]: %v", err)
		
		lh.logger.Error("delete failed", "listing_id", id, "err", err, "request_id", request_id)
		// http.Error(w, "internal error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "internal server error", httpx.CodeInternalError_500)
		return
	}

	count, err := results.RowsAffected()
	if err != nil {
		log.Fatalf("[ERROR - DELETE LISTING ROWS AFFTECTED]: %v", err, )
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}


	lh.logger.Info("Record deleted", "listing_id", id, "request_id", request_id)
	// w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Listing id: %s, count: %d\n", id, count)
}

func (lh ListingHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.GetRequestId(ctx)
	var req CreateListingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJson_400)
		return
	}


	row := lh.db.QueryRowContext(ctx, `
		INSERT INTO listings (title, description, price, city)
		VALUES($1, $2, $3, $4)
		RETURNING id, title, created_at
	`, req.Title, req.Description, req.Price, req.City)

	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.Title, &out.CrteatedAt); err != nil {
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "someting went wrong", httpx.CodeInternalError_500)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", out.ID)
	fmt.Printf("%+v\n", req)

	httpx.JSONP(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": out,
	})
}