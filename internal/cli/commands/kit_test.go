package commands

import (
	"bytes"
	"testing"

	"github.com/shoenig/test"
	"gophers.dev/cmds/taggit/internal/cli"
	"gophers.dev/cmds/taggit/internal/cli/output"
)

type mocks struct {
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	writer     output.Writer
	tagLister  *cli.TagListerMock
	tagCreator *cli.TagCreatorMock
	tagPusher  *cli.TagPusherMock
}

func newMocks(t *testing.T) mocks {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	return mocks{
		stdout:     &stdout,
		stderr:     &stderr,
		writer:     output.NewWriter(&stdout, &stderr),
		tagLister:  cli.NewTagListerMock(t),
		tagCreator: cli.NewTagCreatorMock(t),
		tagPusher:  cli.NewTagPusherMock(t),
	}
}

func (m mocks) assertions(t *testing.T) {
	m.tagLister.MinimockFinish()
	m.tagCreator.MinimockFinish()
	m.tagPusher.MinimockFinish()
}

func Test_NewKit(t *testing.T) {
	mocks := newMocks(t)
	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPusher)

	test.NotNil(t, kit.writer)
	test.NotNil(t, kit.tagLister)
	test.NotNil(t, kit.tagCreator)
	test.NotNil(t, kit.tagPusher)
}
