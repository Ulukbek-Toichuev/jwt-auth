package pkg

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type ctxResult string

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":   "*",
	"Access-Control-Allow-Headers":  "Accept, Authorization, Content-Type, Origin",
	"Access-Control-Allow-Methods":  "GET, POST, PUT, DELETE, OPTIONS",
	"Access-Control-Expose-Headers": "Content-Length, Content-Type",
}

const ResultCtxKey ctxResult = "parsed-token"

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
			WriteResponse(w, http.StatusUnauthorized, "missing 'Authorization' header")
			return
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			WriteResponse(w, http.StatusUnauthorized, "uncorrect token format")
			return
		}

		result, err := VerifyToken(token[1], secretKey)
		if err != nil {
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
		for k, v := range corsHeaders {
			w.Header().Set(k, v)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
