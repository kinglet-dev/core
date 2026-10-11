package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read records which element each shape came from", func() {
	It("numbers elements in document order, starting with the root at 0", func() {
		// Arrange
		body := `<g><rect width="1" height="1"/></g><circle r="1" fill="black"/>`

		// Act
		d, err := read(svg(mm96, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(HaveLen(2))
		Expect(d.Shapes[0].Element).To(Equal(2))
		Expect(d.Shapes[1].Element).To(Equal(3))
	})

	It("counts elements it skips, so numbers match the document", func() {
		// Arrange: svg 0, defs 1, path 2 (skipped inside defs), rect 3.
		body := `<defs><path d="M0 0 H5"/></defs><rect width="1" height="1"/>`

		// Act
		d, err := read(svg(mm96, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(HaveLen(1))
		Expect(d.Shapes[0].Element).To(Equal(3))
	})
})
