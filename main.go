package main

import (
	"os"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/taggit/internal/cli/commands"
)

func main() {
	args := babycli.Arguments()
	rc := commands.Invoke(args)
	os.Exit(rc)
}
