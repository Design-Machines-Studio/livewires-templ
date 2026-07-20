package testutil

import (
	"testing"

	"golang.org/x/net/html"
)

func TestParseFragmentFindElementAndAttrVal(t *testing.T) {
	nodes := ParseFragment(t, `<label data-testid="w"><input type="checkbox" data-on:change="foo()" data-attr:disabled="$busy" aria-label="Hi"/></label>`)

	label := FindElement(nodes, "label")
	if label == nil || label.Type != html.ElementNode {
		t.Fatal("FindElement() did not find the label element")
	}
	input := FindElement(nodes, "input")
	if input == nil || input.Type != html.ElementNode {
		t.Fatal("FindElement() did not find the self-closing input element")
	}

	assertAttrVal(t, input, "data-on:change", "foo()")
	assertAttrVal(t, input, "data-attr:disabled", "$busy")
	assertAttrVal(t, label, "data-testid", "w")
	if value, ok := AttrVal(input, "data-testid"); ok {
		t.Fatalf("AttrVal(input, %q) = %q, true; want false", "data-testid", value)
	}
}

func TestDuplicateAttributesRequireRawTagCounting(t *testing.T) {
	const output = `<input aria-label="a" aria-label="b">`
	input := FindElement(ParseFragment(t, output), "input")
	if input == nil {
		t.Fatal("FindElement() did not find the input element")
	}

	// x/net/html v0.55.0 and newer drop duplicate attributes after the first, so
	// exactly-once assertions use the raw opening tag, not the normalized tree.
	count := 0
	for _, attr := range input.Attr {
		if attr.Key == "aria-label" {
			count++
			if attr.Val != "a" {
				t.Fatalf("parsed aria-label = %q; want first value %q", attr.Val, "a")
			}
		}
	}
	if count != 1 {
		t.Fatalf("parsed aria-label count = %d; want 1", count)
	}
	if got := CountAttr(RawTag(t, output, "input"), "aria-label"); got != 2 {
		t.Fatalf("CountAttr() = %d; want 2", got)
	}
}

func TestCountAttrBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		rawTag string
		attr   string
		want   int
	}{
		{
			name:   "aria label",
			rawTag: `<input aria-label="a" aria-labelledby="b">`,
			attr:   "aria-label",
			want:   1,
		},
		{
			name:   "aria labelledby",
			rawTag: `<input aria-label="a" aria-labelledby="b">`,
			attr:   "aria-labelledby",
			want:   1,
		},
		{
			name:   "left boundary",
			rawTag: `<input data-x-aria-label="y">`,
			attr:   "aria-label",
			want:   0,
		},
		{
			name:   "quoted value",
			rawTag: `<input title="mention aria-label twice aria-label">`,
			attr:   "aria-label",
			want:   0,
		},
		{
			name:   "bare attribute",
			rawTag: `<input checked>`,
			attr:   "checked",
			want:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountAttr(tt.rawTag, tt.attr); got != tt.want {
				t.Fatalf("CountAttr(%q, %q) = %d; want %d", tt.rawTag, tt.attr, got, tt.want)
			}
		})
	}
}

func TestRawTagAllowsGreaterThanInsideQuotedValue(t *testing.T) {
	const output = `<input title="1 > 0" checked><span>after</span>`
	const want = `<input title="1 > 0" checked>`
	if got := RawTag(t, output, "input"); got != want {
		t.Fatalf("RawTag() = %q; want %q", got, want)
	}
}

func TestRawTagMatchesCompleteTagName(t *testing.T) {
	// A prefix-sharing tag (<input-group>) must not be mistaken for <input>.
	const output = `<input-group data-role="wrap"><input type="checkbox" checked></input-group>`
	const want = `<input type="checkbox" checked>`
	if got := RawTag(t, output, "input"); got != want {
		t.Fatalf("RawTag() = %q; want %q", got, want)
	}
}

func assertAttrVal(t *testing.T, n *html.Node, key, want string) {
	t.Helper()
	got, ok := AttrVal(n, key)
	if !ok {
		t.Fatalf("AttrVal(%q) was not found", key)
	}
	if got != want {
		t.Fatalf("AttrVal(%q) = %q; want %q", key, got, want)
	}
}
