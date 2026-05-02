package livewires

import (
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// SingleInitial returns the first initial from a name (e.g., "John Doe" → "J").
// Used for small avatar sizes (xs, sm) where space is limited.
func SingleInitial(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	return Initials(parts[0])
}

const maxInitials = 2

// Initials extracts up to two initials from a name (e.g., "John Doe" → "JD").
func Initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	var b strings.Builder
	b.Grow(maxInitials * utf8.UTFMax)
	for i, part := range parts {
		if i >= maxInitials {
			break
		}
		r, _ := utf8.DecodeRuneInString(part)
		if r != utf8.RuneError {
			b.WriteRune(unicode.ToUpper(r))
		}
	}
	if b.Len() == 0 {
		return "?"
	}
	return b.String()
}

// Title converts a string to title case.
func Title(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	capitalizeNext := true
	for _, r := range s {
		if unicode.IsSpace(r) {
			capitalizeNext = true
			result.WriteRune(r)
		} else if capitalizeNext {
			result.WriteRune(unicode.ToTitle(r))
			capitalizeNext = false
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// dateShortFormats lists formats tried by DateShort.
var dateShortFormats = []string{
	"2006-01-02",
	"2006-01-02T15:04",
	"2006-01-02T15:04:05",
	time.RFC3339,
}

// DateShort formats a date string as "Jan 2, 2006".
func DateShort(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	var t time.Time
	var err error
	for _, fmt := range dateShortFormats {
		t, err = time.Parse(fmt, dateStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 2, 2006")
}

// dateTimeFriendlyFormats lists formats tried by DateTimeFriendly.
var dateTimeFriendlyFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05",
	time.RFC3339,
}

// DateTimeFriendly formats a datetime as "January 2, 2006 at 3:04 PM".
func DateTimeFriendly(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	var t time.Time
	var err error
	for _, fmt := range dateTimeFriendlyFormats {
		t, err = time.Parse(fmt, dateStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		return dateStr
	}
	return t.Format("January 2, 2006 at 3:04 PM")
}

// Itoa converts an int to string.
func Itoa(n int) string {
	return strconv.Itoa(n)
}

// PtrVal returns the string value of a pointer, or empty string if nil.
func PtrVal(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// SplitCSV splits a comma-separated string into trimmed, non-empty items.
func SplitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
