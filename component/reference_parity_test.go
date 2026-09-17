package component

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

// Reference fixtures in testdata/livewires are copied verbatim from
// livewires/public/reference/components (dropdown.html, calendar.html, date-filter.html).

func staticHTML(markup string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, markup)
		return err
	})
}

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func assertMatchesReference(t *testing.T, fixture string, c templ.Component) {
	t.Helper()
	want, err := os.ReadFile("testdata/livewires/" + fixture)
	if err != nil {
		t.Fatal(err)
	}
	got := testutil.CanonicalHTML(t, testutil.RenderToString(t, c))
	if exp := testutil.CanonicalHTML(t, string(want)); got != exp {
		t.Errorf("%s mismatch\n--- got ---\n%s\n--- want ---\n%s", fixture, got, exp)
	}
}

func applyFooter(cancel string) templ.Component {
	return staticHTML(`<button type="reset" class="button button--small">` + cancel + `</button><button type="submit" class="button button--small button--accent">Apply</button>`)
}

func circleOptions() []CheckboxFilterOption {
	return []CheckboxFilterOption{
		{Value: "design", Label: "Design", Count: 12},
		{Value: "engineering", Label: "Engineering", Count: 8},
		{Value: "operations", Label: "Operations", Count: 3},
	}
}

func TestReferenceCheckboxFilter(t *testing.T) {
	assertMatchesReference(t, "checkbox-filter.html", CheckboxFilterComponent(CheckboxFilterProps{
		Label:       "Circles",
		Summary:     "All",
		GroupLabel:  "Filter by circle",
		Name:        "circle",
		CountLabel:  "members",
		Options:     circleOptions(),
		Footer:      applyFooter("Clear all"),
		FooterClass: "cluster cluster-between",
	}))
}

func TestReferenceCheckboxFilterHeader(t *testing.T) {
	assertMatchesReference(t, "checkbox-filter-header.html", CheckboxFilterComponent(CheckboxFilterProps{
		Label:       "Status",
		Summary:     "2",
		GroupLabel:  "Filter by status",
		Name:        "status",
		CountLabel:  "items",
		Variant:     "end",
		Header:      staticHTML(`<strong class="text-sm">Status</strong><button type="reset" class="button button--small">Clear</button>`),
		HeaderClass: "cluster cluster-between",
		Options: []CheckboxFilterOption{
			{Value: "draft", Label: "Draft", Count: 4, Checked: true},
			{Value: "active", Label: "Active", Count: 21, Checked: true},
			{Value: "archived", Label: "Archived", Count: 7},
		},
		Footer:      staticHTML(`<button type="submit" class="button button--small button--accent">Apply</button>`),
		FooterClass: "cluster cluster-end",
	}))
}

func janMarked() []time.Time {
	return []time.Time{day(2027, 1, 1), day(2027, 1, 30), day(2027, 2, 4)}
}

func TestReferenceCalendarBasic(t *testing.T) {
	assertMatchesReference(t, "calendar-basic.html", CalendarComponent(CalendarProps{
		Month: day(2027, 1, 1), ID: "m-basic", Today: day(2027, 1, 26), Marked: janMarked(),
	}))
}

func TestReferenceCalendarSingle(t *testing.T) {
	assertMatchesReference(t, "calendar-single.html", CalendarComponent(CalendarProps{
		Month: day(2027, 1, 1), ID: "m-single", Today: day(2027, 1, 26), Selected: day(2027, 1, 14),
	}))
}

func TestReferenceCalendarRange(t *testing.T) {
	assertMatchesReference(t, "calendar-range.html", CalendarComponent(CalendarProps{
		Month: day(2027, 1, 1), ID: "m-range", Today: day(2027, 1, 26), Marked: janMarked(),
		RangeStart: day(2027, 1, 8), RangeEnd: day(2027, 1, 14),
	}))
}

func TestReferenceCalendarCompact(t *testing.T) {
	assertMatchesReference(t, "calendar-compact.html", CalendarComponent(CalendarProps{
		Month: day(2027, 1, 1), ID: "m-compact", Today: day(2027, 1, 26), Selected: day(2027, 1, 14), Variant: "compact",
	}))
}

func groupMonths(ids [2]string, marked []time.Time) templ.Component {
	props := CalendarProps{Today: day(2027, 1, 26), Marked: marked, RangeStart: day(2027, 1, 8), RangeEnd: day(2027, 1, 14)}
	jan, feb := props, props
	jan.Month, jan.ID, jan.HideNext = day(2027, 1, 1), ids[0], true
	feb.Month, feb.ID, feb.HidePrev = day(2027, 2, 1), ids[1], true
	return templ.Join(CalendarComponent(jan), CalendarComponent(feb))
}

