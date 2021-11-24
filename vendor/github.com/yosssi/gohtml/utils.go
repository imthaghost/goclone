package gohtml

import (
	"bytes"
	"strings"
	"unicode"
)

type formatterTokenType int

const (
	formatterTokenType_Nothing formatterTokenType = iota
	formatterTokenType_Tag
	formatterTokenType_Text
)

type formattedBuffer struct {
	buffer  *bytes.Buffer
	rawMode bool

	indentString string
	indentLevel  int

	lineWrapColumn       int
	lineWrapMaxSpillover int

	curLineLength int
	prevTokenType formatterTokenType
}

func (bf *formattedBuffer) writeLineFeed() {
	if !bf.rawMode {
		// Strip trailing newlines
		bf.buffer = bytes.NewBuffer(bytes.TrimRightFunc(
			bf.buffer.Bytes(),
			func(r rune) bool {
				return r != '\n' && unicode.IsSpace(r)
			},
		))
	}

	bf.buffer.WriteString("\n")
	bf.curLineLength = 0
	bf.prevTokenType = formatterTokenType_Nothing
}

func (bf *formattedBuffer) writeIndent() {
	bf.buffer.WriteString(strings.Repeat(bf.indentString, bf.indentLevel))
	bf.curLineLength += len(bf.indentString) * bf.indentLevel
}

func (bf *formattedBuffer) writeToken(token string, kind formatterTokenType) {
	if bf.rawMode {
		bf.buffer.WriteString(token)
		bf.curLineLength += len(token)
		return
	}

	if bf.prevTokenType == formatterTokenType_Nothing && strings.TrimSpace(token) == "" {
		// It's a whitespace token, but we already have indentation which functions
		// the same, so we ignore it
		return
	}

	toWrite := token
	if kind == formatterTokenType_Text && bf.prevTokenType == formatterTokenType_Text {
		toWrite = " " + token
	}

	if bf.prevTokenType != formatterTokenType_Nothing && bf.lineWrapColumn > 0 {
		switch {
		case bf.curLineLength > bf.lineWrapColumn:
			// Current line is too long
			fallthrough

		case bf.curLineLength+len(toWrite) > bf.lineWrapColumn+bf.lineWrapMaxSpillover:
			// Current line + new token is too long even with allowed spillover
			fallthrough

		case bf.curLineLength+len(toWrite) > bf.lineWrapColumn &&
			bf.curLineLength > bf.lineWrapColumn-bf.lineWrapMaxSpillover:
			// Current line + new token is too long and doesn't quality for spillover

			bf.writeLineFeed()
			bf.writeToken(token, kind)
			return
		}
	}

	if bf.curLineLength == 0 {
		bf.writeIndent()
	}
	bf.buffer.WriteString(toWrite)
	bf.curLineLength += len(toWrite)
	bf.prevTokenType = kind
}

// unifyLineFeed unifies line feeds.
func unifyLineFeed(s string) string {
	return strings.Replace(strings.Replace(s, "\r\n", "\n", -1), "\r", "\n", -1)
}
