// Command taggit publishes new versions of Go modules.
package main

import (
	"context"
	"flag"
	"os"

	git5 "github.com/go-git/go-git/v5"
	"github.com/google/subcommands"

	"gophers.dev/cmds/taggit/internal/cli"
	"gophers.dev/cmds/taggit/internal/cli/commands"
	"gophers.dev/cmds/taggit/internal/cli/output"
)

func main() {
	writer := output.NewWriter(os.Stdout, os.Stderr)

	repository, err := git5.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	tagLister := cli.NewTagLister(repository)
	tagCreator := cli.NewTagCreator(repository)
	tagPusher := cli.NewTagPusher(repository)
	kit := commands.NewKit(
		writer,
		tagLister,
		tagCreator,
		tagPusher,
	)

	listCmd := commands.NewListCmd(kit)
	zeroCmd := commands.NewZeroCmd(kit)
	patchCmd := commands.NewPatchCmd(kit)
	minorCmd := commands.NewMinorCmd(kit)
	majorCmd := commands.NewMajorCmd(kit)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	subs := subcommands.NewCommander(fs, "")
	subs.Register(subs.HelpCommand(), "")
	subs.Register(subs.FlagsCommand(), "")
	subs.Register(listCmd, "")
	subs.Register(zeroCmd, "")
	subs.Register(patchCmd, "")
	subs.Register(minorCmd, "")
	subs.Register(majorCmd, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	ctx := context.Background()
	rc := subs.Execute(ctx, fs.Args())
	os.Exit(int(rc))
}
