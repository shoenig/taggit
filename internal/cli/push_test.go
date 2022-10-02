package cli

import (
	"testing"

	"github.com/shoenig/semantic"
	"github.com/shoenig/test/must"
)

func Test_TagPusher_PushTag(t *testing.T) {
	repo := CreateT(t, []string{
		"v0.0.1",
		"v0.0.2",
	})
	defer CleanupT(t, repo)

	pusher := NewTagPusher(repo)
	err := pusher.PushTag(semantic.New(0, 0, 3))
	must.NoError(t, err)
}
