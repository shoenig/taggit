package commands

import (
	"cattlecloud.net/go/babycli"
	"github.com/shoenig/taggit/internal/tags"
)

func newPatchCommand(kit *Kit) *babycli.Component {
	return &babycli.Component{
		Name:        "patch",
		Help:        "Create an incremented patch version",
		Description: "Create an incremented patch version",
		Flags: babycli.Flags{
			{
				Type:  babycli.StringFlag,
				Long:  "meta",
				Short: "m",
				Help:  "include build +metadata label",
			},
			{
				Type:  babycli.BooleanFlag,
				Long:  "dry",
				Short: "y",
				Help:  "dry-run without making changes",
			},
		},
		Function: patchFunc(kit),
	}
}

func patchFunc(kit *Kit) babycli.Func {
	return func(cmd *babycli.Component) babycli.Code {
		writer := kit.writer
		tagLister := kit.tagLister
		tagCreator := kit.tagCreator
		tagPusher := kit.tagPusher

		meta := cmd.GetString("meta")

		ext := tags.ExtractExtensions(meta, cmd.Arguments())
		writer.Tracef(
			"increment patch version, pre-release: %q, build-metadata: %q",
			ext.PreRelease, ext.BuildMetaData,
		)

		tax, err := tagLister.ListRepoTags()
		if err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		if exists := tags.HasPrevious(tax); !exists {
			writer.Errorf("cannot increment tag because no previous tags exist")
			return babycli.Failure
		}

		latest := tax.Latest()
		next := tags.IncPatch(latest, ext)

		if cmd.GetBool("dry") {
			writer.Writef("would create tag %s", next)
			return babycli.Success
		}

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
