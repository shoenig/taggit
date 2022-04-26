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

func Test_PatchCmd_commandInfo(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	patchCmd := NewPatchCmd(newKit(mocks))

	name := patchCmd.Name()
	test.Eq(t, patchCmdName, name)

	synop := patchCmd.Synopsis()
	test.Eq(t, patchCmdSynopsis, synop)

	usage := patchCmd.Usage()
	test.Eq(t, patchCmdUsage, usage)
}

func Test_PatchCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v1.2.4\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(1, 2, 4)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	mocks.tagCreator.CreateTagMock.When(newTag).Then(nil)
	mocks.tagPusher.PushTagMock.When(newTag).Then(nil)
	patchCmd := NewPatchCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)

	test.Eq(t, subcommands.ExitSuccess, status)
	test.Eq(t, exp, mocks.stdout.String())
	test.Eq(t, "", mocks.stderr.String())
}

func Test_PatchCmd_Execute_noPrevious(t *testing.T) {
	exp := `taggit: cannot increment tag because no previous tags exist
taggit: failure: no previous tags
`

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), nil, // no tags, no error
	)

	patchCmd := NewPatchCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}

func Test_PatchCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	patchCmd := NewPatchCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}

func Test_PatchCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
		}, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(1, 2, 4),
	).Return(
		errors.New("some create error"),
	)

	patchCmd := NewPatchCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}
