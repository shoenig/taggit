package commands

import (
	"cattlecloud.net/go/babycli"
	"github.com/shoenig/taggit/internal/tags"
)

func newMinorCommand(kit *Kit) *babycli.Component {
	return &babycli.Component{
		Name:        "minor",
		Help:        "Create an incremented minor version",
		Description: "Create an incremented minor version",
		Flags: babycli.Flags{
			{
				Type:  babycli.StringFlag,
				Long:  "meta",
				Short: "m",
				Help:  "build metadata label",
			},
		},
		Function: minorFunc(kit),
	}
}

func minorFunc(kit *Kit) babycli.Func {
	return func(cmd *babycli.Component) babycli.Code {
		writer := kit.writer
		tagLister := kit.tagLister
		tagCreator := kit.tagCreator
		tagPusher := kit.tagPusher

		meta := cmd.GetString("meta")

		ext := tags.ExtractExtensions(meta, cmd.Arguments())
		writer.Tracef(
			"increment minor version, pre-release: %q, build-metadata: %q",
			ext.PreRelease, ext.BuildMetaData,
		)

		groups, err := tagLister.ListRepoTags()
		if err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		if exists := tags.HasPrevious(groups); !exists {
			writer.Errorf("cannot increment tag because no previous tags exist")
			return babycli.Failure
		}

		latest := groups.Latest()
		next := tags.IncMinor(latest, ext)

		if err := tagCreator.CreateTag(next); err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		if err := tagPusher.PushTag(next); err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		writer.Writef("created tag %s", next)
		return babycli.Success
	}
}
