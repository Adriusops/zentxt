package versioning

import (
	"database/sql"
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

func CreateFile(db *sql.DB, name string, path string, projectID *string) (*File, error) {
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
