package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read records each shape's fill and stroke", func() {
	It("records a hairline's stroke in lower case, with no fill", func() {
		// Act
		s := only(`<rect width="1" height="1" fill="none" stroke="#FF0000"/>`)

		// Assert
		Expect(s.Stroke).To(Equal("#ff0000"))
		Expect(s.Fill).To(Equal("none"))
	})

	It("records paint inherited from a group", func() {
		// Act
		s := only(`<g stroke="blue" fill="none"><circle r="1"/></g>`)

		// Assert
		Expect(s.Stroke).To(Equal("blue"))
		Expect(s.Fill).To(Equal("none"))
	})

	It("records SVG's defaults: black fill and no stroke", func() {
		// Act
		s := only(`<rect width="1" height="1"/>`)

		// Assert
		Expect(s.Fill).To(Equal("black"))
		Expect(s.Stroke).To(Equal("none"))
	})

	It("records no fill for a line, which has no area", func() {
		// Act
		s := only(`<line x2="1" stroke="black"/>`)

		// Assert
		Expect(s.Fill).To(Equal("none"))
	})
})
