package versioning

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GenerateDiff(content1 string, content2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	return dmp.DiffMain(content1, content2, false)
}
