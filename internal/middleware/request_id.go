package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type contextKey string
const requestIdKey contextKey = "request_id"

func generateRequestID() string {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-ID")

		if requestId == "" {
			requestId = generateRequestID()
		}

		fmt.Println("request_id:", requestId)

		w.Header().Set("X-Request-ID", requestId)
		ctx := context.WithValue(r.Context(), requestIdKey, requestId)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestId(ctx context.Context) string {
	requestId, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		return ""
	}
	return requestId
}