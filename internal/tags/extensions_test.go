package tags

import (
	"flag"
	"testing"

	"github.com/shoenig/test"
)

func Test_ExtractExtensions(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	fs.String("meta", "bm1", "set metadata")

	err := fs.Set("meta", "abc123")
	test.NoError(t, err)

	ext := ExtractExtensions(fs)
	test.Eq(t, "abc123", ext.BuildMetaData)
}
