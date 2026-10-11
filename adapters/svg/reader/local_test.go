package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
)

var _ = Describe("Read keeps each shape's path as the element draws it", func() {
	It("records the path in the element's own units and the transform to millimeters", func() {
		// Arrange: one user unit is a millimeter, and the group doubles it and moves it 5 mm right.
		body := `<g transform="translate(5 0) scale(2)"><path d="M1 1 H3"/></g>`

		// Act
		d, err := read(svg(`width="100mm" height="100mm" viewBox="0 0 100 100"`, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		s := d.Shapes[0]
		Expect(s.Local.Subpaths[0].Start).To(Equal(geom.Point{X: 1, Y: 1}))
		Expect(s.ToMM.Apply(geom.Point{X: 3, Y: 1})).To(Equal(geom.Point{X: 11, Y: 2}))
		Expect(s.Local.Transform(s.ToMM)).To(Equal(s.Path))
	})
})
