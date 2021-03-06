package commands

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/stretchr/testify/require"
	"gophers.dev/pkgs/semantic"
	"oss.indeed.com/go/taggit/internal/tags"
)

func Test_ZeroCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)

	majorCmd := NewZeroCmd(kit)

	name := majorCmd.Name()
	r.Equal(zeroCmdName, name)

	synop := majorCmd.Synopsis()
	r.Equal(zeroCmdSynopsis, synop)

	usage := majorCmd.Usage()
	r.Equal(zeroCmdUsage, usage)
}

func Test_ZeroCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v0.0.0\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	zeroTag := semantic.New(0, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil, // no tags
	)
	mocks.tagCreator.CreateTagMock.When(zeroTag).Then(nil)
	mocks.tagPublisher.PublishMock.When(zeroTag).Then(nil)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	zeroCmd := NewZeroCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_ZeroCmd_Execute_hasPrevious(t *testing.T) {
	exp := "refusing to generate zero tag (v0.0.0) when other semver tags already exist\n"

	r := require.New(t)

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

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	zeroCmd := NewZeroCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Contains(mocks.stderr.String(), exp)
}

func Test_ZeroCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	zeroCmd := NewZeroCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_ZeroCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	r := require.New(t)

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

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	zeroCmd := NewZeroCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_ZeroCmd_Execute_publishErr(t *testing.T) {
	exp := "taggit: failure: some publish error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(0, 0, 0),
	).Return(nil)

	mocks.tagPublisher.PublishMock.Expect(
		semantic.New(0, 0, 0),
	).Return(errors.New("some publish error"))

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	zeroCmd := NewZeroCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
