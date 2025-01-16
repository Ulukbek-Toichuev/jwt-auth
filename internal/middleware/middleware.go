package middleware

import (
	"context"
	"jwt-auth/internal/util"
	"log"
	"net/http"
	"strings"
	"time"
)

type ctxResult int

const ResultCtxKey ctxResult = 469740

type MiddleWare struct {
	secretKey   string
	corsHeaders map[string]string
}

func NewMiddleWare(secretKey string) *MiddleWare {
	middleware := MiddleWare{
		secretKey,
		map[string]string{
			"Access-Control-Allow-Origin":   "*",
			"Access-Control-Allow-Headers":  "Accept, Authorization, Content-Type, Origin",
			"Access-Control-Allow-Methods":  "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Expose-Headers": "Content-Length, Content-Type",
		},
	}

	return &middleware
}

func (mw *MiddleWare) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Since(start)

		log.Printf("Dur %v: %s - %s", end, r.URL, r.Method)
	})
}

func (mw *MiddleWare) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.WriteResponseWithMssg(w, http.StatusUnauthorized, "missing 'Authorization' header")
			return
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			util.WriteResponseWithMssg(w, http.StatusUnauthorized, "uncorrect token format")
			return
		}

		result, err := util.VerifyToken(token[1], mw.secretKey)
		if err != nil {
			util.WriteResponseWithMssg(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), ResultCtxKey, result)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (mw *MiddleWare) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range mw.corsHeaders {
			w.Header().Set(k, v)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
