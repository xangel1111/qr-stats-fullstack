// Package dto holds the request/response payloads that define the public API
// contract. Keeping them in one place makes the contract easy to review.
package dto

import "time"

// LoginRequest is the body of POST /auth/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is returned by POST /auth/login.
type LoginResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

// QRRequest is the body of POST /api/v1/qr.
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// InputInfo describes the original matrix.
type InputInfo struct {
	Matrix [][]float64 `json:"matrix"`
	Rows   int         `json:"rows"`
	Cols   int         `json:"cols"`
}

// QRResult holds the Q and R factor matrices.
type QRResult struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// MatrixStats are the statistics of a single matrix, as computed by the Node
// service. The JSON tags mirror the Node response so it decodes directly.
type MatrixStats struct {
	Name       string  `json:"name,omitempty"`
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	Average    float64 `json:"average"`
	Sum        float64 `json:"sum"`
	IsDiagonal bool    `json:"isDiagonal"`
}

// CombinedStats are statistics aggregated across all matrices.
type CombinedStats struct {
	Max         float64 `json:"max"`
	Min         float64 `json:"min"`
	Average     float64 `json:"average"`
	Sum         float64 `json:"sum"`
	AnyDiagonal bool    `json:"anyDiagonal"`
}

// Statistics is the statistics section exposed by the QR endpoint.
type Statistics struct {
	Q        MatrixStats   `json:"q"`
	R        MatrixStats   `json:"r"`
	Combined CombinedStats `json:"combined"`
}

// ComputationResponse is returned by POST /api/v1/qr and GET .../{id}.
type ComputationResponse struct {
	ID         string     `json:"id"`
	Input      InputInfo  `json:"input"`
	QR         QRResult   `json:"qr"`
	Statistics Statistics `json:"statistics"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// ComputationSummary is a lightweight history item (no matrices in the payload).
type ComputationSummary struct {
	ID        string    `json:"id"`
	Rows      int       `json:"rows"`
	Cols      int       `json:"cols"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

// ComputationList is the paginated history response.
type ComputationList struct {
	Items  []ComputationSummary `json:"items"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
	Total  int                  `json:"total"`
}
