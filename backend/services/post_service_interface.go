package services

import "blog/models"

type PostService interface {
	GetAllPosts() ([]models.Post, error)
	GetPostByID(id uint) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(id uint, postData models.Post) error
	DeletePost(id uint) error
}
