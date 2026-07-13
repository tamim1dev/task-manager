package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/routers"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func main() {
	// Environment vars
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// DB connection
	var dberr error
	database.DB.Pool, dberr = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if dberr != nil {
		fmt.Fprintf(os.Stderr, "Connection pool error: %v\n", err)
		os.Exit(1)
	}
	defer database.DB.Pool.Close()

	// chi initialization
	router := chi.NewRouter()
	// healthcheck
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	// router mounts
	router.Mount("/auth", routers.AuthRouter())

	// start server and gracefull shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}
	utils.StartServerAndGracefullyShutdown(srv, 5*time.Second)
}
