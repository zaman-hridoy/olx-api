package handlers

import (
	"fmt"
	"strings"
	"time"
)

type CreateListingRequest struct {
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Price		int64		`json:"Price"`
	City		string		`json:"city"`
}

type CreateListingResponse struct {
	ID			string		`json:"id"`
	Title		string		`json:"name"`
	CrteatedAt	time.Time	`json:"created_at"`
}


type ValidationError struct {
	Field string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (req CreateListingRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{
			Field: "title",
			Message: "must not be empty",
		}
	}

	if strings.TrimSpace(req.Description) == "" {
		return &ValidationError{
			Field: "description",
			Message: "must not be empty",
		}
	}

	if req.Price <= 0 {
		return &ValidationError{
			Field: "price",
			Message: "please add listing price",
		}
	}

	if strings.TrimSpace(req.City) == "" {
		return &ValidationError{
			Field: "city",
			Message: "must not be empty",
		}
	}

	return nil
}