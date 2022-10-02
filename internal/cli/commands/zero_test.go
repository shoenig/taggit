package commands

import (
	"context"
	"errors"
	"flag"
	"strings"
	"testing"

	"github.com/google/subcommands"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/tags"
	"github.com/shoenig/test/must"
)

func Test_ZeroCmd_commandInfo(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions(t)

	majorCmd := NewZeroCmd(newKit(mocks))

	name := majorCmd.Name()
	must.Eq(t, zeroCmdName, name)

	synop := majorCmd.Synopsis()
	must.Eq(t, zeroCmdSynopsis, synop)

	usage := majorCmd.Usage()
	must.Eq(t, zeroCmdUsage, usage)
}

func Test_ZeroCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v0.0.0\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	zeroTag := semantic.New(0, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil, // no tags
	)
	mocks.tagCreator.CreateTagMock.When(zeroTag).Then(nil)
	mocks.tagPusher.PushTagMock.When(zeroTag).Then(nil)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	must.Eq(t, subcommands.ExitSuccess, status)
	must.Eq(t, exp, mocks.stdout.String())
	must.Eq(t, "", mocks.stderr.String())
}

func Test_ZeroCmd_Execute_hasPrevious(t *testing.T) {
	exp := "refusing to generate zero tag (v0.0.0) when other semver tags already exist\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	// zeroTag := semantic.New(0, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	// mocks.tagCreator.CreateTagMock.When(zeroTag).Then(nil)
	// mocks.tagPublisher.PublishMock.When(zeroTag).Then(nil)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	must.Eq(t, subcommands.ExitFailure, status)
	must.Eq(t, "", mocks.stdout.String())
	must.True(t, strings.Contains(mocks.stderr.String(), exp))
}

func Test_ZeroCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	must.Eq(t, subcommands.ExitFailure, status)
	must.Eq(t, "", mocks.stdout.String())
	must.Eq(t, exp, mocks.stderr.String())
}

func Test_ZeroCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(0, 0, 0),
	).Return(
		errors.New("some create error"),
	)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	must.Eq(t, subcommands.ExitFailure, status)
	must.Eq(t, "", mocks.stdout.String())
	must.Eq(t, exp, mocks.stderr.String())
}
