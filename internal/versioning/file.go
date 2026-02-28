package versioning

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        string
	Name      string
	Path      string
	ProjectID *string
	CreatedAt string
}

var ErrNotFound = errors.New("not found")

func CreateFile(db *sql.DB, name string, path string, projectID *string) (*File, error) {
	if name == "" || path == "" {
		return nil, errors.New("name and path are required")
	}
	// 1. Générer un UUID pour l'id
	id := uuid.New().String()
	// 2. Insérer en DB avec db.Exec
	_, err := db.Exec("INSERT INTO files (id, name, path, project_id, created_at) VALUES (?, ?, ?, ?, ?)", id, name, path, projectID, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	// 3. Retourner le File créé
	return &File{
		ID:        id,
		Name:      name,
		Path:      path,
		ProjectID: projectID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func ListFiles(db *sql.DB) ([]*File, error) {
	rows, err := db.Query("SELECT id, name, path, project_id, created_at FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.ID, &file.Name, &file.Path, &file.ProjectID, &file.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

func GetFile(db *sql.DB, id string) (*File, error) {
	row := db.QueryRow("SELECT id, name, path, project_id, created_at FROM files WHERE id = ?", id)
	var file File
	if err := row.Scan(&file.ID, &file.Name, &file.Path, &file.ProjectID, &file.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &file, nil
}
