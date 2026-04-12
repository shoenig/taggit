package commands

import (
	"cattlecloud.net/go/babycli"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/tags"
)

func newZeroCommand(kit *Kit) *babycli.Component {
	return &babycli.Component{
		Name:        "zero",
		Help:        "Create initial v0.0.0 tag",
		Description: "Create initial v0.0.0 tag",
		Function:    zeroFunc(kit),
	}
}

func zeroFunc(kit *Kit) babycli.Func {
	return func(_ *babycli.Component) babycli.Code {
		writer := kit.writer
		tagLister := kit.tagLister
		tagCreator := kit.tagCreator
		tagPusher := kit.tagPusher

		writer.Tracef("create initial v0.0.0 tag")

		groups, err := tagLister.ListRepoTags()
		if err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		zero := semantic.New(0, 0, 0)

		if exists := tags.HasPrevious(groups); exists {
			writer.Errorf("refusing to generate zero tag (%s) when other semver tags already exist", zero)
			return babycli.Failure
		}

		if err := tagCreator.CreateTag(zero); err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		if err := tagPusher.PushTag(zero); err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		writer.Writef("created tag %s", zero)
		return babycli.Success
	}
}
