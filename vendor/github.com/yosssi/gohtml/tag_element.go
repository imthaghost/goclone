package gohtml

import "bytes"

// A tagElement represents a tag element of an HTML document.
type tagElement struct {
	tagName     string
	startTagRaw string
	endTagRaw   string

	parent   *tagElement
	children []element

	isRaw                 bool
	isChildrenInlineCache *bool
}

// Enable condensing a tag with only inline children onto a single line, or
// completely inlining it with sibling nodes.
// Tags to be treated as inline can be set in `InlineTags`.
// Only inline tags will be completely inlined, while other condensable tags
// will be given their own dedicated (single) line.
var Condense bool

// Tags that are considered inline tags.
// Note: Text nodes are always considered to be inline
var InlineTags = map[string]bool{
	"a":      true,
	"code":   true,
	"em":     true,
	"span":   true,
	"strong": true,
}

// Maximum length of an opening inline tag before it's un-inlined
var InlineTagMaxLength = 40

func (e *tagElement) isInline() bool {
	if e.isRaw || !InlineTags[e.tagName] || len(e.startTagRaw) > InlineTagMaxLength {
		return false
	}
	return e.isChildrenInline()
}

func (e *tagElement) isChildrenInline() bool {
	if !Condense {
		return false
	}
	if e.isChildrenInlineCache != nil {
		return *e.isChildrenInlineCache
	}

	isInline := true
	for _, child := range e.children {
		isInline = isInline && child.isInline()
	}

	e.isChildrenInlineCache = &isInline
	return isInline
}

// write writes a tag to the buffer.
func (e *tagElement) write(bf *formattedBuffer, isPreviousNodeInline bool) bool {
	if e.isRaw {
		if e.parent != nil && !e.parent.isRaw {
			bf.writeLineFeed()
			bf.writeIndent()
			bf.rawMode = true
			defer func() {
				bf.rawMode = false
			}()
		}
		bf.writeToken(e.startTagRaw, formatterTokenType_Tag)
		for _, child := range e.children {
			child.write(bf, true)
		}
		bf.writeToken(e.endTagRaw, formatterTokenType_Tag)
		return false
	}

	if e.isChildrenInline() && (e.endTagRaw != "" || e.isInline()) {
		// Write the condensed output to a separate buffer, in case it doesn't work out
		condensedBuffer := *bf
		condensedBuffer.buffer = &bytes.Buffer{}

		if bf.buffer.Len() > 0 && (!isPreviousNodeInline || !e.isInline()) {
			condensedBuffer.writeLineFeed()
		}
		condensedBuffer.writeToken(e.startTagRaw, formatterTokenType_Tag)
		if !isPreviousNodeInline && e.endTagRaw != "" {
			condensedBuffer.indentLevel++
		}

		for _, child := range e.children {
			child.write(&condensedBuffer, true)
		}
		if e.endTagRaw != "" {
			condensedBuffer.writeToken(e.endTagRaw, formatterTokenType_Tag)
			if !isPreviousNodeInline {
				condensedBuffer.indentLevel--
			}
		}

		if e.isInline() || bytes.IndexAny(condensedBuffer.buffer.Bytes()[1:], "\n") == -1 {
			// If we're an inline tag, or there were no newlines were in the buffer,
			// replace the original with the condensed version
			condensedBuffer.buffer = bytes.NewBuffer(bytes.Join([][]byte{
				bf.buffer.Bytes(), condensedBuffer.buffer.Bytes(),
			}, []byte{}))
			*bf = condensedBuffer

			return e.isInline()
		}
	}

	if bf.buffer.Len() > 0 {
		bf.writeLineFeed()
	}
	bf.writeToken(e.startTagRaw, formatterTokenType_Tag)
	if e.endTagRaw != "" {
		bf.indentLevel++
	}

	isPreviousNodeInline = false
	for _, child := range e.children {
		isPreviousNodeInline = child.write(bf, isPreviousNodeInline)
	}

	if e.endTagRaw != "" {
		if len(e.children) > 0 {
			bf.writeLineFeed()
		}
		bf.indentLevel--
		bf.writeToken(e.endTagRaw, formatterTokenType_Tag)
	}

	return false
}

// appendChild append an element to the element's children.
func (e *tagElement) appendChild(child element) {
	e.children = append(e.children, child)
}
