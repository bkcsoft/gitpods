package pubkey

import (
	"context"
	"database/sql"
	"time"

	"github.com/opentracing/opentracing-go"
)

// Postgres implementation of the Store.
type Postgres struct {
	db *sql.DB
}

// NewPostgresStore returns a Postgres implementation of the Store.
func NewPostgresStore(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// Find a pubkey by its ID.
func (s *Postgres) Find(ctx context.Context, id string) (*PubKey, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pubkey.Postgres.Find")
	span.SetTag("id", id)
	defer span.Finish()

	findByID := `
SELECT
	name,
	content,
	fingerprint,
	created_at,
	updated_at
FROM pubkeys
WHERE id = $1
LIMIT 1;
`

	row := s.db.QueryRowContext(ctx, findByID, id)

	var name string
	var content string
	var fingerprint string
	var created time.Time
	var updated time.Time
	if err := row.Scan(&name, &content, &fingerprint, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &PubKey{
		ID:          id,
		Name:        name,
		Fingerprint: fingerprint,
		Content:     content,
		Created:     created,
		Updated:     updated,
	}, nil
}

// FindByFingerprint is not implemented yet...
func (s *Postgres) FindByFingerprint(ctx context.Context, fingerprint string) (*PubKey, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pubkey.Postgres.FindByFingerprint")
	span.SetTag("fingerprint", fingerprint)
	defer span.Finish()
	panic("implement me!")
}

// Create a user in postgres and return it with the ID set in the store.
func (s *Postgres) Create(ctx context.Context, pk *PubKey) (*PubKey, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pubkey.Postgres.Create")
	span.SetTag("fingerprint", pk.Fingerprint)
	defer span.Finish()

	create := `
INSERT INTO pubkeys (name, fingerprint, content) VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;
`

	err := s.db.QueryRowContext(
		ctx,
		create,
		pk.Name, pk.Fingerprint, pk.Content,
	).Scan(&pk.ID, &pk.Created, &pk.Updated)

	return pk, err
}
