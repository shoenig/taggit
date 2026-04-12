package commands

import (
	"bytes"
	"errors"
	"testing"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/cli"
	"github.com/shoenig/taggit/internal/cli/output"
	"github.com/shoenig/taggit/internal/tags"
	"github.com/shoenig/test/must"
)

type testKit struct {
	writer     output.Writer
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	tagLister  *cli.FakeTagLister
	tagCreator *cli.FakeTagCreator
	tagPusher  *cli.FakeTagPusher
}

func newTestKit() *testKit {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	writer := output.NewWriter(stdout, stderr)
	return &testKit{
		writer:     writer,
		stdout:     stdout,
		stderr:     stderr,
		tagLister:  &cli.FakeTagLister{},
		tagCreator: &cli.FakeTagCreator{},
		tagPusher:  &cli.FakeTagPusher{},
	}
}

func (t *testKit) kit() *Kit {
	return NewKit(t.writer, t.tagLister, t.tagCreator, t.tagPusher)
}

func Test_listFunc_noTags(t *testing.T) {
	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = nil

	cmd := newListCommand(tk.kit())
	code := cmd.Function(cmd)

	must.Eq(t, babycli.Success, code)
	must.Eq(t, "", tk.stdout.String())
}

func Test_listFunc_someTags(t *testing.T) {
	exp := `v0.1.0 |= v0.1.0 v0.1.0-alpha1
v0.2.0 |= v0.2.0-rc1 v0.2.0-r1+linux v0.2.0-r1+darwin
`
	tk := newTestKit()

	tk.tagLister.Taxonomy = tags.Taxonomy{
		tags.NewTriple(0, 1, 0): []semantic.Tag{
			semantic.New(0, 1, 0),
			semantic.New2(0, 1, 0, "alpha1"),
		},
		tags.NewTriple(0, 2, 0): []semantic.Tag{
			semantic.New2(0, 2, 0, "rc1"),
			semantic.New3(0, 2, 0, "r1", "linux"),
			semantic.New3(0, 2, 0, "r1", "darwin"),
		},
	}
	tk.tagLister.Err = nil

	cmd := newListCommand(tk.kit())
	code := cmd.Function(cmd)

	must.Eq(t, babycli.Success, code)
	must.Eq(t, exp, tk.stdout.String())
	must.Eq(t, "", tk.stderr.String())
}

func Test_listFunc_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	tk := newTestKit()

	tk.tagLister.Taxonomy = nil
	tk.tagLister.Err = errors.New("some git error")

	cmd := newListCommand(tk.kit())
	code := cmd.Function(cmd)

	must.Eq(t, babycli.Failure, code)
	must.Eq(t, "", tk.stdout.String())
	must.Eq(t, exp, tk.stderr.String())
}
