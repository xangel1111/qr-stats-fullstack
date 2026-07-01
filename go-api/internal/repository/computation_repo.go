// Package repository is the persistence layer for QR computations, backed by
// PostgreSQL through pgx. It is defined behind an interface so it can be
// swapped or mocked.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
)

// ErrNotFound is returned when a computation id does not exist.
var ErrNotFound = errors.New("computation not found")

// ComputationRepository persists and retrieves QR computations (audit trail).
type ComputationRepository interface {
	Save(ctx context.Context, c *dto.ComputationResponse, username string) error
	GetByID(ctx context.Context, id string) (*dto.ComputationResponse, error)
	List(ctx context.Context, limit, offset int) (items []dto.ComputationSummary, total int, err error)
}

type pgRepository struct {
	pool *pgxpool.Pool
}

// New returns a PostgreSQL-backed ComputationRepository.
func New(pool *pgxpool.Pool) ComputationRepository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) Save(ctx context.Context, c *dto.ComputationResponse, username string) error {
	inputJSON, err := json.Marshal(c.Input.Matrix)
	if err != nil {
		return fmt.Errorf("marshal input matrix: %w", err)
	}
	qJSON, err := json.Marshal(c.QR.Q)
	if err != nil {
		return fmt.Errorf("marshal q matrix: %w", err)
	}
	rJSON, err := json.Marshal(c.QR.R)
	if err != nil {
		return fmt.Errorf("marshal r matrix: %w", err)
	}
	statsJSON, err := json.Marshal(c.Statistics)
	if err != nil {
		return fmt.Errorf("marshal statistics: %w", err)
	}

	const query = `
		INSERT INTO computations
			(id, username, input_matrix, rows, cols, q_matrix, r_matrix, statistics, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = r.pool.Exec(ctx, query,
		c.ID, username, inputJSON, c.Input.Rows, c.Input.Cols, qJSON, rJSON, statsJSON, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert computation: %w", err)
	}
	return nil
}

func (r *pgRepository) GetByID(ctx context.Context, id string) (*dto.ComputationResponse, error) {
	const query = `
		SELECT id, input_matrix, rows, cols, q_matrix, r_matrix, statistics, created_at
		FROM computations
		WHERE id = $1`

	var (
		res                                dto.ComputationResponse
		inputJSON, qJSON, rJSON, statsJSON []byte
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&res.ID, &inputJSON, &res.Input.Rows, &res.Input.Cols, &qJSON, &rJSON, &statsJSON, &res.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query computation: %w", err)
	}

	if err := json.Unmarshal(inputJSON, &res.Input.Matrix); err != nil {
		return nil, fmt.Errorf("unmarshal input matrix: %w", err)
	}
	if err := json.Unmarshal(qJSON, &res.QR.Q); err != nil {
		return nil, fmt.Errorf("unmarshal q matrix: %w", err)
	}
	if err := json.Unmarshal(rJSON, &res.QR.R); err != nil {
		return nil, fmt.Errorf("unmarshal r matrix: %w", err)
	}
	if err := json.Unmarshal(statsJSON, &res.Statistics); err != nil {
		return nil, fmt.Errorf("unmarshal statistics: %w", err)
	}
	return &res, nil
}

func (r *pgRepository) List(ctx context.Context, limit, offset int) ([]dto.ComputationSummary, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM computations`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count computations: %w", err)
	}

	const query = `
		SELECT id, rows, cols, username, created_at
		FROM computations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list computations: %w", err)
	}
	defer rows.Close()

	items := make([]dto.ComputationSummary, 0, limit)
	for rows.Next() {
		var s dto.ComputationSummary
		if err := rows.Scan(&s.ID, &s.Rows, &s.Cols, &s.Username, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan computation summary: %w", err)
		}
		items = append(items, s)
	}
	return items, total, rows.Err()
}
