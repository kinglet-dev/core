package edit

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

// span is where one element lies in the document: from the < of its start
// tag to just after its end tag (or the /> of an empty-element tag).
type span struct{ start, end int64 }

// spans finds every element's span, numbered in document order. The
// decoder's offset before a token is where that token starts, and after
// an end tag is where the element ends.
func spans(src []byte) ([]span, error) {
	dec := xml.NewDecoder(bytes.NewReader(src))
	var out []span
	var open []int // elements whose end tag hasn't been read yet
	for {
		at := dec.InputOffset()
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		switch tok.(type) {
		case xml.StartElement:
			open = append(open, len(out))
			out = append(out, span{start: at})
		case xml.EndElement:
			out[open[len(open)-1]].end = dec.InputOffset()
			open = open[:len(open)-1]
		}
	}
}
