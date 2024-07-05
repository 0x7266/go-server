package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   []User
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewAPIServer(addr string) *APIServer {
	users := []User{}
	return &APIServer{
		addr: addr,
		db:   users,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: " + userID))
	})
	router.HandleFunc(
		"POST /users",
		func(w http.ResponseWriter, r *http.Request) {
			var user User
			json.NewDecoder(r.Body).Decode(&user)
			s.db = append(s.db, user)
			fmt.Println(s.db)
		},
	)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1"), router)

	middlewareChain := MiddlewareChain(RequestLoggerMiddleware, RequireAuthMiddleware)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(router),
	}

	log.Printf("Server running on %s", s.addr)
	return server.ListenAndServe()
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("METHOD: %s\nPATH: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}
