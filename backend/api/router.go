package api

import (
	"blog/controllers"
	"blog/middlewares"
	"blog/repositories"
	"blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS ミドルウェアを適用
	r.Use(middlewares.CORSConfig())

	// リポジトリ、サービス、コントローラーの初期化
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)
	postController := controllers.NewPostController(service)

	// ルートグループを作成
	api := r.Group("/api/posts")
	{
		api.GET("", postController.GetAllPosts)
		api.GET("/:id", postController.GetPostByID)
		api.POST("", postController.CreatePost)
		api.PUT("/:id", postController.UpdatePost)
		api.DELETE("/:id", postController.DeletePost)
		api.GET("/:id/render", postController.RenderMarkdown)
	}

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return r
}
