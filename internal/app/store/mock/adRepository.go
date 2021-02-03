package mock

import "github.com/Vysogota99/advertising/internal/app/models"

// AdRepository - postgres implementation of AddRepository
type AdRepository struct {
	store *Store
}

// Create ...
func (a *AdRepository) Create(models.Ad) (int, error) {
	return 0, nil
}

// GetOne ...
func (a *AdRepository) GetOne(id int, description, photos bool) (*models.Ad, error) {
	result := models.Ad{}
	return &result, nil
}

// GetList ...
func (a *AdRepository) GetList(limit int, offset int, sortBy, sortDirection string) ([]models.Ad, int, error) {
	return nil, 0, nil
}
