package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidID_400 			Code = "invalid_id"
	CodeInternalError_500 		Code = "internal_error"
	CodeNotFound_404			Code = "not_found"
	CodeMalformedJson_400		Code = "malformed_json"
	CodeValidationFailed_422	Code = "validation_failed"
	CodeUnauthenticated_401		Code = "unauthenticated"
	CodeForbidden_403			Code = "forbidden"
	CodeConflict_409			Code = "conflict"
	CodeRateLimited_429			Code = "rate_limited"
)

type errorEnvelope struct {
	Error ErrorPayload	`json:"error"`
}

type ErrorPayload struct {
	Code		Code	`json:"code"`
	Message		string	`json:"message"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	json.NewEncoder(w).Encode(errorEnvelope{
		Error: ErrorPayload{
			Code: code,
			Message: message,
		},
	})
}