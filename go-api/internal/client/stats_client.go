// Package client is the outbound adapter to the Node statistics service.
// It is defined behind an interface so the service layer can be unit-tested
// with a mock instead of a live HTTP call.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
)

// NamedMatrix is a matrix tagged with a name (e.g. "q", "r").
type NamedMatrix struct {
	Name string
	Data [][]float64
}

// StatsResponse mirrors the Node /internal/stats response body.
type StatsResponse struct {
	PerMatrix []dto.MatrixStats `json:"perMatrix"`
	Combined  dto.CombinedStats `json:"combined"`
}

// StatsClient computes statistics over a set of matrices via the Node service.
type StatsClient interface {
	ComputeStats(ctx context.Context, token string, matrices []NamedMatrix) (*StatsResponse, error)
}

type httpStatsClient struct {
	baseURL string
	http    *http.Client
}

// NewStatsClient returns an HTTP-backed StatsClient. The timeout bounds the
// outbound call so a slow Node service cannot hang the Go request.
func NewStatsClient(baseURL string, timeout time.Duration) StatsClient {
	return &httpStatsClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

type statsRequest struct {
	Matrices []statsMatrix `json:"matrices"`
}

type statsMatrix struct {
	Name string      `json:"name"`
	Data [][]float64 `json:"data"`
}

func (c *httpStatsClient) ComputeStats(ctx context.Context, token string, matrices []NamedMatrix) (*StatsResponse, error) {
	payload := statsRequest{Matrices: make([]statsMatrix, len(matrices))}
	for i, m := range matrices {
		payload.Matrices[i] = statsMatrix{Name: m.Name, Data: m.Data}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal stats request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/stats", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build stats request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call stats service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stats service returned status %d", resp.StatusCode)
	}

	var out StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode stats response: %w", err)
	}
	return &out, nil
}
