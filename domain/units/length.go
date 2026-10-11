// Package units parses and represents physical lengths.
package units

import (
	"math"
	"strconv"
	"strings"
)

// Length is a physical length.
type Length struct {
	mm float64
}

// Millimeters returns the length in millimeters.
func (l Length) Millimeters() float64 { return l.mm }

// Millimetres returns the length in millimeters.
//
// Deprecated: use Millimeters; this spelling will be removed in a later
// release.
func (l Length) Millimetres() float64 { return l.Millimeters() }

// unitLength is the length of every supported unit name.
const unitLength = 2

// maxMillimeters (10 km) rejects absurd values before they reach geometry code.
const maxMillimeters = 10_000_000

// millimetersPer lists the supported units. CSS defines 96 px and 72 pt per inch.
var millimetersPer = map[string]float64{
	"mm": 1,
	"cm": 10,
	"in": 25.4,
	"pt": 25.4 / 72,
	"px": 25.4 / 96,
}

// ParseLength parses a length with a unit, such as "0.25mm" or "0.1in".
// A bare number is refused so that a unit is never guessed.
func ParseLength(s string) (Length, error) {
	letters := len(s) - len(strings.TrimRight(s, "abcdefghijklmnopqrstuvwxyz"))
	if letters == 0 {
		return Length{}, ErrMissingUnit
	}
	// Every supported unit is two letters; anything before them is the number.
	split := len(s) - min(letters, unitLength)
	number, unit := s[:split], s[split:]
	perUnit, known := millimetersPer[unit]
	if !known {
		return Length{}, ErrUnknownUnit
	}
	value, err := strconv.ParseFloat(number, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return Length{}, ErrInvalidNumber
	}
	if value < 0 {
		return Length{}, ErrNegative
	}
	mm := value * perUnit
	if mm > maxMillimeters {
		return Length{}, ErrTooLarge
	}
	return Length{mm: mm}, nil
}
