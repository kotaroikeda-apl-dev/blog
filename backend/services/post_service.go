package services

import (
	"blog/models"
	"blog/repositories"
	"time"
)

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetAllPosts() ([]models.Post, error) {
	return s.repo.FindAll()
}

func (s *postService) GetPostByID(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *postService) CreatePost(post *models.Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	return s.repo.Create(post)
}

func (s *postService) UpdatePost(id uint, postData models.Post) error {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	post.Title = postData.Title
	post.Content = postData.Content
	post.UpdatedAt = time.Now()

	return s.repo.Update(post)
}

func (s *postService) DeletePost(id uint) error {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(post)
}
