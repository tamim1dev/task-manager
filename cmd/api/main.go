package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tamim1dev/task-manager/docs"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/middleware"
	"github.com/tamim1dev/task-manager/internal/routers"
	"github.com/tamim1dev/task-manager/internal/utils"
)

// @title Task Manager API
// @version 1.0
// @description A simple task manager backend with JWT auth
// @host localhost:5000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	router.Use(chimw.RequestID)
	router.Use(middleware.RequestLogger)

	// healthcheck
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello")); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	})
	// router mounts
	router.Mount("/auth", routers.AuthRouter())
	router.Mount("/users", routers.UsersRouter())
	router.Mount("/tasks", routers.TasksRouter())
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// start server and gracefull shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}
	utils.StartServerAndGracefullyShutdown(srv, 5*time.Second)
}
