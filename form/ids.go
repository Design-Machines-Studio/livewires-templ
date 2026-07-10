package form

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Form control ids are derived from caller-supplied field names. Those names
// are not guaranteed to be id-safe: "plan tier" or "user[email]" are ordinary
// field names but illegal or fragile as ids. Two things break when an unsafe
// name reaches an id attribute:
//
//   - `id` may not contain whitespace, and an id containing brackets or dots
//     cannot be targeted by a CSS selector or querySelector without escaping.
//   - `aria-describedby` is a space-separated list of IDREFs, so an id with a
//     space in it silently splits into two bogus references and the message it
//     points at is never announced.
//
// Every id emitted by this package therefore passes through sanitizeIDPart.
//
// Ids are derived from the field name alone, so two controls of different kinds
// that share a name (a Field and a Checkbox both named "email") emit the same
// hint and error ids. A component cannot detect that across call sites; give
// controls on one page distinct names.

// nameFallback is used when a name sanitizes to the empty string. An id must
// never be empty or begin with a stray separator.
const nameFallback = "field"

// sanitizeIDPart reduces an arbitrary string to characters safe in an HTML id
// and a CSS selector. Characters outside [A-Za-z0-9-] are escaped as _<hex>_
// rather than folded to a single placeholder, which keeps the mapping
// injective over well-formed UTF-8: two distinct names ("a/b" and "a b") can
// never collide into one id. '_' is itself escaped so it cannot forge an
// escape sequence. Invalid UTF-8 bytes all decode to U+FFFD and so are not
// distinguished, which is immaterial for developer-authored field names.
//
// Distinct names still collide if a caller deliberately picks one that
// reproduces another's derived id (a field named "x-error" collides with the
// error element of a field named "x"). That is inherent to any suffix scheme.
func sanitizeIDPart(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
			b.WriteString(strconv.FormatInt(int64(r), 16))
			b.WriteByte('_')
		}
	}
	return b.String()
}

// sanitizeIDPartOr sanitizes s, falling back to a placeholder when the result
// would be empty.
func sanitizeIDPartOr(s, fallback string) string {
	if id := sanitizeIDPart(s); id != "" {
		return id
	}
	return fallback
}

// controlID is the id of the input a label points at with `for`. Both sides of
// that pair must be sanitized identically or the label stops labelling.
func controlID(name string) string {
	return sanitizeIDPartOr(name, nameFallback)
}

// errorID is the id of a control's error message element.
func errorID(name string) string {
	return controlID(name) + "-error"
}

// hintID is the id of a control's hint element.
func hintID(name string) string {
	return controlID(name) + "-hint"
}

// optionBaseID derives the id prefix for a control whose name is shared across
// several inputs (a radio group, a set of checkboxes). The value disambiguates
// them; sanitizing preserves that, since it never maps two distinct values onto
// one id.
func optionBaseID(name, value string) string {
	return controlID(name) + "-" + sanitizeIDPartOr(value, "option")
}

// labelSpanID and hintSpanID identify the spans nested inside a wrapping label.
func labelSpanID(base string) string { return base + "-label" }
func hintSpanID(base string) string  { return base + "-hint" }

// wrappedAriaAttrs wires an input to the label and hint spans nested inside the
// <label> that wraps it. Such a label folds every descendant string into the
// input's accessible name, so once a hint is nested inside, the input must name
// itself from the label span alone or it announces as "Pro Billed annually".
//
// errID is the id of the error element, or "" when there is no error. Returns
// no attributes when there is nothing to describe.
func wrappedAriaAttrs(base, label, hint, errID string) templ.Attributes {
	if hint == "" && errID == "" {
		return nil
	}

	attrs := templ.Attributes{}
	var refs []string
	if hint != "" {
		// With no label text there is no span to name the input from, and an
		// aria-labelledby pointing at nothing would erase the accessible name.
		if label != "" {
			attrs["aria-labelledby"] = labelSpanID(base)
		}
		refs = append(refs, hintSpanID(base))
	}
	if errID != "" {
		refs = append(refs, errID)
	}
	attrs["aria-describedby"] = strings.Join(refs, " ")
	return attrs
}

// describedByAttrs wires a control to its hint and error text. The hint comes
// first so assistive technology reads the guidance before the failure. Returns
// no attributes when there is nothing to describe, so a control never carries
// an empty aria-describedby.
func describedByAttrs(name, hint, errMsg string) templ.Attributes {
	if hint == "" && errMsg == "" {
		return nil
	}

	var refs []string
	if hint != "" {
		refs = append(refs, hintID(name))
	}
	if errMsg != "" {
		refs = append(refs, errorID(name))
	}
	return templ.Attributes{"aria-describedby": strings.Join(refs, " ")}
}
