package commands

import (
	"github.com/shoenig/taggit/internal/cli"
	"github.com/shoenig/taggit/internal/cli/output"
)

func NewKit(
	writer output.Writer,
	tagLister cli.TagLister,
	tagCreator cli.TagCreator,
	tagPusher cli.TagPusher,
) *Kit {
	return &Kit{
		writer:     writer,
		tagLister:  tagLister,
		tagCreator: tagCreator,
		tagPusher:  tagPusher,
	}
}

// A Kit contains all the things needed for creating and publishing a new tag.
type Kit struct {
	writer     output.Writer
	tagLister  cli.TagLister
	tagCreator cli.TagCreator
	tagPusher  cli.TagPusher
}
