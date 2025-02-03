package controllers_test

import (
	"blog/controllers"
	"blog/models"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) GetAllPosts() ([]models.Post, error) {
	args := m.Called()
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockPostService) GetPostByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *MockPostService) CreatePost(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostService) UpdatePost(id uint, postData models.Post) error {
	args := m.Called(id, postData)
	return args.Error(0)
}

func (m *MockPostService) DeletePost(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAllPosts(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	service.On("GetAllPosts").Return([]models.Post{
		{ID: 1, Title: "Test Post", Content: "Test Content", Author: "Test Author", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil)

	controller.GetAllPosts(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetAllPosts_Error(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	service.On("GetAllPosts").Return([]models.Post{}, errors.New("database error"))

	controller.GetAllPosts(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetPostByID(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	service.On("GetPostByID", uint(1)).Return(&models.Post{ID: 1, Title: "Test Post", Content: "Test Content", Author: "Test Author", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)

	controller.GetPostByID(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetPostByID_NotFound(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "999"})

	service.On("GetPostByID", uint(999)).Return((*models.Post)(nil), errors.New("not found"))

	controller.GetPostByID(ctx)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetPostByID_InvalidID(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "invalid"})

	controller.GetPostByID(ctx)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreatePost(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	// リクエストボディを設定
	postJSON := `{"title": "New Post", "content": "New Content", "author": "New Author"}`
	ctx.Request = httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(postJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	service.On("CreatePost", mock.Anything).Return(nil)

	controller.CreatePost(ctx)

	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreatePost_Error(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	postJSON := `{"title": "New Post", "content": "New Content", "author": "New Author"}`
	ctx.Request = httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(postJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	service.On("CreatePost", mock.Anything).Return(errors.New("insert error"))

	controller.CreatePost(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCreatePost_InvalidJSON(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString("invalid json"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	controller.CreatePost(ctx)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdatePost(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	// JSON データを設定
	postJSON := `{"title": "Updated Title", "content": "Updated Content", "author": "Updated Author"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString(postJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	service.On("UpdatePost", uint(1), mock.Anything).Return(nil)

	controller.UpdatePost(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdatePost_Error(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	postJSON := `{"title": "Updated Title", "content": "Updated Content", "author": "Updated Author"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString(postJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	service.On("UpdatePost", uint(1), mock.Anything).Return(errors.New("update error"))

	controller.UpdatePost(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestUpdatePost_InvalidJSON(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	ctx.Request = httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString("invalid json"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	controller.UpdatePost(ctx)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdatePost_NotFound(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "999"})

	postJSON := `{"title": "Updated Title", "content": "Updated Content", "author": "Updated Author"}`
	ctx.Request = httptest.NewRequest(http.MethodPut, "/posts/999", bytes.NewBufferString(postJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	service.On("UpdatePost", uint(999), mock.Anything).Return(errors.New("not found"))

	controller.UpdatePost(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeletePost(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	service.On("DeletePost", uint(1)).Return(nil)

	controller.DeletePost(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestDeletePost_NotFound(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "999"})

	service.On("DeletePost", uint(999)).Return(errors.New("not found"))

	controller.DeletePost(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeletePost_InvalidID(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "invalid"})

	controller.DeletePost(ctx)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeletePost_Error(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	service.On("DeletePost", uint(1)).Return(errors.New("delete error"))

	controller.DeletePost(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestRenderMarkdown(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

	post := &models.Post{
		ID:      1,
		Title:   "Markdown Test",
		Content: "# Hello\nThis is a test",
	}

	service.On("GetPostByID", uint(1)).Return(post, nil)

	controller.RenderMarkdown(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<h1>Hello</h1>")
}

func TestRenderMarkdown_NotFound(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "999"})

	service.On("GetPostByID", uint(999)).Return((*models.Post)(nil), errors.New("not found"))

	controller.RenderMarkdown(ctx)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Post not found")
}

func TestRenderMarkdown_InvalidID(t *testing.T) {
	service := new(MockPostService)
	controller := controllers.NewPostController(service)
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "abc"})

	controller.RenderMarkdown(ctx)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid ID")
}