func TestReferenceCalendarGroup(t *testing.T) {
	marked := append(janMarked(), day(2027, 2, 8), day(2027, 2, 14))
	assertMatchesReference(t, "calendar-group.html", withChildren(CalendarGroupComponent(CalendarGroupProps{}), groupMonths([2]string{"m-g1", "m-g2"}, marked)))
}

func TestReferenceDateRange(t *testing.T) {
	assertMatchesReference(t, "date-range.html", DateRangeInputs("a-from", "a-to", "2027-01-08", "2027-01-14"))
	assertMatchesReference(t, "date-range-small.html", DateRangeInputsComponent(DateRangeInputsProps{
		StartName: "b-from", EndName: "b-to", StartValue: "2027-01-08", EndValue: "2027-01-14", Size: "small",
	}))
}

func fullPresets(checked string) []PresetOption {
	opts := []PresetOption{
		{Value: "today", Label: "Today"}, {Value: "yesterday", Label: "Yesterday"},
		{Value: "week", Label: "This week"}, {Value: "7d", Label: "Last week"},
		{Value: "month", Label: "This month"}, {Value: "30d", Label: "Last month"},
		{Value: "year", Label: "This year"}, {Value: "12m", Label: "Last 12 months"},
		{Value: "fy", Label: "This fiscal year"}, {Value: "all", Label: "All time"},
	}
	for i := range opts {
		opts[i].Checked = opts[i].Value == checked
	}
	return opts
}

func TestReferencePresets(t *testing.T) {
	assertMatchesReference(t, "presets.html", Presets("p1", fullPresets("7d")))
	assertMatchesReference(t, "presets-inline.html", PresetsComponent(PresetsProps{
		Name: "p2", Inline: true,
		Options: []PresetOption{{Value: "7d", Label: "Last week", Checked: true}, {Value: "30d", Label: "Last month"}, {Value: "12m", Label: "Last year"}},
	}))
}

func withChildren(parent templ.Component, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return parent.Render(templ.WithChildren(ctx, templ.Join(children...)), w)
	})
}

func TestReferenceDateFilterPresetsOnly(t *testing.T) {
	assertMatchesReference(t, "date-filter-presets-only.html", withChildren(
		DateFilterComponent(DateFilterProps{Label: "Last week", Footer: applyFooter("Cancel"), FooterClass: "grid grid-columns-2"}),
		Presets("r1", fullPresets("7d")),
	))
}

func TestReferenceDateFilterCompact(t *testing.T) {
	assertMatchesReference(t, "date-filter-compact-picker.html", withChildren(
		DateFilterComponent(DateFilterProps{Label: "Jan 8 – Jan 14, 2027", Footer: applyFooter("Cancel"), FooterClass: "grid grid-columns-2"}),
		DateRangeInputsComponent(DateRangeInputsProps{StartName: "r2-from", EndName: "r2-to", StartValue: "2027-01-08", EndValue: "2027-01-14", Size: "small"}),
		PresetsComponent(PresetsProps{Name: "r2", Inline: true, Options: []PresetOption{{Value: "7d", Label: "Last week", Checked: true}, {Value: "30d", Label: "Last month"}, {Value: "12m", Label: "Last year"}}}),
		CalendarComponent(CalendarProps{Month: day(2027, 1, 1), ID: "r2-cal", Today: day(2027, 1, 26), Marked: janMarked(), RangeStart: day(2027, 1, 8), RangeEnd: day(2027, 1, 14)}),
	))
}

func TestReferenceDateFilterFull(t *testing.T) {
	marked := append(janMarked(), day(2027, 2, 8), day(2027, 2, 14))
	assertMatchesReference(t, "date-filter-full-picker.html", withChildren(
		DateFilterComponent(DateFilterProps{Label: "Jan 8 – Jan 14, 2027", Sidebar: true, Footer: templ.Join(
			DateRangeInputsComponent(DateRangeInputsProps{StartName: "r3-from", EndName: "r3-to", StartValue: "2027-01-08", EndValue: "2027-01-14", Size: "small"}),
			staticHTML(`<div class="cluster cluster-compact"><button type="reset" class="button button--small">Cancel</button><button type="submit" class="button button--small button--accent">Apply</button></div>`),
		), FooterClass: "cluster cluster-between"}),
		Presets("r3", fullPresets("7d")),
		withChildren(CalendarGroupComponent(CalendarGroupProps{}), groupMonths([2]string{"r3-jan", "r3-feb"}, marked)),
	))
}
