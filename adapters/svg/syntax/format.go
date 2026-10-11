package syntax

import (
	"math"
	"strconv"
	"strings"

	"github.com/kinglet-dev/core/domain/geom"
)

// FormatPath writes path data for p with absolute commands: M, L, C and Z.
// Numbers are rounded to a millionth of a user unit, far below anything a
// laser resolves, so the output stays short.
func FormatPath(p geom.Path) string {
	var b strings.Builder
	for _, sub := range p.Subpaths {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString("M" + pair(sub.Start))
		for _, s := range sub.Segments {
			if s.Curved {
				b.WriteString(" C" + pair(s.C1) + " " + pair(s.C2) + " " + pair(s.To))
			} else {
				b.WriteString(" L" + pair(s.To))
			}
		}
		if sub.Closed {
			b.WriteString(" Z")
		}
	}
	return b.String()
}

// pair writes a point's coordinates.
func pair(p geom.Point) string { return number(p.X) + " " + number(p.Y) }

// number writes v rounded to a millionth; adding 0 turns -0 into 0.
func number(v float64) string {
	return strconv.FormatFloat(math.Round(v*1e6)/1e6+0, 'f', -1, 64)
}
