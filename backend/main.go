package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"      // インストール済みである必要あり
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "blog/models"                 // モデルのパスを合わせる
    "github.com/gomarkdown/markdown"

    "github.com/gorilla/handlers"
)

var db *gorm.DB

func initDB() {
    dsn := "host=localhost user=kotaroikeda password=kotaro220 dbname=blogdb port=5432 sslmode=disable TimeZone=Asia/Tokyo"
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("データベース接続エラー: %v", err)
    }

    // マイグレーション
    err = db.AutoMigrate(&models.Post{})
    if err != nil {
        log.Fatalf("マイグレーションエラー: %v", err)
    }
    log.Println("マイグレーションが成功しました。")
}

// 全ての投稿を取得するハンドラ
func getAllPosts(w http.ResponseWriter, r *http.Request) {
    var posts []models.Post
    if err := db.Find(&posts).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

// 特定の投稿を取得するハンドラ
func getPostByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    if err := db.First(&post, id).Error; err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// 新規投稿を保存するハンドラ
func createPost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := db.Create(&post).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// 特定の投稿を更新するハンドラ
func updatePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    if err := db.First(&post, id).Error; err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    post.ID = uint(id) // 更新する投稿のIDを再設定

    if err := db.Save(&post).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// 特定の投稿を削除するハンドラ
func deletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    if err := db.First(&post, id).Error; err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    if err := db.Delete(&post).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 成功時は204 No Contentを返す
}

// MarkdownコンテンツをHTMLで返すハンドラ
func renderMarkdown(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := db.First(&post, id).Error; err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    // MarkdownをHTMLに変換
    htmlContent := markdown.ToHTML([]byte(post.Content), nil, nil)

    w.Header().Set("Content-Type", "text/html")
    w.Write(htmlContent)
}

func main() {
    initDB()

    r := mux.NewRouter()
    r.HandleFunc("/api/posts", getAllPosts).Methods("GET")
    r.HandleFunc("/api/posts/{id}", getPostByID).Methods("GET")
    r.HandleFunc("/api/posts", createPost).Methods("POST")         // 新規投稿エンドポイント
    r.HandleFunc("/api/posts/{id}", updatePost).Methods("PUT")
    r.HandleFunc("/api/posts/{id}", deletePost).Methods("DELETE")
    r.HandleFunc("/api/posts/{id}/render", renderMarkdown).Methods("GET") // MarkdownをHTMLに変換するエンドポイント

	// CORSの設定を直接 `handlers.CORS()` に渡す
	log.Println("サーバーを起動します。ポート: 8080")
	log.Fatal(http.ListenAndServe(":8080",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"http://localhost:3000"}),  // フロントエンドからのリクエストを許可
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // 許可するHTTPメソッド
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),          // 許可するヘッダー
		)(r)))
}

