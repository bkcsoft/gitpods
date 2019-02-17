package pubkey

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is returned when no pubkey is found
	ErrNotFound = errors.New("pubkey not found")
)

type (
	// Store handles the storage of PubKeys
	Store interface {
		Find(ctx context.Context, id string) (*PubKey, error)
		FindByFingerprint(ctx context.Context, fingerprint string) (*PubKey, error)
		Create(ctx context.Context, pk *PubKey) (*PubKey, error)
	}

	// Service handles all interactions with PubKeys
	Service interface {
		Find(ctx context.Context, id string) (*PubKey, error)
		FindByFingerprint(ctx context.Context, fingerprint string) (*PubKey, error)
		Create(ctx context.Context, pk *PubKey) (*PubKey, error)
	}

	service struct {
		store Store
	}
)

func (s *service) Find(ctx context.Context, id string) (*PubKey, error) {
	return s.store.Find(ctx, id)
}

func (s *service) FindByFingerprint(ctx context.Context, fingerprint string) (*PubKey, error) {
	return s.store.FindByFingerprint(ctx, fingerprint)
}

func (s *service) Create(ctx context.Context, pk *PubKey) (*PubKey, error) {
	if err := validateCreate(pk); err != nil {
		return nil, err
	}
	if err := pk.fillFingerprint(); err != nil {
		return nil, err
	}

	return s.store.Create(ctx, pk)
}
