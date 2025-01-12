package main

import (
	"context"
	"errors"
	"jwt-auth/config"
	"jwt-auth/internal/api"
	"jwt-auth/internal/db"
	"jwt-auth/pkg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := config.NewConfig()
	db := db.NewDB(config.GetDriverName(), config.GetDataSource())

	authHandler := api.NewAuthHandler(db, config)
	todoHandler := api.NewTodoHandler(db)
	userHandler := api.NewUserHandler(db)
	mux := http.NewServeMux()

	{
		mux.HandleFunc("POST /api/sign-in", authHandler.SignIn)
		mux.HandleFunc("POST /api/sign-up", authHandler.SignUp)

	}

	{
		mux.Handle("GET /api/admin/users", pkg.AuthMiddleWare(http.HandlerFunc(userHandler.GetAllUsers), config.Jwt.SecretKey))
		mux.Handle("GET /api/admin/users/{email}", pkg.AuthMiddleWare(http.HandlerFunc(userHandler.GetUserByEmail), config.Jwt.SecretKey))
		mux.Handle("PUT /api/admin/users", pkg.AuthMiddleWare(http.HandlerFunc(userHandler.ChangeUsersRole), config.Jwt.SecretKey))
		mux.Handle("DELETE /api/admin/users", pkg.AuthMiddleWare(http.HandlerFunc(userHandler.DeleteUserByEmail), config.Jwt.SecretKey))
	}

	{
		mux.Handle("GET /api/users/{id}", pkg.AuthMiddleWare(http.HandlerFunc(userHandler.GetUsersOwnDetail), config.Jwt.SecretKey))
		mux.Handle("GET /api/users/todos", pkg.AuthMiddleWare(http.HandlerFunc(todoHandler.GetAllByUser), config.Jwt.SecretKey))
		mux.Handle("GET /api/users/todos/{id}", pkg.AuthMiddleWare(http.HandlerFunc(todoHandler.GetTodoById), config.Jwt.SecretKey))
		mux.Handle("POST /api/users/todos", pkg.AuthMiddleWare(http.HandlerFunc(todoHandler.CreateTodo), config.Jwt.SecretKey))
		mux.Handle("PUT /api/users/todos/{id}", pkg.AuthMiddleWare(http.HandlerFunc(todoHandler.UpdateTodoStatus), config.Jwt.SecretKey))
		mux.Handle("DELETE /api/users/todos/{id}", pkg.AuthMiddleWare(http.HandlerFunc(todoHandler.DeleteTodoById), config.Jwt.SecretKey))
	}

	completedHandler := pkg.CORSMiddleware(
		pkg.LoggerMiddleware(mux),
	)

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
		if err != nil && errors.Is(err, http.ErrServerClosed) {
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
