package repositories_test

import (
	"blog/models"
	"blog/repositories"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// テスト用にDBをモックする
func setupMockDB(t *testing.T) (repositories.PostRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open mock GORM DB: %v", err)
	}

	repo := repositories.NewPostRepository(gormDB)
	return repo, mock
}

func TestFindAll(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "posts"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Test Post 1", "Content 1", "Author 1", time.Now(), time.Now(), nil).
			AddRow(2, "Test Post 2", "Content 2", "Author 2", time.Now(), time.Now(), nil))

	posts, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Test Post 1", posts[0].Title)
	assert.Equal(t, "Content 1", posts[0].Content)
	assert.Equal(t, "Author 1", posts[0].Author)
	assert.Equal(t, "Test Post 2", posts[1].Title)
	assert.Equal(t, "Content 2", posts[1].Content)
	assert.Equal(t, "Author 2", posts[1].Author)
}

func TestFindAll_Error(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "posts"`).
		WillReturnError(errors.New("database error"))

	posts, err := repo.FindAll()
	assert.Error(t, err)
	assert.Nil(t, posts)
}

func TestFindByID(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"."id" = \$1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(1, 1). // `LIMIT $2` を考慮
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Test Post", "Test Content", "Test Author", time.Now(), time.Now(), nil))

	post, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.Equal(t, "Test Author", post.Author)
}

func TestFindByID_NotFound(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"."id" = \$1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"})) // 空の結果

	post, err := repo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestCreate(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin() // トランザクション開始

	mock.ExpectQuery(`INSERT INTO "posts" \("title","content","author","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
		WithArgs("New Post", "New Content", "New Author", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit() // トランザクションコミット

	post := &models.Post{Title: "New Post", Content: "New Content", Author: "New Author"}
	err := repo.Create(post)
	assert.NoError(t, err)
}

func TestCreate_Error(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin()

	// データベースエラーを発生させる
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs("New Post", "New Content", "New Author", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(errors.New("failed to insert post"))

	mock.ExpectRollback()

	post := &models.Post{Title: "New Post", Content: "New Content", Author: "New Author"}
	err := repo.Create(post)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin() // トランザクション開始

	mock.ExpectExec(`UPDATE "posts" SET "title"=\$1,"content"=\$2,"author"=\$3,"created_at"=\$4,"updated_at"=\$5,"deleted_at"=\$6 WHERE "posts"."deleted_at" IS NULL AND "id" = \$7`).
		WithArgs("Updated Post", "Updated Content", "Updated Author", sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit() // トランザクションコミット

	post := &models.Post{ID: 1, Title: "Updated Post", Content: "Updated Content", Author: "Updated Author"}
	err := repo.Update(post)
	assert.NoError(t, err)
}

func TestUpdate_Error(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin()

	// データベースエラーを発生させる
	mock.ExpectExec(`UPDATE "posts"`).
		WithArgs("Updated Post", "Updated Content", "Updated Author", sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1).
		WillReturnError(errors.New("failed to update post"))

	mock.ExpectRollback()

	post := &models.Post{ID: 1, Title: "Updated Post", Content: "Updated Content", Author: "Updated Author"}
	err := repo.Update(post)
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin() // トランザクション開始

	mock.ExpectExec(`UPDATE "posts" SET "deleted_at"=\$1 WHERE "posts"."id" = \$2 AND "posts"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit() // トランザクションコミット

	post := &models.Post{ID: 1}
	err := repo.Delete(post)
	assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
	repo, mock := setupMockDB(t)

	mock.ExpectBegin()

	// データベースエラーを発生させる
	mock.ExpectExec(`UPDATE "posts" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnError(errors.New("failed to delete post"))

	mock.ExpectRollback()

	post := &models.Post{ID: 1}
	err := repo.Delete(post)
	assert.Error(t, err)
}
