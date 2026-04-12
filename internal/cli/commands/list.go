package commands

import (
	"fmt"
	"strings"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/cli/output"
	"github.com/shoenig/taggit/internal/tags"
)

func newListCommand(kit *Kit) *babycli.Component {
	return &babycli.Component{
		Name:        "list",
		Help:        "List tagged versions.",
		Description: "List tagged versions.",
		Function:    listFunc(kit),
	}
}

func listFunc(kit *Kit) babycli.Func {
	return func(_ *babycli.Component) babycli.Code {
		writer := kit.writer
		tagLister := kit.tagLister

		tax, err := tagLister.ListRepoTags()
		if err != nil {
			writer.Errorf("failure: %v", err)
			return babycli.Failure
		}

		list(tax, writer)
		return babycli.Success
	}
}

func list(groups tags.Taxonomy, writer output.Writer) {
	writer.Tracef("listing tags in git repository")

	triples := groups.Bases()

	for _, triple := range triples {
		tagsOfTriple := groups[triple]
		line := outputLineForTriple(triple, tagsOfTriple)
		writer.Directf("%s", line)
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
