package livewires

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Initials extracts initials from a name (e.g., "John Doe" → "JD").
func Initials(name string) string {
	parts := strings.Fields(name)
	initials := ""
	for _, part := range parts {
		if len(part) > 0 {
			initials += strings.ToUpper(string(part[0]))
		}
	}
	if initials == "" {
		return "?"
	}
	return initials
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

// DateShort formats a date string as "Jan 2, 2006".
func DateShort(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}
	var t time.Time
	var err error
	for _, fmt := range formats {
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

// Year extracts the year from a date string.
func Year(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return ""
		}
	}
	return t.Format("2006")
}

// DateTimeFriendly formats a datetime as "January 2, 2006 at 3:04 PM".
func DateTimeFriendly(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	var t time.Time
	var err error
	for _, fmt := range formats {
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

// PtrValOr returns the string value of a pointer, or a default if nil/empty.
func PtrValOr(s *string, def string) string {
	if s != nil && *s != "" {
		return *s
	}
	return def
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
