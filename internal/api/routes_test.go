package api_test

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/Adriusops/zentxt/internal/api"
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
