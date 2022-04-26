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

func newKit(mocks mocks) *Kit {
	return NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPusher)
}

func Test_ListCmd_commandInfo(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	listCmd := NewListCmd(newKit(mocks))

	name := listCmd.Name()
	test.Eq(t, listCmdName, name)

	synop := listCmd.Synopsis()
	test.Eq(t, listCmdSynopsis, synop)

	usage := listCmd.Usage()
	test.Eq(t, listCmdUsage, usage)
}

func Test_ListCmd_Execute_noTags(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), // no tags to parse
		nil,
	)

	listCmd := NewListCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	test.Eq(t, subcommands.ExitSuccess, status)
}

func Test_ListCmd_Execute_someTags(t *testing.T) {
	exp := `v0.1.0 |= v0.1.0 v0.1.0-alpha1
v0.2.0 |= v0.2.0-rc1 v0.2.0-r1+linux v0.2.0-r1+darwin
`
	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(0, 1, 0): []semantic.Tag{
				semantic.New(0, 1, 0),
				semantic.New2(0, 1, 0, "alpha1"),
			},
			tags.NewTriple(0, 2, 0): []semantic.Tag{
				semantic.New2(0, 2, 0, "rc1"),
				semantic.New3(0, 2, 0, "r1", "linux"),
				semantic.New3(0, 2, 0, "r1", "darwin"),
			},
		},
		nil,
	)

	listCmd := NewListCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	test.Eq(t, subcommands.ExitSuccess, status)
	test.Eq(t, exp, mocks.stdout.String())
	test.Eq(t, "", mocks.stderr.String())
}

func Test_ListCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, errors.New("some git error"),
	)

	listCmd := NewListCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	test.Eq(t, subcommands.ExitFailure, status)
	test.Eq(t, "", mocks.stdout.String())
	test.Eq(t, exp, mocks.stderr.String())
}
