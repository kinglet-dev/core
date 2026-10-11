package edit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/adapters/svg/edit"
)

var _ = Describe("Apply changes the elements it is told to and nothing else", func() {
	It("returns the document unchanged when there are no changes", func() {
		// Arrange
		src := []byte(`<svg width="10mm"><rect width="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, nil)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(string(src)))
	})

	It("removes an element by its number, the root being 0", func() {
		// Arrange
		src := []byte(`<svg><rect width="1"/><circle r="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 2}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(`<svg><rect width="1"/></svg>`))
	})
	It("replaces an element with new markup", func() {
		// Arrange
		src := []byte(`<svg><rect width="1"/><circle r="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 1, With: `<path d="M0 0 H1"/>`}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(`<svg><path d="M0 0 H1"/><circle r="1"/></svg>`))
	})

	It("removes a group with everything inside it", func() {
		// Arrange
		src := []byte(`<svg><g id="a"><rect width="1"></rect><circle r="1"/></g><line x2="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 1}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(`<svg><line x2="1"/></svg>`))
	})

	It("keeps comments, text and spacing it doesn't change", func() {
		// Arrange
		src := []byte("<?xml version=\"1.0\"?>\n<!-- made by hand -->\n<svg>\n  <rect width=\"1\"/>\n  <text>A &amp; B</text>\n</svg>\n")

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 1}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal("<?xml version=\"1.0\"?>\n<!-- made by hand -->\n<svg>\n  \n  <text>A &amp; B</text>\n</svg>\n"))
	})

	It("applies changes in any order", func() {
		// Arrange
		src := []byte(`<svg><rect width="1"/><circle r="1"/><line x2="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 3}, {Element: 1}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(`<svg><circle r="1"/></svg>`))
	})

	It("changes elements that sit right next to each other", func() {
		// Arrange
		src := []byte(`<svg><rect width="1"/><circle r="1"/></svg>`)

		// Act
		out, err := edit.Apply(src, []edit.Change{{Element: 1, With: "<g/>"}, {Element: 2}})

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(Equal(`<svg><g/></svg>`))
	})
})
