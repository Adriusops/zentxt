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

func TestSaveVersion_VersionNumber(t *testing.T) {
	db := setupTestDB(t)

	file, _ := CreateFile(db, "test.txt", "/tmp/test.txt", nil)

	version1, err := SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Initial commit", "yes")
	assert.NoError(t, err)
	version2, err := SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Second commit", "yes")
	assert.NoError(t, err)
	version3, err := SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Third commit", "yes")
	assert.NoError(t, err)

	assert.Equal(t, 1, version1.VersionNumber)
	assert.Equal(t, 2, version2.VersionNumber)
	assert.Equal(t, 3, version3.VersionNumber)
}

func TestListVersions(t *testing.T) {
	db := setupTestDB(t)

	file, _ := CreateFile(db, "test.txt", "/tmp/test.txt", nil)

	_, err := SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Initial commit", "yes")
	assert.NoError(t, err)
	_, err = SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Second commit", "yes")
	assert.NoError(t, err)
	_, err = SaveVersion(db, file.ID, "/tmp/test.txt", "tonton", "Third commit", "yes")
	assert.NoError(t, err)

	versions, err := ListVersions(db, file.ID)
	assert.Equal(t, 3, len(versions))

}
