package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	git5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/shoenig/test"
)

var signature = &object.Signature{
	Name:  "Testing",
	Email: "testing@example.com",
	When:  time.Date(2020, 11, 8, 14, 18, 0, 0, time.UTC),
}

func CleanupT(t *testing.T, r *git5.Repository) {
	w, err := r.Worktree()
	test.NoError(t, err)
	root := w.Filesystem.Root()

	err = os.RemoveAll(root)
	test.NoError(t, err)
}

func CreateT(t *testing.T, tags []string) *git5.Repository {
	dir, err := ioutil.TempDir("", "taggit-")
	test.NoError(t, err)

	r, err := git5.PlainInit(dir, false)
	test.NoError(t, err)

	w, err := r.Worktree()
	test.NoError(t, err)

	for i, tag := range tags {
		msg := fmt.Sprintf("commit #%d", i)
		h, err := w.Commit(msg, &git5.CommitOptions{
			Author: signature,
		})
		test.NoError(t, err)
		_, err = r.CreateTag(tag, h, nil)
		test.NoError(t, err)
	}

	return r
}
