package cli

import (
	"testing"

	"github.com/shoenig/semantic"
	"github.com/shoenig/test/must"
)

func Test_TagCreator_CreateTag(t *testing.T) {
	repo := CreateT(t, []string{
		"v0.0.1",
		"v0.0.2",
	})
	defer CleanupT(t, repo)

	creator := NewTagCreator(repo)
	err := creator.CreateTag(semantic.Tag{
		Major: 0,
		Minor: 0,
		Patch: 3,
	})
	must.NoError(t, err)
}
