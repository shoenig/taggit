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
	writer := output.NewWriter(os.Stdout, os.Stderr)

	repository, err := git5.PlainOpen(".")
	if err != nil {
		writer.Errorf("unable to open repository: %v", err)
		return babycli.Failure
	}

	tagLister := cli.NewTagLister(repository)
	tagCreator := cli.NewTagCreator(repository)
	tagPusher := cli.NewTagPusher(repository)
	kit := NewKit(
		writer,
		tagLister,
		tagCreator,
		tagPusher,
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
			Function: func(cmd *babycli.Component) babycli.Code {
				if cmd.GetBool("version") {
					fmt.Println(version.Version)
					return babycli.Success
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
