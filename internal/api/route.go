package api

import (
	"jwt-auth/internal/middleware"
	"net/http"
)

type Route struct {
	mdlw *middleware.MiddleWare
}

func NewRoute(secretKey string) *Route {
	middleware := middleware.NewMiddleWare(secretKey)
	return &Route{middleware}
}

func (r *Route) SetAuthRoute(mux *http.ServeMux, authHandler *AuthHandler) {
	mux.HandleFunc("POST /api/sign-in", authHandler.SignIn)
	mux.HandleFunc("POST /api/sign-up", authHandler.SignUp)
}

func (r *Route) SetUserRoute(mux *http.ServeMux, userHandler *UserHandler) {
	mux.Handle("GET /api/admin/users", r.mdlw.AuthMiddleWare(http.HandlerFunc(userHandler.GetAllUsers)))
	mux.Handle("GET /api/admin/users/{email}", r.mdlw.AuthMiddleWare(http.HandlerFunc(userHandler.GetUserByEmail)))
	mux.Handle("PUT /api/admin/users", r.mdlw.AuthMiddleWare(http.HandlerFunc(userHandler.ChangeUsersRole)))
	mux.Handle("DELETE /api/admin/users", r.mdlw.AuthMiddleWare(http.HandlerFunc(userHandler.DeleteUserByEmail)))
}

func (r *Route) SetTodoRoute(mux *http.ServeMux, todoHandler *TodoHandler) {
	mux.Handle("GET /api/users/todos", r.mdlw.AuthMiddleWare(http.HandlerFunc(todoHandler.GetAllByUser)))
	mux.Handle("GET /api/users/todos/{id}", r.mdlw.AuthMiddleWare(http.HandlerFunc(todoHandler.GetTodoByIdAndByUser)))
	mux.Handle("POST /api/users/todos", r.mdlw.AuthMiddleWare(http.HandlerFunc(todoHandler.CreateTodo)))
	mux.Handle("PUT /api/users/todos/{id}", r.mdlw.AuthMiddleWare(http.HandlerFunc(todoHandler.UpdateTodoStatus)))
	mux.Handle("DELETE /api/users/todos/{id}", r.mdlw.AuthMiddleWare(http.HandlerFunc(todoHandler.DeleteTodoById)))
}

func (r *Route) CompletedHandler(mux *http.ServeMux) http.Handler {
	return r.mdlw.CORSMiddleware(
		r.mdlw.LoggerMiddleware(mux),
	)
}
