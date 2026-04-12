package tags

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_ExtractExtensions(t *testing.T) {
	ext := ExtractExtensions("abc123", nil)
	must.Eq(t, "abc123", ext.BuildMetaData)
}

func Test_ExtractExtensions_withPreRelease(t *testing.T) {
	ext := ExtractExtensions("abc123", []string{"beta"})
	must.Eq(t, "beta", ext.PreRelease)
	must.Eq(t, "abc123", ext.BuildMetaData)
}
