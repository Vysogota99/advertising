package store

import (
	"github.com/Vysogota99/advertising/internal/app/models"
)

// Store - интерфейс для работы с хранилищем
type Store interface {
	Add() AdRepository
}

// AdRepository - интерфейс, содержащий методы для получаения информации об объявлениях
type AdRepository interface {
	Create(models.Ad) (int, error)
	GetOne(id int, description, photos bool) (*models.Ad, error)
	GetList(limit int, offset int, sortBy, sortDirection string) ([]models.Ad, int, error)
}
