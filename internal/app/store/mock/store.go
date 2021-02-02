package mock

import (
	"github.com/Vysogota99/advertising/internal/app/store"
)

// Store - postgres implementation of store
type Store struct {
	addRepository store.AdRepository
}

// New - helper to init Store
func New() *Store {
	return &Store{}
}

// Add - implementation of AdRepository interface
func (s *Store) Add() store.AdRepository {
	if s.addRepository == nil {
		s.addRepository = &AdRepository{
			store: s,
		}
	}

	return s.addRepository
}
