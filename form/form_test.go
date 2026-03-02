package form

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderToString is a test helper that renders a templ component to a string.
func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return buf.String()
}

func TestFieldRenders(t *testing.T) {
	html := renderToString(t, FieldText("Name", "name", "", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSelectRenders(t *testing.T) {
	opts := []SelectOption{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
	}
	html := renderToString(t, SelectSimple("Country", "country", "Choose...", opts, false))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestCheckboxRenders(t *testing.T) {
	html := renderToString(t, CheckboxSimple("agree", "I agree", false))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestRadioGroupRenders(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := renderToString(t, RadioGroup(data))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSwitchRenders(t *testing.T) {
	html := renderToString(t, Switch("notifications", "notif", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTextareaRenders(t *testing.T) {
	html := renderToString(t, TextareaSimple("Message", "message", "", "Write...", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSearchRenders(t *testing.T) {
	html := renderToString(t, SearchSimple("q", "", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFilterRenders(t *testing.T) {
	data := FilterProps{
		Title: "Status",
		Name:  "status",
		Options: []FilterOption{
			{Value: "active", Label: "Active"},
			{Value: "archived", Label: "Archived"},
		},
	}
	html := renderToString(t, Filter(data))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestDateRangeRenders(t *testing.T) {
	data := DateRangeProps{
		StartName: "from",
		EndName:   "to",
		Label:     "Date Range",
	}
	html := renderToString(t, DateRange(data))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
