package api_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Adriusops/zentxt/internal/api"
	"github.com/Adriusops/zentxt/internal/versioning"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

func setupTestApp(t *testing.T) (*fiber.App, *sql.DB) {
	// Open memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// Execute migrations
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (id TEXT PRIMARY KEY,name TEXT NOT NULL,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS files (id TEXT PRIMARY KEY,name TEXT NOT NULL,path TEXT NOT NULL,current_version_id TEXT REFERENCES versions(id) NULL,project_id TEXT REFERENCES projects(id),created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS versions (id TEXT PRIMARY KEY, file_id TEXT REFERENCES files(id),version_number INTEGER NOT NULL,path TEXT NOT NULL, author TEXT NULL,message TEXT NULL,content TEXT NOT NULL,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		t.Fatal(err)
	}
	// Create app
	app := fiber.New()

	api.SetupRoutes(app, db)

	// Return connection
	return app, db
}

func TestCreateFile_Success(t *testing.T) {
	app, _ := setupTestApp(t)

	req, err := http.NewRequest("POST", "/api/files", strings.NewReader(`{"name": "test", "path": "/tmp/test.txt"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateFile_ErrorNoParam(t *testing.T) {
	app, _ := setupTestApp(t)

	req, err := http.NewRequest("POST", "/api/files", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateFile_ErrorName(t *testing.T) {
	app, _ := setupTestApp(t)

	req, err := http.NewRequest("POST", "/api/files", strings.NewReader(`{"name": "", "path": "/tmp/test.txt"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateFile_ErrorPath(t *testing.T) {
	app, _ := setupTestApp(t)

	req, err := http.NewRequest("POST", "/api/files", strings.NewReader(`{"name": "test", "path": ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSaveVersion_Success(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files/%s/versions", file.ID), strings.NewReader(`{"path": "/tmp/test.txt", "author": "You", "message": "Initial Commit", "content": "blablabla" }`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSaveVersion_ErrorNoPath(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files/%s/versions", file.ID), strings.NewReader(`{"path": "", "author": "You", "message": "Initial Commit", "content": "blablabla" }`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSaveVersion_ErrorNoContent(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files/%s/versions", file.ID), strings.NewReader(`{"path": "/tmp/test.txt", "author": "You", "message": "Initial Commit", "content": "" }`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSaveVersion_ErrorInvalidFileID(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	_, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/files/100/versions", strings.NewReader(`{"path": "/tmp/test.txt", "author": "You", "message": "Initial Commit", "content": "blablabla" }`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRestoreVersion_Success(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	version, err := versioning.SaveVersion(db, file.ID, "/tmp/test.txt", "You", "Initial Commit", "blablabla")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/files/%s/restore/%s", file.ID, version.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestoreVersion_InvalidFileID(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	version, err := versioning.SaveVersion(db, file.ID, "/tmp/test.txt", "You", "Initial Commit", "blablabla")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/files/59/restore/%s", version.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRestoreVersion_InvalidVersionID(t *testing.T) {
	app, db := setupTestApp(t)

	// Create file in db
	file, err := versioning.CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = versioning.SaveVersion(db, file.ID, "/tmp/test.txt", "You", "Initial Commit", "blablabla")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/files/%s/restore/89", file.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
