package storage

import (
	"database/sql"
	"embed"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(migrationsFS embed.FS) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "zentxt.db")
	if err != nil {
		return nil, err
	}

	migrations, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		runMigration, err := migrationsFS.ReadFile(filepath.Join("migrations", migration.Name()))
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
