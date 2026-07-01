// Package service contains the QR use case: it orchestrates validation, QR
// factorization, the call to the Node statistics service and persistence.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/client"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/matrix"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/repository"
)

// QRService orchestrates the full QR computation pipeline.
type QRService struct {
	stats client.StatsClient
	repo  repository.ComputationRepository

	// clock and idgen are injected so tests can make results deterministic.
	clock func() time.Time
	idgen func() string
}

// NewQRService wires the service with its collaborators and production
// clock/id generators.
func NewQRService(stats client.StatsClient, repo repository.ComputationRepository) *QRService {
	return &QRService{
		stats: stats,
		repo:  repo,
		clock: func() time.Time { return time.Now().UTC() },
		idgen: uuid.NewString,
	}
}

// Compute validates the input, factorizes it, fetches statistics from Node,
// persists the result and returns the full response.
func (s *QRService) Compute(ctx context.Context, input [][]float64, username, token string) (*dto.ComputationResponse, error) {
	if err := matrix.Validate(input); err != nil {
		return nil, apperr.Validation(err.Error())
	}

	factors := matrix.Factorize(input)

	statsResp, err := s.stats.ComputeStats(ctx, token, []client.NamedMatrix{
		{Name: "q", Data: factors.Q},
		{Name: "r", Data: factors.R},
	})
	if err != nil {
		return nil, apperr.Upstream("statistics service unavailable")
	}

	statistics, err := mapStatistics(statsResp)
	if err != nil {
		return nil, apperr.Upstream("unexpected statistics response")
	}

	res := &dto.ComputationResponse{
		ID: s.idgen(),
		Input: dto.InputInfo{
			Matrix: input,
			Rows:   len(input),
			Cols:   len(input[0]),
		},
		QR:         dto.QRResult{Q: factors.Q, R: factors.R},
		Statistics: statistics,
		CreatedAt:  s.clock(),
	}

	// Persistence is synchronous and required: an audit log that silently
	// drops records is worse than a failed request in this domain.
	if err := s.repo.Save(ctx, res, username); err != nil {
		return nil, apperr.Internal("could not persist computation")
	}

	return res, nil
}

// mapStatistics maps the Node per-matrix list (keyed by name) into the shaped
// Statistics response.
func mapStatistics(r *client.StatsResponse) (dto.Statistics, error) {
	byName := make(map[string]dto.MatrixStats, len(r.PerMatrix))
	for _, m := range r.PerMatrix {
		byName[m.Name] = m
	}

	q, okQ := byName["q"]
	rr, okR := byName["r"]
	if !okQ || !okR {
		return dto.Statistics{}, fmt.Errorf("stats response missing q or r")
	}

	return dto.Statistics{Q: q, R: rr, Combined: r.Combined}, nil
}
