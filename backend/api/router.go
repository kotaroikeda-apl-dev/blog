package api

import (
	"blog/controllers"
	"blog/middlewares"
	"blog/repositories"
	"blog/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	r.Use(middlewares.CORSConfig())

	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)
	postController := controllers.NewPostController(service)

	r.HandleFunc("/api/posts", postController.GetAllPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id:[0-9]+}", postController.GetPostByID).Methods("GET")
	r.HandleFunc("/api/posts", postController.CreatePost).Methods("POST")
	r.HandleFunc("/api/posts/{id:[0-9]+}", postController.UpdatePost).Methods("PUT")
	r.HandleFunc("/api/posts/{id:[0-9]+}", postController.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/posts/{id:[0-9]+}/render", postController.RenderMarkdown).Methods("GET")

	// **OPTIONS リクエストを許可**
	r.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	r.HandleFunc("/api/posts/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// ヘルスチェック
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}).Methods("GET")

	return r
}
