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
	GetOne(int) (models.Ad, error)
	GetList(curr int, limit int, offset int) ([]models.Ad, error)
}
