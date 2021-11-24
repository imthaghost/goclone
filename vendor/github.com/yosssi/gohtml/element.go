package gohtml

// An element represents an HTML element.
type element interface {
	isInline() bool
	write(*formattedBuffer, bool) bool
}
