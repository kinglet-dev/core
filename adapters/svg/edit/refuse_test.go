package edit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/edit"
)

var _ = Describe("Apply refuses changes it can't make safely", func() {
	var src []byte

	BeforeEach(func() {
		// Arrange (shared): svg 0, g 1, rect 2 inside the group, circle 3.
		src = []byte(`<svg><g><rect width="1"/></g><circle r="1"/></svg>`)
	})

	DescribeTable("refuses an element number the document doesn't have",
		func(n int) {
			// Act
			_, err := edit.Apply(src, []edit.Change{{Element: n}})

			// Assert
			Expect(err).To(MatchError(edit.ErrNoSuchElement))
		},
		Entry("one past the last element", 4),
		Entry("a negative number", -1),
	)

	It("refuses to change an element and one inside it", func() {
		// Act
		_, err := edit.Apply(src, []edit.Change{{Element: 2}, {Element: 1}})

		// Assert
		Expect(err).To(MatchError(edit.ErrOverlap))
	})

	It("refuses to change the same element twice", func() {
		// Act
		_, err := edit.Apply(src, []edit.Change{{Element: 3}, {Element: 3, With: "<g/>"}})

		// Assert
		Expect(err).To(MatchError(edit.ErrOverlap))
	})

	It("refuses a document that isn't well-formed XML", func() {
		// Act
		_, err := edit.Apply([]byte(`<svg><rect></svg>`), []edit.Change{{Element: 1}})

		// Assert
		Expect(err).To(MatchError(edit.ErrBadXML))
	})

	It("leaves the caller's list of changes as it was", func() {
		// Arrange
		changes := []edit.Change{{Element: 3}, {Element: 2}}

		// Act
		_, err := edit.Apply(src, changes)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(changes).To(Equal([]edit.Change{{Element: 3}, {Element: 2}}))
	})
})
