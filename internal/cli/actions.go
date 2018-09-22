package cli

import (
	"fmt"
	"strings"

	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/tags"
)

func List(tags []tags.Tag) {
	var b strings.Builder
	for _, tag := range tags {
		b.WriteString(tag.String())
		b.WriteString("\n")
	}
	fmt.Print(b.String())
}

func Zero(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) > 0 {
		return fmt.Errorf("refusing to generate zero tag (v0.0.0) when other semver tags already exist")
	}

	return git.CreateTag(cmd, tags.ZeroValue)
}

func Patch(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump patch with no pre-existing tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor,
		Patch: latest.Patch + 1,
	})
}

func Minor(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump minor with no pre-existing tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor + 1,
		Patch: 0,
	})
}

func Major(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump major with no pre-existnig tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major + 1,
		Minor: 0,
		Patch: 0,
	})
}
