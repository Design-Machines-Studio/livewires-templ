// Package testutil provides shared test helpers for rendering templ components.
package testutil

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ParseFragment parses an HTML fragment in a body context, failing the test on error.
func ParseFragment(t *testing.T, fragment string) []*html.Node {
	t.Helper()
	context := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	nodes, err := html.ParseFragment(strings.NewReader(fragment), context)
	if err != nil {
		t.Fatalf("parse fragment failed: %v", err)
	}
	return nodes
}

// FindElement returns the first depth-first element matching tag across nodes.
func FindElement(nodes []*html.Node, tag string) *html.Node {
	for _, node := range nodes {
		if found := findElement(node, tag); found != nil {
			return found
		}
	}
	return nil
}

func findElement(node *html.Node, tag string) *html.Node {
	if node.Type == html.ElementNode && node.Data == tag {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElement(child, tag); found != nil {
			return found
		}
	}
	return nil
}

// AttrVal returns an attribute value by key. x/net/html lowercases attribute
// keys, so callers pass lowercase keys. Colon keys survive on non-foreign elements.
func AttrVal(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

// RawTag returns the first raw opening tag matching tag, failing the test if absent.
func RawTag(t *testing.T, output, tag string) string {
	t.Helper()
	needle := "<" + tag
	// Match a complete tag name: the character after the name must be a tag-name
	// boundary (whitespace, `>`, or `/`), so a request for "input" never selects
	// an earlier "<input-group>".
	start := -1
	for from := 0; ; {
		idx := strings.Index(output[from:], needle)
		if idx == -1 {
			break
		}
		abs := from + idx
		after := abs + len(needle)
		if after < len(output) {
			if c := output[after]; isHTMLSpace(c) || c == '>' || c == '/' {
				start = abs
				break
			}
		}
		from = abs + len(needle)
	}
	if start == -1 {
		t.Fatalf("raw opening tag %q not found", needle)
	}

	var quote byte
	for i := start + len(needle); i < len(output); i++ {
		switch output[i] {
		case '\'', '"':
			if quote == 0 {
				quote = output[i]
			} else if quote == output[i] {
				quote = 0
			}
		case '>':
			if quote == 0 {
				return output[start : i+1]
			}
		}
	}
	t.Fatalf("raw opening tag %q is not terminated", needle)
	return ""
}

// CountAttr counts complete attribute-name tokens outside quoted values in rawTag.
func CountAttr(rawTag, name string) int {
	if name == "" {
		return 0
	}

	count := 0
	var quote byte
	for i := 0; i < len(rawTag); i++ {
		switch rawTag[i] {
		case '\'', '"':
			if quote == 0 {
				quote = rawTag[i]
			} else if quote == rawTag[i] {
				quote = 0
			}
			continue
		}
		if quote != 0 || i == 0 || !isHTMLSpace(rawTag[i-1]) || !strings.HasPrefix(rawTag[i:], name) {
			continue
		}

		after := i + len(name)
		if after < len(rawTag) && (rawTag[after] == '=' || rawTag[after] == '/' || rawTag[after] == '>' || isHTMLSpace(rawTag[after])) {
			count++
			i = after - 1
		}
	}
	return count
}

func isHTMLSpace(value byte) bool {
	switch value {
	case ' ', '\t', '\n', '\f', '\r':
		return true
	default:
		return false
	}
}
