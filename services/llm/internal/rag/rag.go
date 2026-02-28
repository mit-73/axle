// Package rag provides Retrieval-Augmented Generation helpers.
// Currently a skeleton — real vector-store retrieval wired when pgvector is available.
package rag

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Retriever fetches relevant document chunks for a given query.
type Retriever struct {
	pool *pgxpool.Pool
}

// New creates a new Retriever backed by a pgvector-enabled Postgres pool.
func New(pool *pgxpool.Pool) *Retriever {
	return &Retriever{pool: pool}
}

// Chunk is a retrieved document excerpt.
type Chunk struct {
	ID      string
	Content string
	Score   float64
}

// Retrieve returns the top-k most relevant chunks for the given query.
// TODO: replace stub with pgvector similarity search once schema is migrated.
func (r *Retriever) Retrieve(_ context.Context, query string, topK int) ([]Chunk, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("rag: db pool not configured")
	}
	if query == "" {
		return nil, fmt.Errorf("rag: query is empty")
	}
	if topK <= 0 {
		topK = 5
	}

	// Stub — returns empty slice until pgvector schema is in place.
	_ = topK
	return []Chunk{}, nil
}
