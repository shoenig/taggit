package commands

import (
	"bytes"
	"testing"

	"github.com/shoenig/taggit/internal/cli"
	"github.com/shoenig/taggit/internal/cli/output"
	"github.com/shoenig/test/must"
)

type mocks struct {
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	writer     output.Writer
	tagLister  *cli.FakeTagLister
	tagCreator *cli.FakeTagCreator
	tagPusher  *cli.FakeTagPusher
}

func newMocks(_ *testing.T) mocks {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	return mocks{
		stdout:     &stdout,
		stderr:     &stderr,
		writer:     output.NewWriter(&stdout, &stderr),
		tagLister:  &cli.FakeTagLister{},
		tagCreator: &cli.FakeTagCreator{},
		tagPusher:  &cli.FakeTagPusher{},
	}
}

func Test_NewKit(t *testing.T) {
	mocks := newMocks(t)
	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPusher)

	must.NotNil(t, kit.writer)
	must.NotNil(t, kit.tagLister)
	must.NotNil(t, kit.tagCreator)
	must.NotNil(t, kit.tagPusher)
}
