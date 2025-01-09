package pkg

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type ctxResult string

const (
	headerKey    string    = "Content-Type"
	headerVal    string    = "application/json"
	ResultCtxKey ctxResult = "parsed-token"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Since(start)

		log.Printf("Dur %v: %s - %s", end, r.URL, r.Method)
	})
}

func AuthMiddleWare(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set(headerKey, headerVal)
			WriteResponse(w, http.StatusUnauthorized, "missing 'Authorization' header")
			return
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			w.Header().Set(headerKey, headerVal)
			WriteResponse(w, http.StatusUnauthorized, "uncorrect token format")
			return
		}

		result, err := VerifyToken(token[1], secretKey)
		if err != nil {
			w.Header().Set(headerKey, headerVal)
			WriteResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), ResultCtxKey, result)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
