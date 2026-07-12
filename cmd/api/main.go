package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/handlers"
	"github.com/tamim1dev/task-manager/internal/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	var dberr error
	database.DB.Pool, dberr = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if dberr != nil {
		fmt.Fprintf(os.Stderr, "Connection pool error: %v\n", err)
		os.Exit(1)
	}
	defer database.DB.Pool.Close()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	// register new user
	router.Post("/register", handlers.RegisterUser)
	// login user
	router.Post("/login", handlers.LoginUser)
	// jwt middleware check
	router.Get("/me", middleware.AuthMiddleware(handlers.GetMe))

	serverStartError := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if serverStartError != nil {
		fmt.Println("Failed to start server")
		os.Exit(1)
	}
}
