package examples_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/reader"
	"github.com/kinglet-dev/core/domain/crossing"
)

// The README's example, kept here so its output stays true.
var _ = Describe("The README example", func() {
	It("finds the line a grid of separate squares burns twice", func() {
		// Arrange
		f, err := os.Open("grid.svg")
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(f.Close)

		// Act
		d, err := reader.Read(f)
		Expect(err).NotTo(HaveOccurred())
		r, err := crossing.Find(d.ScoredLines(0.005), crossing.Params{
			Merge: 0.1, Cell: 1, MaxSegments: 1_000_000, MaxEntries: 10_000_000,
		})
		Expect(err).NotTo(HaveOccurred())

		// Assert
		Expect(fmt.Sprintf("%.1f × %.1f mm: %.1f mm of line burns twice, %d crossings",
			d.Width, d.Height, r.Stacked, len(r.Crossings))).
			To(Equal("100.0 × 80.0 mm: 340.0 mm of line burns twice, 0 crossings"))
	})
})
