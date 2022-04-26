package commands

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/shoenig/test"
	"gophers.dev/cmds/taggit/internal/tags"
	"gophers.dev/pkgs/semantic"
)

func Test_MajorCmd_commandInfo(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	majorCmd := NewMajorCmd(newKit(mocks))

	name := majorCmd.Name()
	test.Eq(t, majorCmdName, name)

	synop := majorCmd.Synopsis()
	test.Eq(t, majorCmdSynopsis, synop)

	usage := majorCmd.Usage()
	test.Eq(t, majorCmdUsage, usage)
}

func Test_MajorCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v2.0.0\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(2, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	mocks.tagCreator.CreateTagMock.When(newTag).Then(nil)
	mocks.tagPusher.PushTagMock.When(newTag).Then(nil)

	majorCmd := NewMajorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)

	test.Eq(t, subcommands.ExitSuccess, status)
	test.Eq(t, exp, mocks.stdout.String())
	test.Eq(t, "", mocks.stderr.String())
}

func Test_MajorCmd_Execute_noPrevious(t *testing.T) {
	exp := `taggit: cannot increment tag because no previous tags exist
taggit: failure: no previous tags
`

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), nil, // no tags, no error
	)

	majorCmd := NewMajorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}

func Test_MajorCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	majorCmd := NewMajorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}

func Test_MajorCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
		}, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(2, 0, 0),
	).Return(
		errors.New("some create error"),
	)

	majorCmd := NewMajorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}
