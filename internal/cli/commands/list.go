package commands

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/cli"
	"github.com/shoenig/taggit/internal/cli/output"
	"github.com/shoenig/taggit/internal/tags"
)

const (
	listCmdName     = "list"
	listCmdSynopsis = "List tagged versions."
	listCmdUsage    = "list"
)

func NewListCmd(kit *Kit) subcommands.Command {
	return &listCmd{
		writer:    kit.writer,
		tagLister: kit.tagLister,
	}
}

type listCmd struct {
	writer    output.Writer
	tagLister cli.TagLister
}

func (lc *listCmd) Name() string {
	return listCmdName
}

func (lc *listCmd) Synopsis() string {
	return listCmdSynopsis
}

func (lc *listCmd) Usage() string {
	return listCmdUsage
}

func (lc *listCmd) SetFlags(_ *flag.FlagSet) {
}

func (lc *listCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := lc.execute(); err != nil {
		lc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (lc *listCmd) execute() error {
	tax, err := lc.tagLister.ListRepoTags()
	if err != nil {
		return err
	}

	lc.list(tax)
	return nil
}

func (lc *listCmd) list(groups tags.Taxonomy) {
	lc.writer.Tracef("listing tags in git repository")

	triples := groups.Bases()

	for _, triple := range triples {
		tagsOfTriple := groups[triple]
		line := outputLineForTriple(triple, tagsOfTriple)
		lc.writer.Directf("%s", line)
	}
}

func outputLineForTriple(triple tags.Triple, associated []semantic.Tag) string {
	asString := associatedList(associated)
	s := fmt.Sprintf("%s |= %s", triple, strings.Join(asString, " "))
	return s
}

func associatedList(associated []semantic.Tag) []string {
	asStrings := make([]string, 0, len(associated))
	for _, aTag := range associated {
		asStrings = append(asStrings, aTag.String())
	}
	return asStrings
}
