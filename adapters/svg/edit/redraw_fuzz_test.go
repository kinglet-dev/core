package edit_test

import (
	"encoding/xml"
	"testing"

	"github.com/kinglet-dev/core/adapters/svg/edit"
)

// FuzzRedraw checks that redrawing never panics on hostile input and that
// a redrawn document is still well-formed XML.
func FuzzRedraw(f *testing.F) {
	f.Add([]byte(`<svg><rect id='a' x = "1" width="5"><title>t</title></rect></svg>`), 1, "M0 0 L1 1")
	f.Add([]byte(`<svg:svg xmlns:svg="http://www.w3.org/2000/svg"><svg:path d="M0 0" a:b="&quot;>" xmlns:a="x"/></svg:svg>`), 1, `M0 0 "&<`)
	f.Fuzz(func(t *testing.T, src []byte, n int, d string) {
		out, err := edit.Apply(src, []edit.Change{{Element: n, Redraw: d}})
		if err != nil || d == "" {
			return
		}
		var v any
		if err := xml.Unmarshal(out, &v); err != nil {
			t.Fatalf("the redrawn document isn't well-formed: %v\n%s", err, out)
		}
	})
}
