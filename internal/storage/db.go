package storage

import (
	"database/sql"
	"embed"
	"os"
	"path"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(migrationsFS embed.FS) (*sql.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "zentxt")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "zentxt.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	migrations, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		runMigration, err := migrationsFS.ReadFile(path.Join("migrations", migration.Name()))
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(string(runMigration))
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
