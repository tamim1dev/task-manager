package handlers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Println("failed to start postgres container:", err)
		os.Exit(1)
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Println("failed to get connection string:", err)
		os.Exit(1)
	}

	if err := runTern(connString); err != nil {
		fmt.Println("failed to migrate:", err)
		os.Exit(1)
	}

	// DB connection
	var dberr error
	database.DB.Pool, dberr = pgxpool.New(context.Background(), connString)
	if dberr != nil {
		fmt.Fprintf(os.Stderr, "Connection pool error: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	database.DB.Pool.Close()
	pgContainer.Terminate(ctx)

	os.Exit(code)
}

func runTern(connString string) error {
	cmd := exec.Command("tern", "migrate", "--conn-string", connString)
	cmd.Dir = "../../migrations"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
