package versioning

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDiff(t *testing.T) {
	// Cas 1 : contenus identiques
	diffs := GenerateDiff("hello", "hello")
	assert.Len(t, diffs, 1) // un seul élément DiffEqual

	// Cas 2 : contenus différents
	diffs = GenerateDiff("hello", "world")
	assert.Greater(t, len(diffs), 1) // plusieurs éléments
}
