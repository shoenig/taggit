package cli

import (
	"github.com/shoenig/semantic"
	"github.com/shoenig/taggit/internal/tags"
)

type FakeTagLister struct {
	Taxonomy tags.Taxonomy
	Err      error
}

func (f *FakeTagLister) ListRepoTags() (tags.Taxonomy, error) {
	return f.Taxonomy, f.Err
}

type FakeTagCreator struct {
	Tag semantic.Tag
	Err error
}

func (f *FakeTagCreator) CreateTag(tag semantic.Tag) error {
	f.Tag = tag
	return f.Err
}

type FakeTagPusher struct {
	Tag semantic.Tag
	Err error
}

func (f *FakeTagPusher) PushTag(tag semantic.Tag) error {
	f.Tag = tag
	return f.Err
}
