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
func (a *AdRepository) GetOne(int) (models.Ad, error) {
	result := models.Ad{}
	return result, nil
}

// GetList ...
func (a *AdRepository) GetList(curr int, limit int, offset int) ([]models.Ad, error) {
	return nil, nil
}
