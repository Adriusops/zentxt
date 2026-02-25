package versioning

import (
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

type Version struct {
	ID            string
	FileID        string
	VersionNumber int
	Path          string
	Author        string
	Message       string
	Content       string
	CreatedAt     string
}

func SaveVersion(db *sql.DB, fileID string, path string, author string, message string, content string) (*Version, error) {
	// 1. Générer un UUID pour l'id
	id := uuid.New().String()
	// 2. Insérer en DB avec db.Exec nouvelle version dans version
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM versions WHERE file_id = ?", fileID).Scan(&count)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("INSERT INTO versions (id, file_id, version_number, path, author, message, content, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", id, fileID, count+1, path, author, message, content, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	// 3. Update le current version id dans files
	_, err = db.Exec("UPDATE files SET current_version_id = ? WHERE id = ?", id, fileID)
	if err != nil {
		return nil, err
	}

	return &Version{
		ID:            id,
		FileID:        fileID,
		VersionNumber: count + 1,
		Path:          path,
		Author:        author,
		Message:       message,
		Content:       content,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}, nil
}

func ListVersions(db *sql.DB, fileID string) ([]*Version, error) {
	rows, err := db.Query("SELECT id, file_id, version_number, path, author, message, content, created_at FROM versions WHERE file_id = ?", fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*Version
	for rows.Next() {
		version := &Version{}
		err := rows.Scan(&version.ID, &version.FileID, &version.VersionNumber, &version.Path, &version.Author, &version.Message, &version.Content, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func GetVersion(db *sql.DB, id string) (*Version, error) {
	row := db.QueryRow("SELECT id, file_id, version_number, path, author, message, content, created_at FROM versions WHERE id = ?", id)
	var version Version
	if err := row.Scan(&version.ID, &version.FileID, &version.VersionNumber, &version.Path, &version.Author, &version.Message, &version.Content, &version.CreatedAt); err != nil {
		return nil, err
	}
	return &version, nil
}

func RestoreVersion(db *sql.DB, fileID string, versionID string) (*Version, error) {
	currentVersion, err := GetVersion(db, versionID)
	if err != nil {
		return nil, err
	}

	currentFile, err := GetFile(db, fileID)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("UPDATE files SET current_version_id = ? WHERE id = ?", currentVersion.ID, currentFile.ID)
	if err != nil {
		return nil, err
	}

	src, err := os.Open(currentVersion.Path)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(currentFile.Path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	return currentVersion, nil
}
