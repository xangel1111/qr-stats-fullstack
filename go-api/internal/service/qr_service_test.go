package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/client"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
)

// --- fakes for the service's collaborators ---

type fakeStats struct {
	resp   *client.StatsResponse
	err    error
	called bool
}

func (f *fakeStats) ComputeStats(context.Context, string, []client.NamedMatrix) (*client.StatsResponse, error) {
	f.called = true
	return f.resp, f.err
}

type fakeRepo struct {
	saveErr   error
	saved     *dto.ComputationResponse
	savedUser string
}

func (r *fakeRepo) Save(_ context.Context, c *dto.ComputationResponse, username string) error {
	r.saved = c
	r.savedUser = username
	return r.saveErr
}
func (r *fakeRepo) GetByID(context.Context, string) (*dto.ComputationResponse, error) { return nil, nil }
func (r *fakeRepo) List(context.Context, int, int) ([]dto.ComputationSummary, int, error) {
	return nil, 0, nil
}

func okStats() *client.StatsResponse {
	return &client.StatsResponse{
		PerMatrix: []dto.MatrixStats{
			{Name: "q", Sum: 2},
			{Name: "r", Sum: 12},
		},
		Combined: dto.CombinedStats{Sum: 14},
	}
}

// newService injects deterministic clock/id so results are assertable.
func newService(stats client.StatsClient, repo *fakeRepo) *QRService {
	s := NewQRService(stats, repo)
	s.clock = func() time.Time { return time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC) }
	s.idgen = func() string { return "fixed-id" }
	return s
}

func assertStatus(t *testing.T, err error, want int) {
	t.Helper()
	var apiErr *apperr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is not *apperr.APIError: %v", err)
	}
	if apiErr.Status != want {
		t.Fatalf("status = %d, want %d", apiErr.Status, want)
	}
}

func TestCompute_HappyPath(t *testing.T) {
	stats := &fakeStats{resp: okStats()}
	repo := &fakeRepo{}
	s := newService(stats, repo)

	res, err := s.Compute(context.Background(), [][]float64{{1, 2}, {3, 4}, {5, 6}}, "demo", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != "fixed-id" {
		t.Errorf("ID = %q, want fixed-id", res.ID)
	}
	if res.Input.Rows != 3 || res.Input.Cols != 2 {
		t.Errorf("dims = %dx%d, want 3x2", res.Input.Rows, res.Input.Cols)
	}
	if len(res.QR.Q) != 3 || len(res.QR.R) != 2 {
		t.Errorf("unexpected QR dims: Q=%d rows, R=%d rows", len(res.QR.Q), len(res.QR.R))
	}
	if res.Statistics.Q.Sum != 2 || res.Statistics.R.Sum != 12 {
		t.Errorf("statistics not mapped: Q.Sum=%v R.Sum=%v", res.Statistics.Q.Sum, res.Statistics.R.Sum)
	}
	if !stats.called {
		t.Error("stats client was not called")
	}
	if repo.saved == nil || repo.savedUser != "demo" {
		t.Errorf("computation not persisted with username; saved=%v user=%q", repo.saved, repo.savedUser)
	}
}

func TestCompute_ValidationError(t *testing.T) {
	s := newService(&fakeStats{resp: okStats()}, &fakeRepo{})
	// Wide matrix (m < n) is rejected before any downstream call.
	_, err := s.Compute(context.Background(), [][]float64{{1, 2, 3}}, "demo", "tok")
	assertStatus(t, err, 422)
}

func TestCompute_UpstreamError(t *testing.T) {
	s := newService(&fakeStats{err: errors.New("node down")}, &fakeRepo{})
	_, err := s.Compute(context.Background(), [][]float64{{1, 2}, {3, 4}}, "demo", "tok")
	assertStatus(t, err, 502)
}

func TestCompute_MissingStats(t *testing.T) {
	// Stats response missing the "r" matrix -> mapping fails -> upstream error.
	stats := &fakeStats{resp: &client.StatsResponse{PerMatrix: []dto.MatrixStats{{Name: "q"}}}}
	s := newService(stats, &fakeRepo{})
	_, err := s.Compute(context.Background(), [][]float64{{1, 2}, {3, 4}}, "demo", "tok")
	assertStatus(t, err, 502)
}

func TestCompute_PersistError(t *testing.T) {
	repo := &fakeRepo{saveErr: errors.New("db down")}
	s := newService(&fakeStats{resp: okStats()}, repo)
	_, err := s.Compute(context.Background(), [][]float64{{1, 2}, {3, 4}}, "demo", "tok")
	assertStatus(t, err, 500)
}
