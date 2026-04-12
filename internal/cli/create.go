package cli

import (
	"fmt"

	git5 "github.com/go-git/go-git/v5"
	"github.com/shoenig/semantic"
)

// A TagCreator is used to create git tags.
type TagCreator interface {
	// CreateTag creates the given tag in the repository.
	CreateTag(semantic.Tag) error
}

type tagCreator struct {
	repository *git5.Repository
}

func NewTagCreator(r *git5.Repository) TagCreator {
	return &tagCreator{
		repository: r,
	}
}

func (tc *tagCreator) CreateTag(tag semantic.Tag) error {
	head, err := tc.repository.Head()
	if err != nil {
		return fmt.Errorf("could not create tag: %w", err)
	}

	if _, err := tc.repository.CreateTag(
		tag.String(),
		head.Hash(),
		nil, // options, maybe include a message
	); err != nil {
		return fmt.Errorf("could not create tag: %w", err)
	}

	return nil
}
