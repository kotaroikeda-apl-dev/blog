package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// インストール済みである必要あり
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"blog/api"
	"blog/models" // モデルのパスを合わせる
)

var db *gorm.DB

var (
	host     = os.Getenv("DATABASE_HOST")
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")
	dbName   = os.Getenv("DATABASE_NAME")
	port     = "5432"
	dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
)

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB接続エラー: %v", err)
	}

	// マイグレーション
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("マイグレーションエラー: %v", err)
	}
	log.Println("マイグレーションが成功しました。")

	// ルートの登録
	r := api.RegisterRoutes(db)

	// サーバーの起動
	log.Println("サーバーを起動します。ポート: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
