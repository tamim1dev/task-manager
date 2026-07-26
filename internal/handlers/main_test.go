package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	os.Setenv("JWT_SECRET", "4vGux6CHn8zcE0NBCi5KHNPNjdfsbv^89hkk=")

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

func truncateTables(t *testing.T) {
	t.Helper()
	_, err := database.DB.Pool.Exec(context.Background(),
		"TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func seedUserViaRegister(t *testing.T, name, email, password string) {
	t.Helper()

	payload := fmt.Sprintf(`{"name":%q,"email":%q,"password":%q}`, name, email, password)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	RegisterUser(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("setup: failed to seed user via register, status %d, body: %s", res.StatusCode, body)
	}
}

func getJwt(t *testing.T, email, password string) string {
	t.Helper()

	payload := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-type", "application/json")
	rec := httptest.NewRecorder()
	LoginUser(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("setup: failed to login, status %d, body: %s", res.StatusCode, body)
	}

	var tokenResp models.TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("setup: failed to decode login response: %v", err)
	}

	if tokenResp.Token == "" {
		t.Fatal("setup: login succeeded but token was empty")
	}

	return tokenResp.Token
}
