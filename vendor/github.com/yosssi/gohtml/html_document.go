package gohtml

import "bytes"

// Column to wrap lines to (disabled by default)
var LineWrapColumn = 0

// Maxmimum characters a long word can extend past LineWrapColumn without wrapping
var LineWrapMaxSpillover = 5

// An htmlDocument represents an HTML document.
type htmlDocument struct {
	elements []element
}

// html generates an HTML source code and returns it.
func (htmlDoc *htmlDocument) html() string {
	return string(htmlDoc.bytes())
}

// bytes reads from htmlDocument's internal array of elements and returns HTML source code
func (htmlDoc *htmlDocument) bytes() []byte {
	bf := &formattedBuffer{
		buffer: &bytes.Buffer{},

		lineWrapColumn:       LineWrapColumn,
		lineWrapMaxSpillover: LineWrapMaxSpillover,

		indentString: defaultIndentString,
		indentLevel:  startIndent,
	}

	isPreviousNodeInline := true
	for _, child := range htmlDoc.elements {
		isPreviousNodeInline = child.write(bf, isPreviousNodeInline)
	}
	return bf.buffer.Bytes()
}

// append appends an element to the htmlDocument.
func (htmlDoc *htmlDocument) append(e element) {
	htmlDoc.elements = append(htmlDoc.elements, e)
}
