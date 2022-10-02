package tags

import (
	"flag"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_ExtractExtensions(t *testing.T) {
	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	fs.String("meta", "bm1", "set metadata")

	err := fs.Set("meta", "abc123")
	must.NoError(t, err)

	ext := ExtractExtensions(fs)
	must.Eq(t, "abc123", ext.BuildMetaData)
}
