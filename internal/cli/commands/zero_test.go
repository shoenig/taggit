package commands

import (
	"errors"
	"strings"
	"testing"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/tags"
	"github.com/shoenig/test/must"
)

func Test_zeroFunc_normal(t *testing.T) {
	exp := "taggit: created tag v0.0.0\n"

	tk := newTestKit()

	zeroTag := semantic.New(0, 0, 0)

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = nil
	tk.tagCreator.Err = nil
	tk.tagPusher.Err = nil

	code := runZeroFunc(tk.kit())

	must.Eq(t, babycli.Success, code)
	must.Eq(t, exp, tk.stdout.String())
	must.Eq(t, "", tk.stderr.String())
	must.Eq(t, zeroTag, tk.tagCreator.Tag)
	must.Eq(t, zeroTag, tk.tagPusher.Tag)
}

func Test_zeroFunc_hasPrevious(t *testing.T) {
	exp := "refusing to generate zero tag (v0.0.0) when other semver tags already exist"

	tk := newTestKit()

	oldTag := semantic.New(1, 2, 3)
	tk.tagLister.Taxonomy = tags.Taxonomy{
		tags.NewTriple(1, 2, 3): {oldTag},
	}
	tk.tagLister.Err = nil

	code := runZeroFunc(tk.kit())

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.True(t, strings.Contains(tk.stderr.String(), exp))
}

func Test_zeroFunc_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = errors.New("some git error")

	code := runZeroFunc(tk.kit())

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}

func Test_zeroFunc_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = nil
	tk.tagCreator.Err = errors.New("some create error")
	tk.tagPusher.Err = nil

	code := runZeroFunc(tk.kit())

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}

func runZeroFunc(kit *Kit) babycli.Code {
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
