package tags

import (
	"strings"
)

type Extensions struct {
	PreRelease    string
	BuildMetaData string
}

func ExtractExtensions(meta string, preReleaseArgs []string) Extensions {
	var (
		preRelease    string
		buildMetadata string
	)

	if len(preReleaseArgs) > 0 {
		preRelease = preReleaseArgs[0]
	}

	buildMetadata = meta

	return Extensions{
		PreRelease:    clean(preRelease),
		BuildMetaData: clean(buildMetadata),
	}
}

func clean(orig string) string {
	noDash := strings.TrimPrefix(orig, "-")
	noPlus := strings.TrimPrefix(noDash, "+")
	return noPlus
}
