package edit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/edit"
)

// redraw redraws element 1 of an svg holding body with path data d.
func redraw(body, d string) string {
	out, err := edit.Apply([]byte(`<svg>`+body+`</svg>`), []edit.Change{{Element: 1, Redraw: d}})
	Expect(err).NotTo(HaveOccurred())
	return string(out)
}

var _ = Describe("Apply redraws an element as a path", func() {
	It("turns a rect into a path, keeping its other attributes as written", func() {
		// Act
		got := redraw(`<rect id='a' x="1" y="1" width="5" height="5" rx="1" style="stroke:#000" transform="scale(2)"/>`, "M1 1 L6 1")

		// Assert
		Expect(got).To(Equal(`<svg><path id='a' style="stroke:#000" transform="scale(2)" d="M1 1 L6 1"/></svg>`))
	})

	DescribeTable("drops each kind of element's own geometry",
		func(body, want string) {
			// Act
			got := redraw(body, "M0 0 L1 0")

			// Assert
			Expect(got).To(Equal(`<svg>` + want + `</svg>`))
		},
		Entry("circle", `<circle cx="1" cy="1" r="1" fill="none"/>`, `<path fill="none" d="M0 0 L1 0"/>`),
		Entry("ellipse", `<ellipse cx="1" cy="1" rx="1" ry="2"/>`, `<path d="M0 0 L1 0"/>`),
		Entry("line", `<line x1="0" y1="0" x2="1" y2="1" stroke="red"/>`, `<path stroke="red" d="M0 0 L1 0"/>`),
		Entry("polyline", `<polyline points="0 0 1 1"/>`, `<path d="M0 0 L1 0"/>`),
		Entry("polygon", `<polygon points="0 0 1 1 1 0"/>`, `<path d="M0 0 L1 0"/>`),
	)

	It("replaces a path's data where it was written", func() {
		// Act
		got := redraw(`<path id="p" d="M0 0 H5" stroke="#000"/>`, "M0 0 L2 0")

		// Assert
		Expect(got).To(Equal(`<svg><path id="p" d="M0 0 L2 0" stroke="#000"/></svg>`))
	})

	It("keeps the element's contents and renames its end tag", func() {
		// Act
		got := redraw(`<rect width="5" height="5"><title>Lid</title></rect>`, "M0 0 L5 0")

		// Assert
		Expect(got).To(Equal(`<svg><path d="M0 0 L5 0"><title>Lid</title></path></svg>`))
	})

	It("keeps a namespace prefix on the element's name", func() {
		// Act
		got := redraw(`<svg:rect xmlns:svg="http://www.w3.org/2000/svg" width="5" height="5"></svg:rect>`, "M0 0 L5 0")

		// Assert
		Expect(got).To(Equal(`<svg><svg:path xmlns:svg="http://www.w3.org/2000/svg" d="M0 0 L5 0"></svg:path></svg>`))
	})

	It("drops editor data that would redraw the old shape", func() {
		// Arrange: Inkscape redraws arcs from sodipodi attributes and path effects from original-d.
		body := `<path xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" ` +
			`xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" ` +
			`sodipodi:type="arc" sodipodi:cx="1" inkscape:original-d="M0 0" inkscape:path-effect="#e" ` +
			`inkscape:label="Lid" d="M0 0 A1 1 0 0 1 2 0"/>`

		// Act
		got := redraw(body, "M0 0 L2 0")

		// Assert
		Expect(got).To(Equal(`<svg><path xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" ` +
			`xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" inkscape:label="Lid" d="M0 0 L2 0"/></svg>`))
	})

	It("keeps attributes written with spaces around the equals sign and over several lines", func() {
		// Act
		got := redraw("<rect\n  id = \"a\"\n  width=\"5\"\n  height=\"5\" />", "M0 0")

		// Assert
		Expect(got).To(Equal("<svg><path\n  id = \"a\" d=\"M0 0\"/></svg>"))
	})

	It("keeps attributes in other namespaces that share a geometry attribute's name", func() {
		// Act
		got := redraw(`<rect xmlns:x="urn:x" x:width="7" x:d="keep" width="5"/>`, "M0 0")

		// Assert
		Expect(got).To(Equal(`<svg><path xmlns:x="urn:x" x:width="7" x:d="keep" d="M0 0"/></svg>`))
	})

	It("reads attributes separated by tabs and Windows line endings", func() {
		// Act
		got := redraw("<rect\tid=\"a\"\r\n\twidth=\"5\"\r\n/>", "M0 0")

		// Assert
		Expect(got).To(Equal("<svg><path\tid=\"a\" d=\"M0 0\"/></svg>"))
	})

	It("refuses to redraw an element that isn't a shape", func() {
		// Act
		_, err := edit.Apply([]byte(`<svg><g/></svg>`), []edit.Change{{Element: 1, Redraw: "M0 0"}})

		// Assert
		Expect(err).To(MatchError(edit.ErrNotAShape))
	})

	DescribeTable("refuses path data with characters path data never has",
		func(d string) {
			// Act
			_, err := edit.Apply([]byte(`<svg><rect/></svg>`), []edit.Change{{Element: 1, Redraw: d}})

			// Assert
			Expect(err).To(MatchError(edit.ErrBadPathData))
		},
		Entry("a quote", `M0 0" onload="x`),
		Entry("markup", "M0 0<script/>"),
		Entry("a control character", "M0 0\x00"),
		Entry("bytes that aren't text", "M0 0\xff"),
	)
})
