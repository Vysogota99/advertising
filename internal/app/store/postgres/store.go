package postgres

import (
	"github.com/Vysogota99/advertising/internal/app/store"
	"github.com/jmoiron/sqlx"
)

// Store - postgres implementation of store
type Store struct {
	DB            *sqlx.DB
	adRepository store.AdRepository
}

// New - helper to init Store
func New(db *sqlx.DB) *Store {
	return &Store{
		DB: db,
	}
}

// Add - implementation of AddRepository interface
func (s *Store) Add() store.AdRepository {
	if s.adRepository == nil {
		s.adRepository = &AdRepository{
			store: s,
		}
	}

	return s.adRepository
}
