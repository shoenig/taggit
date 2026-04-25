package commands

import (
	"fmt"
	"os"

	"cattlecloud.net/go/babycli"
	git5 "github.com/go-git/go-git/v5"
	"github.com/shoenig/taggit/internal/cli"
	"github.com/shoenig/taggit/internal/cli/output"
	"github.com/shoenig/taggit/version"
)

func Invoke(args []string) babycli.Code {
	repository, err := git5.PlainOpen(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open repository: %v", err)
		return babycli.Failure
	}

	kit := NewKit(
		output.NewWriter(os.Stdout, os.Stderr),
		cli.NewTagLister(repository),
		cli.NewTagCreator(repository),
		cli.NewTagPusher(repository),
	)

	return babycli.New(&babycli.Configuration{
		Arguments: args,
		Version:   version.Version,
		Globals: babycli.Flags{
			{
				Type:  babycli.BooleanFlag,
				Long:  "version",
				Short: "v",
				Help:  "print version information",
			},
		},
		Top: &babycli.Component{
			Name:        "taggit",
			Description: "Publish new versions of Go modules.",
			Function: func(c *babycli.Component) babycli.Code {
				if c.GetBool("version") {
					fmt.Println(version.Version)
				}
				return babycli.Success
			},
			Components: babycli.Components{
				newListCommand(kit),
				newZeroCommand(kit),
				newPatchCommand(kit),
				newMinorCommand(kit),
				newMajorCommand(kit),
			},
		},
	}).Run()
}
