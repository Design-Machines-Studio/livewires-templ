package component

import (
	"strings"
	"testing"
	"time"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestDropdownPanelDefaults(t *testing.T) {
	html := testutil.RenderToString(t, DropdownPanel("Circles"))
	nodes := testutil.ParseFragment(t, html)
	for _, tag := range []string{"header", "footer"} {
		if testutil.FindElement(nodes, tag) != nil {
			t.Errorf("expected no <%s> when nil", tag)
		}
	}
	if testutil.FindElementByClass(nodes, "text-muted") != nil {
		t.Error("expected no summary span without Summary or SummaryAttrs")
	}
	details := testutil.FindElement(nodes, "details")
	if _, ok := testutil.AttrVal(details, "open"); ok {
		t.Error("expected closed panel by default")
	}
	form := testutil.FindElement(nodes, "form")
	if form == nil {
		t.Fatal("expected form.menu by default")
	}
	if _, ok := testutil.AttrVal(form, "action"); ok {
		t.Error("expected no action when empty")
	}
}

func TestDropdownPanelOptions(t *testing.T) {
	html := testutil.RenderToString(t, DropdownPanelComponent(DropdownPanelProps{
		ID:           "p",
		TriggerID:    "p-trigger",
		Label:        "Circles",
		SummaryAttrs: templ.Attributes{"data-text": "$summary"},
		TriggerClass: "button--subtle",
		Open:         true,
		MenuTag:      "div",
		Method:       "post",
		MenuAttrs:    templ.Attributes{"data-on:keydown": "fn()"},
		BodyAttrs:    templ.Attributes{"id": "p-body"},
		Attrs:        templ.Attributes{"data-x": "1"},
	}))
	nodes := testutil.ParseFragment(t, html)
	details := testutil.FindElementByID(nodes, "p")
	if _, ok := testutil.AttrVal(details, "open"); !ok {
		t.Error("expected open attribute")
	}
	if v, _ := testutil.AttrVal(details, "data-x"); v != "1" {
		t.Error("expected Attrs on details")
	}
	if v, _ := testutil.AttrVal(testutil.FindElementByID(nodes, "p-trigger"), "class"); v != "trigger button button--small button--subtle" {
		t.Errorf("trigger class = %q", v)
	}
	if span := testutil.FindElementByClass(nodes, "text-muted"); span == nil {
		t.Error("expected summary span when SummaryAttrs set")
	}
	if testutil.FindElement(nodes, "form") != nil || strings.Contains(html, `method=`) {
		t.Error("expected div.menu without method when MenuTag is div")
	}
	menu := testutil.FindElementByClass(nodes, "menu")
	if v, _ := testutil.AttrVal(menu, "data-on:keydown"); v != "fn()" {
		t.Error("expected MenuAttrs on .menu")
	}
	if testutil.FindElementByID(nodes, "p-body") == nil {
		t.Error("expected BodyAttrs on .body")
	}
}

func TestDropdownPanelFormAction(t *testing.T) {
	html := testutil.RenderToString(t, DropdownPanelComponent(DropdownPanelProps{Label: "X", Method: "post", Action: "/filter"}))
	form := testutil.FindElement(testutil.ParseFragment(t, html), "form")
	if v, _ := testutil.AttrVal(form, "method"); v != "post" {
		t.Errorf("method = %q", v)
	}
	if v, _ := testutil.AttrVal(form, "action"); v != "/filter" {
		t.Errorf("action = %q", v)
	}
}

func TestCheckboxFilterWrapperAndIDs(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxFilterComponent(CheckboxFilterProps{
		OptionsID:  "opts",
		Label:      "Circles",
		Name:       "circle",
		InputAttrs: templ.Attributes{"data-on:change": "go()", "name": "evil"},
		Options:    []CheckboxFilterOption{{ID: "c-1", Value: "ops", Label: "Ops", Count: 0}},
	}))
	nodes := testutil.ParseFragment(t, html)
	body := testutil.FindElementByID(nodes, "opts")
	if v, _ := testutil.AttrVal(body, "aria-label"); v != "Circles" {
		t.Errorf("group label should default to Label, got %q", v)
	}
	input := testutil.FindElementByID(nodes, "c-1")
	if v, _ := testutil.AttrVal(input, "name"); v != "circle" {
		t.Errorf("component name should win, got %q", v)
	}
	if v, _ := testutil.AttrVal(input, "data-on:change"); v != "go()" {
		t.Error("expected InputAttrs on input")
	}
	if strings.Contains(html, "visually-hidden") {
		t.Error("expected no hidden noun without CountLabel")
	}
	if !strings.Contains(html, `<span class="text-muted">0</span>`) {
		t.Error("expected zero count rendered")
	}

	wrapper := testutil.RenderToString(t, CheckboxFilter("Circles", "circle", circleOptions()))
	if strings.Count(wrapper, `name="circle"`) != 3 {
		t.Error("expected wrapper to render every option")
	}
}

