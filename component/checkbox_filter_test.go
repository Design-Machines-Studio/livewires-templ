package component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func checkboxFilterStatic(markup string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, markup)
		return err
	})
}

func checkboxFilterOptions() []CheckboxFilterOption {
	return []CheckboxFilterOption{
		{ID: "circle-0", Value: "ops", Label: "Operations", Count: 3, Checked: true},
		{ID: "circle-1", Value: "board", Label: "Board", Count: 1},
	}
}

// Members-equivalent props: the Baseplate Members Circles filter expressed through the component.
const checkboxFilterGoldenMembers = `<details class="dropdown dropdown--panel" data-on:click__outside="el.open = false"><summary class="trigger button--subtle button--small" aria-controls="member-circle-options"><svg class="icon"><use href="#circle"></use></svg><span>Circles</span> <span class="text-muted ml-025" data-text="summaryExpr">1 selected</span> <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg></summary><div class="menu"><div id="member-circle-options" class="dropdown-body stack stack-compact px-05 py-025" role="group" aria-label="Filter by circle"><label class="checkbox" for="member-circle-1"><input type="checkbox" id="member-circle-1" name="circle" value="ops" checked data-indicator:member-filter-pending="" data-on:change="changeAct"> <span>Ops</span> <span class="text-muted">3<span class="sr-only">members</span></span></label><label class="checkbox" for="member-circle-2"><input type="checkbox" id="member-circle-2" name="circle" value="board" data-indicator:member-filter-pending="" data-on:change="changeAct"> <span>Board</span> <span class="text-muted">1<span class="sr-only">members</span></span></label></div><footer class="cluster cluster-start pb-0 pt-025"><button type="submit" name="clear" value="circles" class="button button--text button--small">Clear all</button></footer></div></details>`

func TestCheckboxFilterMembersEquivalent(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		OptionsID:    "member-circle-options",
		Label:        "Circles",
		Summary:      "1 selected",
		GroupLabel:   "Filter by circle",
		Name:         "circle",
		CountLabel:   "members",
		Icon:         checkboxFilterStatic(`<svg class="icon"><use href="#circle"></use></svg>`),
		SummaryAttrs: templ.Attributes{"data-text": "summaryExpr"},
		InputAttrs:   templ.Attributes{"data-indicator:member-filter-pending": "", "data-on:change": "changeAct"},
		Options: []CheckboxFilterOption{
			{ID: "member-circle-1", Value: "ops", Label: "Ops", Count: 3, Checked: true},
			{ID: "member-circle-2", Value: "board", Label: "Board", Count: 1},
		},
		BodyClass:   "stack stack-compact px-05 py-025",
		FooterClass: "cluster cluster-start pb-0 pt-025",
		Footer:      checkboxFilterStatic(`<button type="submit" name="clear" value="circles" class="button button--text button--small">Clear all</button>`),
		Attrs:       templ.Attributes{"data-on:click__outside": "el.open = false"},
	}))
	if html != checkboxFilterGoldenMembers {
		t.Errorf("Members-equivalent output mismatch\n got: %s\nwant: %s", html, checkboxFilterGoldenMembers)
	}
}

func TestCheckboxFilterStructure(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		ID:        "types",
		TriggerID: "types-trigger",
		OptionsID: "type-options",
		Label:     "Proposal types",
		Summary:   "All",
		Name:      "type",
		Options:   checkboxFilterOptions(),
		Class:     "extra",
	}))
	nodes := testutil.ParseFragment(t, html)

	details := testutil.FindElement(nodes, "details")
	if details == nil {
		t.Fatal("expected details element")
	}
	if v, _ := testutil.AttrVal(details, "class"); v != "dropdown dropdown--panel extra" {
		t.Errorf("details class = %q", v)
	}
	if v, _ := testutil.AttrVal(details, "id"); v != "types" {
		t.Errorf("details id = %q", v)
	}

	summary := testutil.FindElement(nodes, "summary")
	if v, _ := testutil.AttrVal(summary, "class"); v != "trigger button--subtle button--small" {
		t.Errorf("summary class = %q", v)
	}
	if v, _ := testutil.AttrVal(summary, "id"); v != "types-trigger" {
		t.Errorf("summary id = %q", v)
	}
	if v, _ := testutil.AttrVal(summary, "aria-controls"); v != "type-options" {
		t.Errorf("summary aria-controls = %q", v)
	}
	if muted := testutil.FindElementByClass(nodes, "ml-025"); muted == nil || testutil.NodeText(muted) != "All" {
		t.Error("expected summary text in muted span")
	}

	if testutil.FindElementByClass(nodes, "menu") == nil {
		t.Error("expected .menu wrapper")
	}
	body := testutil.FindElementByID(nodes, "type-options")
	if body == nil {
		t.Fatal("expected options group with OptionsID")
	}
	if v, _ := testutil.AttrVal(body, "class"); v != "dropdown-body" {
		t.Errorf("body class = %q, want only dropdown-body when BodyClass is empty", v)
	}
	if v, _ := testutil.AttrVal(body, "role"); v != "group" {
		t.Errorf("body role = %q", v)
	}
	if v, _ := testutil.AttrVal(body, "aria-label"); v != "Proposal types" {
		t.Errorf("group aria-label = %q, want Label fallback", v)
	}
	if strings.Contains(html, "stack") || strings.Contains(html, "cluster") {
		t.Error("component must not emit layout primitives by default")
	}
}

