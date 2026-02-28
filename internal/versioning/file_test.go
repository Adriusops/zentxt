package versioning

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
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
	// Return connection
	return db
}

func TestCreateFile(t *testing.T) {
	db := setupTestDB(t)
	file, err := CreateFile(db, "test.txt", "/tmp/test.txt", nil)
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", file.Name)
	assert.NotEmpty(t, file.ID)
}

func TestCreateFile_EmptyName(t *testing.T) {
	db := setupTestDB(t)
	file, err := CreateFile(db, "", "/tmp/test.txt", nil)
	assert.Error(t, err)
	assert.Nil(t, file)
}
