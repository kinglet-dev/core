package edit

import (
	"bytes"
	"encoding/xml"
	"slices"
	"strings"
)

// Namespaces of editor data that would redraw an element's old shape.
const (
	sodipodiNS = "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
	inkscapeNS = "http://www.inkscape.org/namespaces/inkscape"
)

// geometry lists each shape's own geometry attributes, which a path
// doesn't have; a path's d is replaced where it stands.
var geometry = map[string][]string{
	"rect": {"x", "y", "width", "height", "rx", "ry"}, "circle": {"cx", "cy", "r"},
	"ellipse": {"cx", "cy", "rx", "ry"}, "line": {"x1", "y1", "x2", "y2"},
	"polyline": {"points"}, "polygon": {"points"}, "path": {},
}

// pathChars are the only characters path data has, so it needs no
// escaping inside an attribute.
const pathChars = "0123456789.,+-eE MmLlHhVvCcSsQqTtAaZz\t\n\r"

// redraw returns an element's markup as a path with path data d: every
// other attribute exactly as written and in order, and its contents.
func redraw(src []byte, s span, d string) (string, error) {
	if strings.ContainsFunc(d, func(r rune) bool { return !strings.ContainsRune(pathChars, r) }) {
		return "", ErrBadPathData
	}
	drop, ok := geometry[s.tag.Name.Local]
	if !ok {
		return "", ErrNotAShape
	}
	tag := src[s.start:s.tagEnd]
	rawName := tag[1:nameEnd(tag)]
	name := "path"
	if prefix, _, ok := bytes.Cut(rawName, []byte(":")); ok {
		name = string(prefix) + ":" + name
	}
	var b strings.Builder
	b.WriteString("<" + name)
	placed := false
	for i, raw := range rawAttrs(tag) {
		a := s.tag.Attr[i].Name
		switch {
		case a.Space == "" && a.Local == "d":
			b.WriteString(string(raw[:bytes.IndexByte(raw, 'd')]) + `d="` + d + `"`)
			placed = true
		case a.Space == "" && slices.Contains(drop, a.Local), dropped(a):
		default:
			b.Write(raw)
		}
	}
	if !placed {
		b.WriteString(` d="` + d + `"`)
	}
	if s.tagEnd == s.end {
		b.WriteString("/>")
		return b.String(), nil
	}
	body := src[s.tagEnd:s.end]
	b.WriteString(">")
	b.Write(body[:bytes.LastIndex(body, []byte("</"))])
	b.WriteString("</" + name + ">")
	return b.String(), nil
}

// dropped reports editor data that would redraw the old shape: Inkscape
// rebuilds arcs and stars from sodipodi attributes, and path effects from
// original-d.
func dropped(a xml.Name) bool {
	return a.Space == sodipodiNS || (a.Space == inkscapeNS && (a.Local == "original-d" || a.Local == "path-effect"))
}
