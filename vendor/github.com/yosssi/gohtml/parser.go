package gohtml

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// parse parses a stirng and converts it into an html.
func parse(r io.Reader) *htmlDocument {
	htmlDoc := &htmlDocument{}
	tokenizer := html.NewTokenizer(r)
	for {
		if errorToken, _, _ := parseToken(tokenizer, htmlDoc, nil); errorToken {
			break
		}
	}
	return htmlDoc
}

// Function that identifies which tags will be treated as containing preformatted
// content. Such tags will have the formatting of all its contents preserved
// unchanged.
// The opening tag html.Token is passed to this function.
// By default, only <pre> and <textarea> tags are considered preformatted.
var IsPreformatted = func(token html.Token) bool {
	return token.Data == "pre" || token.Data == "textarea"
}

func parseToken(tokenizer *html.Tokenizer, htmlDoc *htmlDocument, parent *tagElement) (bool, bool, string) {
	tokenType := tokenizer.Next()
	switch tokenType {
	case html.ErrorToken:
		return true, false, ""
	case html.TextToken:
		text := string(tokenizer.Raw())
		if strings.TrimSpace(text) == "" && (parent == nil || !parent.isRaw) {
			break
		}
		textElement := &textElement{text: text, parent: parent}
		appendElement(htmlDoc, parent, textElement)
	case html.StartTagToken:
		raw := string(tokenizer.Raw())
		token := tokenizer.Token()
		tagElement := &tagElement{
			tagName:     string(token.Data),
			startTagRaw: raw,
			isRaw:       IsPreformatted(token) || (parent != nil && parent.isRaw),
			parent:      parent,
		}
		appendElement(htmlDoc, parent, tagElement)
		for {
			errorToken, parentEnded, unsetEndTag := parseToken(tokenizer, htmlDoc, tagElement)
			if errorToken {
				return true, false, ""
			}
			if parentEnded {
				if unsetEndTag != "" {
					return false, false, unsetEndTag
				}
				break
			}
			if unsetEndTag != "" {
				return false, false, setEndTagRaw(tokenizer, tagElement, unsetEndTag)
			}
		}
	case html.EndTagToken:
		return false, true, setEndTagRaw(tokenizer, parent, getTagName(tokenizer))
	case html.DoctypeToken, html.SelfClosingTagToken, html.CommentToken:
		tagElement := &tagElement{
			tagName:     getTagName(tokenizer),
			startTagRaw: string(tokenizer.Raw()),
			isRaw:       parent != nil && parent.isRaw,
			parent:      parent,
		}
		appendElement(htmlDoc, parent, tagElement)
	}
	return false, false, ""
}

// appendElement appends the element to the htmlDocument or parent tagElement.
func appendElement(htmlDoc *htmlDocument, parent *tagElement, e element) {
	if parent != nil {
		parent.appendChild(e)
	} else {
		htmlDoc.append(e)
	}
}

// getTagName gets a tagName from tokenizer.
func getTagName(tokenizer *html.Tokenizer) string {
	tagName, _ := tokenizer.TagName()
	return string(tagName)
}

// setEndTagRaw sets an endTagRaw to the parent.
func setEndTagRaw(tokenizer *html.Tokenizer, parent *tagElement, tagName string) string {
	if parent != nil && parent.tagName == tagName {
		parent.endTagRaw = string(tokenizer.Raw())
		return ""
	}
	return tagName
}
