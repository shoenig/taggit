package output

import (
	"bytes"
	"testing"

	"github.com/shoenig/test/must"
)

func testWriters() (*bytes.Buffer, *bytes.Buffer) {
	var (
		normal  bytes.Buffer
		failure bytes.Buffer
	)
	return &normal, &failure
}

func Test_Writer_Directf(t *testing.T) {
	normal, failure := testWriters()

	w := NewWriter(normal, failure)
	w.Directf("foo %s %d", "bar", 1)

	expN := "foo bar 1\n"
	must.Eq(t, expN, normal.String())
	must.Eq(t, expN, normal.String())
	must.Eq(t, "", failure.String())
}

func Test_Writer_Writef(t *testing.T) {
	normal, failure := testWriters()

	w := NewWriter(normal, failure)
	w.Writef("foo %s %d", "bar", 1)

	expN := "taggit: foo bar 1\n"
	must.Eq(t, expN, normal.String())
	must.Eq(t, "", failure.String())
}

func Test_Writer_Errorf(t *testing.T) {
	normal, failure := testWriters()

	w := NewWriter(normal, failure)
	w.Errorf("foo %s %d", "bar", 1)

	expN := "taggit: foo bar 1\n"
	must.Eq(t, "", normal.String())
	must.Eq(t, expN, failure.String())
}

func Test_Writer_Tracef(t *testing.T) {
	t.Setenv(tracesEnv, "1")
	normal, failure := testWriters()

	w := NewWriter(normal, failure)
	w.Tracef("foo %s %d", "bar", 1)

	expN := "trace: foo bar 1\n"
	must.Eq(t, expN, normal.String())
	must.Eq(t, "", failure.String())
}