func TestCheckboxFilterLabelInputAssociation(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilter("opts", "Circles", "circle", checkboxFilterOptions()))
	nodes := testutil.ParseFragment(t, html)
	for _, opt := range checkboxFilterOptions() {
		input := testutil.FindElementByID(nodes, opt.ID)
		if input == nil {
			t.Fatalf("expected input %q", opt.ID)
		}
		if v, _ := testutil.AttrVal(input, "type"); v != "checkbox" {
			t.Errorf("input type = %q", v)
		}
		if v, _ := testutil.AttrVal(input, "name"); v != "circle" {
			t.Errorf("input name = %q", v)
		}
		if v, _ := testutil.AttrVal(input, "value"); v != opt.Value {
			t.Errorf("input value = %q, want %q", v, opt.Value)
		}
		label := input.Parent
		if label.Data != "label" {
			t.Fatalf("expected input nested in label, got %q", label.Data)
		}
		if v, _ := testutil.AttrVal(label, "for"); v != opt.ID {
			t.Errorf("label for = %q, want %q", v, opt.ID)
		}
		if v, _ := testutil.AttrVal(label, "class"); v != "checkbox" {
			t.Errorf("label class = %q", v)
		}
	}
}

func TestCheckboxFilterCheckedState(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilter("opts", "Circles", "circle", checkboxFilterOptions()))
	nodes := testutil.ParseFragment(t, html)
	if _, ok := testutil.AttrVal(testutil.FindElementByID(nodes, "circle-0"), "checked"); !ok {
		t.Error("expected checked option to have checked attribute")
	}
	if _, ok := testutil.AttrVal(testutil.FindElementByID(nodes, "circle-1"), "checked"); ok {
		t.Error("expected unchecked option to omit checked attribute")
	}
}

func TestCheckboxFilterCounts(t *testing.T) {
	withLabel := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		Label: "Types", Name: "type", CountLabel: "decisions", Options: checkboxFilterOptions(),
	}))
	if !strings.Contains(withLabel, `<span class="text-muted">3<span class="sr-only">decisions</span></span>`) {
		t.Error("expected count with screen-reader label")
	}
	without := testutil.RenderToString(t, CheckboxFilter("opts", "Types", "type", checkboxFilterOptions()))
	if strings.Contains(without, "sr-only") {
		t.Error("expected no sr-only span when CountLabel is empty")
	}
	if !strings.Contains(without, `<span class="text-muted">3</span>`) {
		t.Error("expected bare count when CountLabel is empty")
	}
}

func TestCheckboxFilterEscaping(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		Label:   `<b>Types</b>`,
		Summary: `"x" & <y>`,
		Name:    "type",
		Options: []CheckboxFilterOption{{ID: "t-0", Value: `a"><script>`, Label: `<script>alert(1)</script>`}},
	}))
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Errorf("expected labels and values escaped, got %s", html)
	}
	nodes := testutil.ParseFragment(t, html)
	if v, _ := testutil.AttrVal(testutil.FindElementByID(nodes, "t-0"), "value"); v != `a"><script>` {
		t.Errorf("value round-trip = %q", v)
	}
}

func TestCheckboxFilterOptionalIconAndFooter(t *testing.T) {
	bare := testutil.RenderToString(t, CheckboxFilter("opts", "Circles", "circle", checkboxFilterOptions()))
	if strings.Contains(bare, "<footer") {
		t.Error("expected footer omitted when Footer is nil")
	}
	if strings.Count(bare, "<svg") != 1 {
		t.Error("expected only the chevron svg when Icon is nil")
	}
	if strings.Contains(bare, "ml-025") {
		t.Error("expected summary span omitted when Summary and SummaryAttrs are empty")
	}

	full := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		Label:   "Circles",
		Name:    "circle",
		Options: checkboxFilterOptions(),
		Icon:    checkboxFilterStatic(`<svg class="icon" data-lead></svg>`),
		Footer:  checkboxFilterStatic(`<a href="/clear">Clear all</a>`),
	}))
	if !strings.Contains(full, `<summary class="trigger button--subtle button--small"><svg class="icon" data-lead></svg><span>Circles</span>`) {
		t.Error("expected icon rendered before the label")
	}
	if !strings.Contains(full, `<footer><a href="/clear">Clear all</a></footer>`) {
		t.Error("expected classless footer with content when FooterClass is empty")
	}
}

func TestCheckboxFilterAttrsSpread(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		Label:        "Circles",
		Name:         "circle",
		Options:      checkboxFilterOptions()[:1],
		Attrs:        templ.Attributes{"data-testid": "filter"},
		SummaryAttrs: templ.Attributes{"data-text": "$count"},
		InputAttrs:   templ.Attributes{"data-on:change": "go()", "name": "ignored"},
	}))
	if !strings.Contains(testutil.RawTag(t, html, "details"), `data-testid="filter"`) {
		t.Error("expected Attrs on details")
	}
	nodes := testutil.ParseFragment(t, html)
	if muted := testutil.FindElementByClass(nodes, "ml-025"); muted == nil {
		t.Error("expected summary span rendered for SummaryAttrs")
	} else if v, _ := testutil.AttrVal(muted, "data-text"); v != "$count" {
		t.Errorf("summary data-text = %q", v)
	}
	input := testutil.RawTag(t, html, "input")
	if !strings.Contains(input, `data-on:change="go()"`) {
		t.Error("expected InputAttrs on input")
	}
	if strings.Index(input, `name="circle"`) > strings.Index(input, `name="ignored"`) {
		t.Error("expected component name before InputAttrs so the component wins")
	}
}
