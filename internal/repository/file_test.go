// repository/file_repository_test.go
package repository_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"chat-apps/internal/util"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var file = domain.File{
	UserID:  1,
	FileURL: "http://example.com/file",
}

/*
	func TestUploadFile(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect database: %v", err)
		}
		db.AutoMigrate(&domain.File{})

		fileRepo := repository.NewFileRepository(db)

		createdFile, err := fileRepo.UploadFile(file)
		if err != nil {
			t.Fatalf("Failed to upload file: %v", err)
		}

		assert.NoError(t, err)
		assert.Equal(t, file.UserID, createdFile.UserID)
		assert.Equal(t, file.FileURL, createdFile.FileURL)
		assert.WithinDuration(t, time.Now(), createdFile.UploadedAt, time.Second)
	}
*/
func TestGetFileByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	db.AutoMigrate(&domain.File{})

	fileRepo := repository.NewFileRepository(db)
	createdFile, _ := fileRepo.UploadFile(file)

	retrievedFile, err := fileRepo.GetFileByID(createdFile.ID)
	if err != nil {
		t.Fatalf("Failed to get file by ID: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, createdFile.ID, retrievedFile.ID)
	assert.Equal(t, createdFile.UserID, retrievedFile.UserID)
	assert.Equal(t, createdFile.FileURL, retrievedFile.FileURL)
}

func TestUploadFileWithSQLMock(t *testing.T) {
	t.Run("test positif upload file with sqlmock", func(t *testing.T) {
		db, mock := util.SetupMockDB()

		file.UploadedAt = time.Now()
		mock.ExpectBegin()
		// rows := sqlmock.NewRows([]string{"id", "user_id", "file_url", "uploaded_at"}).
		// 	AddRow(1, file.UserID, file.FileURL, file.UploadedAt)

		mock.ExpectQuery(`INSERT INTO "files" ("user_id","file_url","uploaded_at") VALUES ($1,$2,$3) RETURNING "id"`).
			// WithArgs(file.UserID, file.FileURL, file.UploadedAt).
			WithArgs(file.UserID, file.FileURL, sqlmock.AnyArg()). // Use sqlmock.AnyArg() for UploadedAt
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		fileRepo := repository.NewFileRepository(db)

		createdFile, err := fileRepo.UploadFile(file)
		if err != nil {
			t.Fatalf("Failed to upload file: %v", err)
		}

		t.Log("[createdFile]:", createdFile)
		t.Run("test store data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
			assert.Equal(t, file.UserID, createdFile.UserID)
			assert.Equal(t, file.FileURL, createdFile.FileURL)
			assert.WithinDuration(t, time.Now(), createdFile.UploadedAt, time.Second)
		})
	})
}

func TestGetFileWithSQLMock(t *testing.T) {
	t.Run("test positif get file with sqlmock", func(t *testing.T) {
		db, mock := util.SetupMockDB()

		file.UploadedAt = time.Now()
		// mock.ExpectBegin()
		id := 1
		rows := sqlmock.NewRows([]string{"id", "user_id", "file_url", "uploaded_at"}).
			AddRow(id, file.UserID, file.FileURL, file.UploadedAt)
		mock.ExpectQuery(`SELECT * FROM "files" WHERE "files"."id" = $1 ORDER BY "files"."id" LIMIT $2`).
			WithArgs(id, 1). // Use sqlmock.AnyArg() for UploadedAt
			WillReturnRows(rows)
		// mock.ExpectCommit()

		fileRepo := repository.NewFileRepository(db)
		retrievedFile, err := fileRepo.GetFileByID(id)
		if err != nil {
			t.Fatalf("Failed to get file by ID: %v", err)
		}

		t.Log("[retrievedFile]:", retrievedFile)
		t.Run("test store data with no error", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, id, retrievedFile.ID)
			assert.Equal(t, file.UserID, retrievedFile.UserID)
			assert.Equal(t, file.FileURL, retrievedFile.FileURL)
		})
	})
}

func TestUploadFile_Error(t *testing.T) {
	t.Run("test negatif - upload file", func(t *testing.T) {
		db, mock := util.SetupMockDB()
		repo := repository.NewFileRepository(db)
		mock.ExpectExec("INSERT INTO files(user_id, file) VALUES (3, file.png)").
			WillReturnError(assert.AnError)
		files, err := repo.UploadFile(file)
		t.Run("test return error upload file", func(t *testing.T) {
			assert.Error(t, err)
			assert.Equal(t, domain.File{}, files)
		})
	})
}

func TestGetFileByID_Error(t *testing.T) {
	t.Run("test negatif - get file by id", func(t *testing.T) {
		id := 3
		db, mock := util.SetupMockDB()
		repo := repository.NewFileRepository(db)
		mock.ExpectQuery("SELECT * FROM files where id = ?").
			WithArgs(id).
			WillReturnError(assert.AnError)
		files, err := repo.GetFileByID(id)
		t.Run("test return error upload file", func(t *testing.T) {
			assert.Error(t, err)
			assert.Equal(t, domain.File{}, files)
		})
	})
}
