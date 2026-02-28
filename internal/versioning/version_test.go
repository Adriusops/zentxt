package versioning

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

func TestSaveVersion(t *testing.T) {
	db := setupTestDB(t)

	file, _ := CreateFile(db, "test.txt", "/tmp/test.txt", nil)

	version, err := SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Initial commit", "yes")
	assert.NoError(t, err)
	assert.Equal(t, file.ID, version.FileID)
	assert.NotEmpty(t, version.ID)
}
