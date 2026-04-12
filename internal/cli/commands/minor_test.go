package commands

import (
	"errors"
	"testing"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/tags"
	"github.com/shoenig/test/must"
)

func Test_minorFunc_normal(t *testing.T) {
	exp := "taggit: created tag v1.3.0\n"

	tk := newTestKit()

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(1, 3, 0)

	tk.tagLister.Taxonomy = tags.Taxonomy{
		tags.NewTriple(1, 2, 3): {oldTag},
	}
	tk.tagLister.Err = nil
	tk.tagCreator.Err = nil
	tk.tagPusher.Err = nil

	code := runMinorFunc(tk.kit(), "")

	must.Eq(t, babycli.Success, code)
	must.Eq(t, exp, tk.stdout.String())
	must.Eq(t, "", tk.stderr.String())
	must.Eq(t, newTag, tk.tagCreator.Tag)
	must.Eq(t, newTag, tk.tagPusher.Tag)
}

func Test_minorFunc_noPrevious(t *testing.T) {
	exp := "taggit: cannot increment tag because no previous tags exist\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = nil

	code := runMinorFunc(tk.kit(), "")

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}

func Test_minorFunc_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = errors.New("some git error")

	code := runMinorFunc(tk.kit(), "")

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}

func Test_minorFunc_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = tags.Taxonomy{
		tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
	}
	tk.tagLister.Err = nil
	tk.tagCreator.Err = errors.New("some create error")
	tk.tagPusher.Err = nil

	code := runMinorFunc(tk.kit(), "")

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}

func runMinorFunc(kit *Kit, meta string) babycli.Code {
	writer := kit.writer
	tagLister := kit.tagLister
	tagCreator := kit.tagCreator
	tagPusher := kit.tagPusher

	ext := tags.ExtractExtensions(meta, nil)
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
