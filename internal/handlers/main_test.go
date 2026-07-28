package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/handlers"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/routers"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	if err := os.Setenv("JWT_SECRET", "4vGux6CHn8zcE0NBCi5KHNPNjdfsbv^89hkk="); err != nil {
		fmt.Println("failed to set JWT_SECRET:", err)
		os.Exit(1)
	}

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
	if err := pgContainer.Terminate(ctx); err != nil {
		fmt.Println("failed to terminate container:", err)
	}

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
	handlers.RegisterUser(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}()

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
	handlers.LoginUser(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}()

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

func newAuthedRequest(method, path, body, token string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func newTaskRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/tasks", routers.TasksRouter())
	return r
}

func setupTaskForUser(t *testing.T, token string) string {
	t.Helper()

	router := newTaskRouter()

	payload := `{"title":"Buy groceries","description":"Milk, eggs, bread","days":7}`
	req := newAuthedRequest(http.MethodPost, "/tasks", payload, token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("setup: failed to create task, status %d, body: %s", res.StatusCode, body)
	}

	var task models.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("setup: failed to decode created task: %v", err)
	}

	if task.Id.String() == "" {
		t.Fatal("setup: created task but got empty ID")
	}

	return task.Id.String()
}

type testUser struct {
	Name     string
	Email    string
	Password string
	Token    string
}

func setupAuthedUser(t *testing.T) testUser {
	t.Helper()

	user := testUser{
		Name:     "abc",
		Email:    "abc@gmail.com",
		Password: "abcpass",
	}

	seedUserViaRegister(t, user.Name, user.Email, user.Password)
	user.Token = getJwt(t, user.Email, user.Password)

	return user
}
