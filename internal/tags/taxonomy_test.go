package tags

import (
	"testing"

	"github.com/shoenig/semantic"
	"github.com/shoenig/test/must"
)

const sampleTags = `
deploy/2017-03-17
deploy/2017-05-12
deploy/2017-10-12_09-05-04
deploy/2017-10-12_09-53-03
deploy/2017-10-12_10-42-32
v0.0.0
v0.0.1
v0.0.1-alpha
v0.0.1-alpha2
v0.0.5
v0.0.6
v1.0.0
v1.0.0-rc1
v1.1.0
v1.1.1
0.0.3
`

func Test_Taxonomy_Add(t *testing.T) {
	tax := Taxonomy{
		NewTriple(1, 2, 3): Tags{
			semantic.New(1, 2, 3),
		},
	}

	tax.Add(semantic.New(1, 3, 0))

	exp := Taxonomy{
		NewTriple(1, 2, 3): Tags{
			semantic.New(1, 2, 3),
		},
		NewTriple(1, 3, 0): Tags{
			semantic.New(1, 3, 0),
		},
	}

	must.MapEq(t, exp, tax)
}

func Test_Taxonomy_Sort(t *testing.T) {
	orig := Taxonomy{
		NewTriple(1, 2, 3): []semantic.Tag{
			semantic.New(1, 2, 3),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New2(1, 2, 3, "alpha1"),
		},
		NewTriple(1, 3, 0): []semantic.Tag{
			semantic.New(1, 3, 0),
			semantic.New3(1, 3, 0, "rc1", "bm1"),
		},
	}

	exp := Taxonomy{
		NewTriple(1, 2, 3): []semantic.Tag{
			semantic.New(1, 2, 3),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New2(1, 2, 3, "alpha1"),
		},
		NewTriple(1, 3, 0): []semantic.Tag{
			semantic.New(1, 3, 0),
			semantic.New3(1, 3, 0, "rc1", "bm1"),
		},
	}

	orig.Sort() // in place

	must.MapEq(t, exp, orig)
}

func Test_Taxonomy_Bases(t *testing.T) {
	orig := Taxonomy{
		NewTriple(1, 2, 3): []semantic.Tag{
			semantic.New(1, 2, 3),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New2(1, 2, 3, "alpha1"),
		},
		NewTriple(1, 3, 0): []semantic.Tag{
			semantic.New(1, 3, 0),
			semantic.New3(1, 3, 0, "rc1", "bm1"),
		},
	}

	bases := orig.Bases()
	must.Eq(t, []Triple{
		NewTriple(1, 2, 3),
		NewTriple(1, 3, 0),
	}, bases)
}

func Test_Taxonomy_Latest(t *testing.T) {
	orig := Taxonomy{
		NewTriple(1, 2, 3): []semantic.Tag{
			semantic.New(1, 2, 3),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New2(1, 2, 3, "alpha1"),
		},
		NewTriple(1, 3, 0): []semantic.Tag{
			semantic.New(1, 3, 0),
			semantic.New3(1, 3, 0, "rc1", "bm1"),
		},
	}

	latest := orig.Latest()
	must.Eq(t, semantic.New(1, 3, 0), latest)
}
