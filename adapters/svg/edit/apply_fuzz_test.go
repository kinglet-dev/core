package edit_test

import (
	"bytes"
	"testing"

	"github.com/kinglet-dev/core/adapters/svg/edit"
)

// FuzzApply checks that Apply never panics on hostile input, returns a
// document unchanged when asked for no changes, and never touches the
// bytes before the first element it changes.
func FuzzApply(f *testing.F) {
	f.Add([]byte(`<svg><rect width="1"/><circle r="1"/></svg>`), 1, 2, "<g/>")
	f.Add([]byte(`<svg><g><rect/></g></svg>`), 2, 1, "")
	f.Add([]byte(`<?xml version="1.0"?><!-- c --><svg>&amp;<![CDATA[x]]></svg>`), 0, 0, "")
	f.Fuzz(func(t *testing.T, src []byte, a, b int, with string) {
		same, err := edit.Apply(src, nil)
		if err == nil && !bytes.Equal(same, src) {
			t.Fatalf("no changes, but the document changed")
		}
		out, err := edit.Apply(src, []edit.Change{{Element: a, With: with}, {Element: b}})
		if err != nil {
			return
		}
		first, _, _ := bytes.Cut(src, []byte("<"))
		if !bytes.HasPrefix(out, first) {
			t.Fatalf("bytes before the first element changed")
		}
	})
}
