package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/client"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/repository"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/server"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/service"
)

const testSecret = "test-secret"

// --- fakes ---

type stubStats struct{ err error }

func (s stubStats) ComputeStats(context.Context, string, []client.NamedMatrix) (*client.StatsResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &client.StatsResponse{
		PerMatrix: []dto.MatrixStats{{Name: "q", Sum: 2}, {Name: "r", Sum: 8}},
		Combined:  dto.CombinedStats{Sum: 10},
	}, nil
}

type memRepo struct {
	items map[string]*dto.ComputationResponse
}

func newMemRepo() *memRepo { return &memRepo{items: map[string]*dto.ComputationResponse{}} }

func (r *memRepo) Save(_ context.Context, c *dto.ComputationResponse, _ string) error {
	r.items[c.ID] = c
	return nil
}
func (r *memRepo) GetByID(_ context.Context, id string) (*dto.ComputationResponse, error) {
	if c, ok := r.items[id]; ok {
		return c, nil
	}
	return nil, repository.ErrNotFound
}
func (r *memRepo) List(context.Context, int, int) ([]dto.ComputationSummary, int, error) {
	return []dto.ComputationSummary{}, 0, nil
}

type okPinger struct{}

func (okPinger) Ping(context.Context) error { return nil }

// --- helpers ---

func newApp(stats client.StatsClient) *fiber.App {
	repo := newMemRepo()
	return server.New(server.Deps{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		JWTSecret:    testSecret,
		JWTExpiry:    time.Hour,
		AuthUsername: "demo",
		AuthPassword: "demo123",
		QRService:    service.NewQRService(stats, repo),
		Repo:         repo,
		DB:           okPinger{},
	})
}

func doJSON(t *testing.T, app *fiber.App, method, path, token string, body any) (*http.Response, map[string]any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	var out map[string]any
	data, _ := io.ReadAll(resp.Body)
	if len(data) > 0 {
		_ = json.Unmarshal(data, &out)
	}
	return resp, out
}

func login(t *testing.T, app *fiber.App) string {
	t.Helper()
	resp, body := doJSON(t, app, http.MethodPost, "/auth/login", "",
		map[string]string{"username": "demo", "password": "demo123"})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want 200", resp.StatusCode)
	}
	token, _ := body["token"].(string)
	if token == "" {
		t.Fatal("login returned empty token")
	}
	return token
}

func errCode(body map[string]any) string {
	e, _ := body["error"].(map[string]any)
	if e == nil {
		return ""
	}
	code, _ := e["code"].(string)
	return code
}

// --- tests ---

func TestQR_RequiresAuth(t *testing.T) {
	app := newApp(stubStats{})
	resp, body := doJSON(t, app, http.MethodPost, "/api/v1/qr", "",
		map[string]any{"matrix": [][]float64{{1, 2}, {3, 4}}})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", resp.StatusCode)
	}
	if errCode(body) != "UNAUTHORIZED" {
		t.Errorf("code = %q, want UNAUTHORIZED", errCode(body))
	}
}

func TestQR_HappyPath(t *testing.T) {
	app := newApp(stubStats{})
	token := login(t, app)

	resp, body := doJSON(t, app, http.MethodPost, "/api/v1/qr", token,
		map[string]any{"matrix": [][]float64{{1, 2}, {3, 4}, {5, 6}}})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%v)", resp.StatusCode, body)
	}

	qr, _ := body["qr"].(map[string]any)
	if qr == nil || qr["q"] == nil || qr["r"] == nil {
		t.Errorf("response missing qr.q / qr.r: %v", body)
	}
	if body["statistics"] == nil {
		t.Errorf("response missing statistics")
	}
	if id, _ := body["id"].(string); id == "" {
		t.Errorf("response missing id")
	}
}

func TestQR_ValidationError(t *testing.T) {
	app := newApp(stubStats{})
	token := login(t, app)

	resp, body := doJSON(t, app, http.MethodPost, "/api/v1/qr", token,
		map[string]any{"matrix": [][]float64{{1, 2, 3}}}) // wide (m < n)
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", resp.StatusCode)
	}
	if errCode(body) != "VALIDATION_ERROR" {
		t.Errorf("code = %q, want VALIDATION_ERROR", errCode(body))
	}
}

func TestQR_UpstreamError(t *testing.T) {
	app := newApp(stubStats{err: errors.New("node down")})
	token := login(t, app)

	resp, body := doJSON(t, app, http.MethodPost, "/api/v1/qr", token,
		map[string]any{"matrix": [][]float64{{1, 2}, {3, 4}}})
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
	if errCode(body) != "UPSTREAM_ERROR" {
		t.Errorf("code = %q, want UPSTREAM_ERROR", errCode(body))
	}
}

func TestComputation_NotFound(t *testing.T) {
	app := newApp(stubStats{})
	token := login(t, app)

	resp, body := doJSON(t, app, http.MethodGet, "/api/v1/computations/does-not-exist", token, nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.StatusCode)
	}
	if errCode(body) != "NOT_FOUND" {
		t.Errorf("code = %q, want NOT_FOUND", errCode(body))
	}
}

func TestComputation_GetAfterCreate(t *testing.T) {
	app := newApp(stubStats{})
	token := login(t, app)

	_, created := doJSON(t, app, http.MethodPost, "/api/v1/qr", token,
		map[string]any{"matrix": [][]float64{{1, 2}, {3, 4}}})
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("no id returned from create")
	}

	resp, body := doJSON(t, app, http.MethodGet, "/api/v1/computations/"+id, token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if body["id"] != id {
		t.Errorf("id = %v, want %s", body["id"], id)
	}
}

func TestHealth(t *testing.T) {
	app := newApp(stubStats{})
	resp, body := doJSON(t, app, http.MethodGet, "/health", "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %v, want ok", body["status"])
	}
}
