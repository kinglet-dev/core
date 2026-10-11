package edit

// isSpace reports XML whitespace.
func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' }

// nameEnd is where the element name in a start tag ends, or the / or >
// that closes a tag with no attributes (the name is only used for its
// prefix, which comes before either). A start tag always ends with >.
func nameEnd(tag []byte) int {
	i := 1
	for !isSpace(tag[i]) && tag[i] != '>' {
		i++
	}
	return i
}

// rawAttrs splits a well-formed start tag into its attributes as written,
// each with the whitespace before it, in document order: the order the
// decoder lists them in.
func rawAttrs(tag []byte) [][]byte {
	var out [][]byte
	i := nameEnd(tag)
	for {
		from := i
		for isSpace(tag[i]) {
			i++
		}
		if tag[i] == '/' || tag[i] == '>' {
			return out
		}
		for tag[i] != '=' {
			i++
		}
		i++
		for isSpace(tag[i]) {
			i++
		}
		quote := tag[i]
		i++
		for tag[i] != quote {
			i++
		}
		i++
		out = append(out, tag[from:i])
	}
}
