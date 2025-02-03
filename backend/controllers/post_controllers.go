package controllers

import (
	"blog/models"
	"blog/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service: service}
}

// 全ての投稿を取得
func (c *PostController) GetAllPosts(ctx *gin.Context) {
	posts, err := c.service.GetAllPosts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

// ID から投稿を取得
func (c *PostController) GetPostByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	post, err := c.service.GetPostByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// 新規投稿を作成
func (c *PostController) CreatePost(ctx *gin.Context) {
	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := c.service.CreatePost(&post); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}

// 投稿を更新
func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := c.service.UpdatePost(id, post); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// 投稿を削除
func (c *PostController) DeletePost(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeletePost(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// Markdown を HTML に変換して表示
func (c *PostController) RenderMarkdown(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	post, err := c.service.GetPostByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	htmlContent := markdown.ToHTML([]byte(post.Content), nil, nil)
	ctx.Data(http.StatusOK, "text/html", htmlContent)
}

// Gin のパスパラメータから ID を取得
func parseID(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
