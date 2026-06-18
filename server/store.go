package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Document struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"clientId"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	SizeBytes   int64     `json:"sizeBytes"`
	StorageKey  string    `json:"-"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Store struct {
	pool      *pgxpool.Pool
	uploadDir string
}

type Client struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewStore(pool *pgxpool.Pool, uploadDir string) (*Store, error) {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &Store{pool: pool, uploadDir: uploadDir}, nil
}

func (s *Store) SaveFile(storageKey string, r io.Reader) (int64, error) {
	path := filepath.Join(s.uploadDir, storageKey)
	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		_ = os.Remove(path)
		return 0, fmt.Errorf("write file: %w", err)
	}
	return n, nil
}

func (s *Store) CreateDocument(ctx context.Context, doc *Document) error {
	const q = `
		INSERT INTO documents (client_id, filename, content_type, size_bytes, storage_key, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return s.pool.QueryRow(ctx, q,
		doc.ClientID, doc.Filename, doc.ContentType,
		doc.SizeBytes, doc.StorageKey, doc.Status,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
}

const docColumns = `
	id,
	client_id,
	filename,
	content_type,
	size_bytes,
	storage_key,
	status,
	created_at,
	updated_at
`

func (s *Store) getDocumentsByClient(ctx context.Context, clientID string) ([]Document, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+docColumns+` FROM documents WHERE client_id = $1 ORDER BY created_at DESC`,
		clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectDocs(rows)
}

func (s *Store) getAllDocuments(ctx context.Context) ([]Document, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+docColumns+` FROM documents ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectDocs(rows)
}

func collectDocs(rows pgx.Rows) ([]Document, error) {
	docs := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.ID, &d.ClientID, &d.Filename, &d.ContentType,
			&d.SizeBytes, &d.StorageKey, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (s *Store) GetClientByEmail(ctx context.Context, email string) (*Client, error) {
	const q = `SELECT id, name, email FROM clients WHERE email = $1`
	var c Client
	err := s.pool.QueryRow(ctx, q, email).Scan(&c.ID, &c.Name, &c.Email)
	return &c, err 
}