package repositories

import (
	"blog/models"
	"fmt"

	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	return posts, nil
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	return &post, nil
}

func (r *postRepository) Create(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *postRepository) Update(post *models.Post) error {
	if err := r.db.Save(post).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *postRepository) Delete(post *models.Post) error {
	if err := r.db.Delete(post).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}
