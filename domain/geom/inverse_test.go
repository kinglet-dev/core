package geom_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
)

var _ = Describe("Inverse undoes a transform", func() {
	It("maps every point back where it came from", func() {
		// Arrange: rotate, scale unevenly, skew and move.
		m := geom.Matrix{A: 0.8, B: 0.6, C: -1.2, D: 1.6, E: 7, F: -3}
		p := geom.Point{X: 3.5, Y: -2}

		// Act
		inv, ok := m.Inverse()

		// Assert
		Expect(ok).To(BeTrue())
		back := inv.Apply(m.Apply(p))
		Expect(back.X).To(BeNumerically("~", p.X, near))
		Expect(back.Y).To(BeNumerically("~", p.Y, near))
	})

	It("has none for a transform that flattens everything onto a line", func() {
		// Act
		_, ok := geom.Matrix{A: 1, B: 2, C: 2, D: 4}.Inverse()

		// Assert
		Expect(ok).To(BeFalse())
	})
})
