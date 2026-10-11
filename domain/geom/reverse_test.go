package geom_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
)

var _ = Describe("Reverse runs a subpath the other way", func() {
	pt := func(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

	It("starts at the end and visits the same points backward", func() {
		// Arrange
		sub := geom.Subpath{Start: pt(0, 0), Segments: []geom.Segment{
			geom.LineTo(pt(10, 0)),
			geom.CubicTo(pt(12, 0), pt(14, 2), pt(14, 4)),
		}}

		// Act
		got := sub.Reverse()

		// Assert: the curve's control points swap.
		Expect(got).To(Equal(geom.Subpath{Start: pt(14, 4), Segments: []geom.Segment{
			geom.CubicTo(pt(14, 2), pt(12, 0), pt(10, 0)),
			geom.LineTo(pt(0, 0)),
		}}))
	})

	It("keeps a closed subpath closed", func() {
		// Arrange
		sub := geom.Subpath{Start: pt(0, 0), Segments: []geom.Segment{geom.LineTo(pt(1, 0)), geom.LineTo(pt(1, 1))}, Closed: true}

		// Act
		got := sub.Reverse()

		// Assert
		Expect(got.Closed).To(BeTrue())
		Expect(got.Reverse()).To(Equal(sub))
	})

	It("leaves a lone point where it is", func() {
		// Arrange
		sub := geom.Subpath{Start: pt(3, 4), Segments: []geom.Segment{}}

		// Act
		got := sub.Reverse()

		// Assert
		Expect(got).To(Equal(sub))
	})
})
