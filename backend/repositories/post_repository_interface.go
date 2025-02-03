package repositories

import "blog/models"

type PostRepository interface {
	FindAll() ([]models.Post, error)
	FindByID(id uint) (*models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(post *models.Post) error
}
