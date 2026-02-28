package versioning

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GenerateDiff(content1 string, content2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()

	// Convert to line-based diff
	a, b, lineArray := dmp.DiffLinesToChars(content1, content2)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)
	diffs = dmp.DiffCleanupSemantic(diffs) // Clean up for readability

	return diffs
}
