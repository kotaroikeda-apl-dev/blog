package services_test

import (
	"blog/models"
	"blog/services"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) FindAll() ([]models.Post, error) {
	args := m.Called()
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockPostRepository) FindByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *MockPostRepository) Create(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) Update(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func TestGetAllPosts(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	repo.On("FindAll").Return([]models.Post{{ID: 1, Title: "Test Post", Content: "Test Content", Author: "Test Author", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil)

	posts, err := service.GetAllPosts()
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Test Post", posts[0].Title)
	assert.Equal(t, "Test Content", posts[0].Content)
	assert.Equal(t, "Test Author", posts[0].Author)
}

func TestGetPostByID(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	repo.On("FindByID", uint(1)).Return(&models.Post{ID: 1, Title: "Test Post", Content: "Test Content", Author: "Test Author", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)

	post, err := service.GetPostByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.Equal(t, "Test Author", post.Author)
}

func TestGetPostByID_NotFound(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)

	repo.On("FindByID", uint(99)).Return((*models.Post)(nil), errors.New("post not found"))

	post, err := service.GetPostByID(99)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestCreatePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	post := &models.Post{Title: "New Post", Content: "New Content", Author: "New Author"}
	repo.On("Create", mock.Anything).Return(nil)

	err := service.CreatePost(post)
	assert.NoError(t, err)
	assert.NotZero(t, post.CreatedAt)
	assert.NotZero(t, post.UpdatedAt)
	assert.Equal(t, "New Content", post.Content)
	assert.Equal(t, "New Author", post.Author)
}

func TestCreatePost_Error(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	post := &models.Post{Title: "New Post", Content: "New Content", Author: "New Author"}

	repo.On("Create", mock.Anything).Return(errors.New("failed to create post"))

	err := service.CreatePost(post)
	assert.Error(t, err)
}

func TestUpdatePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	existingPost := &models.Post{ID: 1, Title: "Old Title", Content: "Old Content", Author: "Old Author"}
	repo.On("FindByID", uint(1)).Return(existingPost, nil)
	repo.On("Update", mock.Anything).Return(nil)

	updatedPost := models.Post{Title: "Updated Title", Content: "Updated Content", Author: "Updated Author"}

	existingPost.Title = updatedPost.Title
	existingPost.Content = updatedPost.Content
	existingPost.Author = updatedPost.Author
	existingPost.UpdatedAt = time.Now()

	err := service.UpdatePost(1, updatedPost)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", existingPost.Title)
	assert.Equal(t, "Updated Content", existingPost.Content)
	assert.Equal(t, "Updated Author", existingPost.Author)
}

func TestUpdatePost_NotFound(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)

	repo.On("FindByID", uint(99)).Return((*models.Post)(nil), errors.New("post not found"))

	updatedPost := models.Post{Title: "Updated Title", Content: "Updated Content", Author: "Updated Author"}
	err := service.UpdatePost(99, updatedPost)

	assert.Error(t, err)
}

func TestUpdatePost_Error(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)

	existingPost := &models.Post{ID: 1, Title: "Old Title", Content: "Old Content", Author: "Old Author"}
	repo.On("FindByID", uint(1)).Return(existingPost, nil)
	repo.On("Update", mock.Anything).Return(errors.New("failed to update post"))

	updatedPost := models.Post{Title: "Updated Title", Content: "Updated Content", Author: "Updated Author"}
	err := service.UpdatePost(1, updatedPost)

	assert.Error(t, err)
}

func TestDeletePost(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)
	existingPost := &models.Post{ID: 1, Title: "Test Post", Content: "Test Content", Author: "Test Author"}
	repo.On("FindByID", uint(1)).Return(existingPost, nil)
	repo.On("Delete", mock.Anything).Return(nil)

	err := service.DeletePost(1)
	assert.NoError(t, err)
}

func TestDeletePost_NotFound(t *testing.T) {
	repo := new(MockPostRepository)
	service := services.NewPostService(repo)

	repo.On("FindByID", uint(99)).Return((*models.Post)(nil), errors.New("post not found"))

	err := service.DeletePost(99)
	assert.Error(t, err)
}
