// Package edit changes chosen elements of an SVG document and leaves every
// other byte as it was, so ids, styles, comments and layout survive.
// Elements are numbered as the reader numbers them: in document order,
// the root svg element first, from 0.
package edit

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// Errors Apply returns, wrapped with details.
var (
	ErrNoSuchElement = errors.New("the document has no such element")
	ErrOverlap       = errors.New("two changes overlap")
	ErrBadXML        = errors.New("the document isn't well-formed XML")
)

// Change replaces one element, start tag to end tag, with new markup;
// empty markup removes it.
type Change struct {
	Element int    // the element's number
	With    string // its replacement; "" removes it
}

// Apply returns the document with the changes made.
func Apply(src []byte, changes []Change) ([]byte, error) {
	all, err := spans(src)
	if err != nil {
		return nil, fmt.Errorf("%w (%v)", ErrBadXML, err)
	}
	// Element numbers follow document order, so sorting by number puts the
	// changes in the order they appear. Sort a copy; the caller's is theirs.
	changes = slices.SortedFunc(slices.Values(changes), func(a, b Change) int { return cmp.Compare(a.Element, b.Element) })
	var out []byte
	from := int64(0)
	for _, c := range changes {
		if c.Element < 0 || c.Element >= len(all) {
			return nil, fmt.Errorf("%w: %d (it has %d)", ErrNoSuchElement, c.Element, len(all))
		}
		// Spans are nested or apart; in document order, one that starts
		// before the last change ended lies inside it, or is the same.
		s := all[c.Element]
		if s.start < from {
			return nil, fmt.Errorf("%w: element %d is inside or the same as one already changed", ErrOverlap, c.Element)
		}
		out = append(out, src[from:s.start]...)
		out = append(out, c.With...)
		from = s.end
	}
	return append(out, src[from:]...), nil
}
