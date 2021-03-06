package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Cmd -s _mock.go

// An Cmd executes git commands.
type Cmd interface {
	Run(args []string, timeout time.Duration) (string, error)
}

// New creates an Cmd backed by the given executable.
func New(executable string) Cmd {
	return &cmd{
		executable: executable,
	}
}

type cmd struct {
	executable string
}

func (c *cmd) Run(args []string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.executable, args...)
	cmd.Env = os.Environ() // use the tty's environment
	bs, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "subprocess command failed: %v\n", err)
		_, _ = fmt.Fprintf(os.Stderr, "output: %s\n", string(bs))
		return "", errors.Wrap(err, "subprocess exec failed")
	}
	return string(bs), nil
}