func TestCalendarDefaults(t *testing.T) {
	html := testutil.RenderToString(t, Calendar(day(2027, 1, 15), time.Time{}))
	nodes := testutil.ParseFragment(t, html)
	title := testutil.FindElementByID(nodes, "cal-2027-01")
	if title == nil || testutil.NodeText(title) != "January 2027" {
		t.Fatal("expected default title id and text")
	}
	if strings.Count(html, "<tr>") != 7 || strings.Count(html, "<td") != 42 {
		t.Error("expected header row plus six weeks of seven days")
	}
	if strings.Contains(html, "aria-current") || strings.Contains(html, "aria-pressed") {
		t.Error("expected no today or selection state with zero times")
	}
}

func TestCalendarOptions(t *testing.T) {
	html := testutil.RenderToString(t, CalendarComponent(CalendarProps{
		Month:       day(2027, 2, 1),
		Title:       "Février 2027",
		SundayFirst: true,
		RangeStart:  day(2027, 2, 10),
		HidePrev:    true,
		NextAttrs:   templ.Attributes{"data-on:click": "next()"},
		DayAttrs: func(d time.Time) templ.Attributes {
			return templ.Attributes{"data-date": d.Format("2006-01-02"), "type": "submit"}
		},
	}))
	if !strings.Contains(html, "Février 2027") {
		t.Error("expected Title override")
	}
	if !strings.Contains(html, `<th scope="col" abbr="Sunday">Su</th>`) || strings.Index(html, "Sunday") > strings.Index(html, "Monday") {
		t.Error("expected Sunday first column")
	}
	// Feb 1 2027 is a Monday, so a Sunday-first grid starts on Jan 31. The component's type comes first and wins.
	if !strings.Contains(html, `<td class="outside"><button type="button" tabindex="-1" data-date="2027-01-31"`) {
		t.Error("expected leading outside day with DayAttrs after component attrs")
	}
	if strings.Contains(html, `class="prev"`) || !strings.Contains(html, `data-on:click="next()"`) {
		t.Error("expected prev hidden and NextAttrs applied")
	}
	if strings.Contains(html, "in-range") || !strings.Contains(html, `<td class="range-start"><button type="button" aria-pressed="true"`) {
		t.Error("expected range-start only when RangeEnd is zero")
	}
}

func TestCalendarGroupWrapper(t *testing.T) {
	html := testutil.RenderToString(t, CalendarGroup(CalendarProps{Month: day(2027, 12, 20), ID: "g", Title: "Dec"}, 2))
	nodes := testutil.ParseFragment(t, html)
	if testutil.FindElementByClass(nodes, "calendar-group") == nil {
		t.Fatal("expected calendar-group")
	}
	if strings.Count(html, `class="prev"`) != 1 || strings.Count(html, `class="next"`) != 1 {
		t.Error("expected navigation only on outer edges")
	}
	second := testutil.FindElementByID(nodes, "g-2")
	if second == nil || testutil.NodeText(second) != "January 2028" {
		t.Error("expected second month id suffix and default title across the year boundary")
	}
}

func TestDateFilterOptions(t *testing.T) {
	html := testutil.RenderToString(t, DateFilter("Last week", "range", []PresetOption{{Value: "7d", Label: "Last week", Checked: true}}))
	if !strings.Contains(html, `class="dropdown dropdown--panel date-filter"`) || !strings.Contains(html, `<rect x="3"`) {
		t.Error("expected date-filter class and default calendar icon")
	}
	if strings.Contains(html, "sidebar") {
		t.Error("expected no sidebar by default")
	}

	custom := testutil.RenderToString(t, DateFilterComponent(DateFilterProps{
		Label: "Range", Icon: staticHTML(`<svg class="icon custom"></svg>`), Sidebar: true, BodyClass: "extra", Variant: "end",
	}))
	if strings.Contains(custom, `<rect x="3"`) || !strings.Contains(custom, "icon custom") {
		t.Error("expected custom icon to replace default")
	}
	if !strings.Contains(custom, `class="body sidebar extra"`) || !strings.Contains(custom, "dropdown--end date-filter") {
		t.Error("expected sidebar body and end variant")
	}
}

func TestDateRangeInputsOptions(t *testing.T) {
	html := testutil.RenderToString(t, DateRangeInputsComponent(DateRangeInputsProps{
		StartName: "from", EndName: "to", Type: "datetime-local", StartLabel: "From", EndLabel: "To",
		StartAttrs: templ.Attributes{"min": "2027-01-01T00:00"}, EndAttrs: templ.Attributes{"required": true},
	}))
	if strings.Count(html, `type="datetime-local"`) != 2 || strings.Contains(html, "value=") {
		t.Error("expected datetime-local inputs without empty values")
	}
	if !strings.Contains(html, `aria-label="From" min="2027-01-01T00:00"`) || !strings.Contains(html, `aria-label="To" required`) {
		t.Error("expected label overrides and per-input attrs")
	}
}

func TestPresetsInputAttrs(t *testing.T) {
	html := testutil.RenderToString(t, PresetsComponent(PresetsProps{
		Name: "range", Label: "Periods", InputAttrs: templ.Attributes{"data-on:change": "apply()"},
		Options: []PresetOption{{Value: "7d", Label: "Last week"}},
	}))
	if !strings.Contains(html, `aria-label="Periods"`) || !strings.Contains(html, `value="7d" data-on:change="apply()"`) {
		t.Error("expected label override and InputAttrs after component attrs")
	}
}
