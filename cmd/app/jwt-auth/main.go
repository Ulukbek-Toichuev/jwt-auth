package main

import (
	"context"
	"errors"
	"jwt-auth/config"
	"jwt-auth/internal/api"
	"jwt-auth/internal/db"
	"jwt-auth/internal/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
)

func main() {
	config := config.NewConfig()

	db := db.NewDB(config.GetDriverName(), config.GetDataSource(), config.GetMigrationsPath())

	validate := validator.New()
	validate.RegisterValidation("check_status", util.ValidateTodoStatus)

	authHandler := api.NewAuthHandler(db, config, validate)
	userHandler := api.NewUserHandler(db, validate)
	todoHandler := api.NewTodoHandler(db, validate)

	mux := http.NewServeMux()

	route := api.NewRoute(config.GetSecretKey())

	route.SetAuthRoute(mux, authHandler)
	route.SetUserRoute(mux, userHandler)
	route.SetTodoRoute(mux, todoHandler)

	completedHandler := route.CompletedHandler(mux)

	server := &http.Server{
		Addr:    config.GetPort(),
		Handler: completedHandler,
	}

	run(server, config.GetPort())
	shutDown(server, func() { db.Close() })
}

func run(server *http.Server, port string) {
	go func() {
		log.Printf("Starting server on port - %s", port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Printf("Stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan
}

func shutDown(server *http.Server, funcToClose ...func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	for _, f := range funcToClose {
		f()
	}
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Printf("Server graceful shutdown!")
}
