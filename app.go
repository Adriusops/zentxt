package zentxt

import (
	"database/sql"

	"github.com/Adriusops/zentxt/internal/versioning"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type App struct{ db *sql.DB }

func NewApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

func (a *App) ListFiles() ([]*versioning.File, error) {
	return versioning.ListFiles(a.db)
}

func (a *App) CreateFile(name, path string) (*versioning.File, error) {
	return versioning.CreateFile(a.db, name, path, nil)
}

func (a *App) GetFile(id string) (*versioning.File, error) {
	return versioning.GetFile(a.db, id)
}

func (a *App) SaveVersion(fileID, path, author, message, content string) (*versioning.Version, error) {
	return versioning.SaveVersion(a.db, fileID, path, author, message, content)
}

func (a *App) ListVersions(fileID string) ([]*versioning.Version, error) {
	return versioning.ListVersions(a.db, fileID)

}

func (a *App) GetVersion(fileID string) (*versioning.Version, error) {
	return versioning.GetVersion(a.db, fileID)

}

func (a *App) GenerateDiff(v1, v2 string) ([]diffmatchpatch.Diff, error) {
	version1, err := versioning.GetVersion(a.db, v1)
	if err != nil {
		return nil, err
	}
	version2, err := versioning.GetVersion(a.db, v2)
	if err != nil {
		return nil, err
	}
	return versioning.GenerateDiff(version1.Content, version2.Content), nil

}

func (a *App) RestoreVersion(fileID, versionID string) (*versioning.Version, error) {
	return versioning.RestoreVersion(a.db, fileID, versionID)

}

func (a *App) DeleteVersion(versionID string) error {
	return versioning.DeleteVersion(a.db, versionID)

}
