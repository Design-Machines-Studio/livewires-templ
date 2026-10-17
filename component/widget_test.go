package component

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

// This complete producer document is copied verbatim from Live Wires v2.24.0,
// commit 3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3, public/reference/components/widgets.html.
func TestWidgetPublishedReferenceParity(t *testing.T) {
	b, err := os.ReadFile("testdata/livewires/widgets-reference.html")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	heading := testutil.FindElementByID([]*html.Node{doc}, "membership-heading")
	if heading == nil {
		t.Fatal("missing producer membership example")
	}
	frame := heading.Parent.Parent.Parent
	var expected bytes.Buffer
	if err := html.Render(&expected, frame); err != nil {
		t.Fatal(err)
	}
	actual := withChildren(WidgetComponent(WidgetProps{Title: "Membership at a glance", HeadingID: "membership-heading", HeadingLevel: "h3", Attrs: templ.Attributes{"data-example": "metric"}, Footer: commentFixture(`<a class="more" href="#source-notes">Member directory</a>`)}), withChildren(WidgetMetrics(""),
		WidgetMetricItemComponent(WidgetMetricProps{Label: "Active", Value: "42", Variant: "short"}),
		WidgetMetricItemComponent(WidgetMetricProps{Label: "Probationary", Value: "0", Variant: "short"}),
		WidgetMetricItemComponent(WidgetMetricProps{Label: "Joined in the last 30 days", Value: "3", Variant: "short"})))
	if got, want := testutil.CanonicalHTML(t, testutil.RenderToString(t, actual)), testutil.CanonicalHTML(t, expected.String()); got != want {
		t.Fatalf("published structure differs\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func widgetElement(t *testing.T, output, class string) *html.Node {
	t.Helper()
	n := testutil.FindElementByClass(testutil.ParseFragment(t, output), class)
	if n == nil {
		t.Fatalf("missing .%s in %s", class, output)
	}
	return n
}

func TestWidgetFrameSlotsAttributesAndEscaping(t *testing.T) {
	props := WidgetProps{ID: "overview:a", Title: `<script>Title & label</script>`, HeadingID: "overview:a-title", HeadingLevel: "h3", Description: `<b>description</b>`,
		Action: commentFixture(`<a class="widget-action" href="#details">Details</a>`), Footer: commentFixture(`<a href="#more">More</a>`),
		PlainHeader: true, Class: "scheme-accent", Attrs: templ.Attributes{"data-owner": `a"<&`},
		HeadingAttrs: templ.Attributes{"data-heading": "yes"}, ContentAttrs: templ.Attributes{"data-content": "yes"}, FooterAttrs: templ.Attributes{"data-footer": "yes"}}
	output := testutil.RenderToString(t, withChildren(WidgetComponent(props), commentFixture(`<p>Caller content</p>`)))
	nodes := testutil.ParseFragment(t, output)
	frame := widgetElement(t, output, "widget")
	if frame.Data != "section" {
		t.Fatal(output)
	}
	for key, want := range map[string]string{"id": "overview:a", "aria-labelledby": "overview:a-title", "data-owner": `a"<&`, "class": "widget widget--plain-header scheme-accent"} {
		if got, _ := testutil.AttrVal(frame, key); got != want {
			t.Fatalf("%s: %q != %q", key, got, want)
		}
	}
	heading := testutil.FindElementByID(nodes, props.HeadingID)
	if heading == nil || heading.Data != "h3" || testutil.NodeText(heading) != props.Title {
		t.Fatal(output)
	}
	if testutil.FindElement(nodes, "script") != nil || testutil.FindElement(nodes, "b") != nil {
		t.Fatal("text interpreted as markup", output)
	}
	for class, attr := range map[string]string{"widget-title": "data-heading", "widget-content": "data-content", "widget-footer": "data-footer"} {
		n := widgetElement(t, output, class)
		if v, _ := testutil.AttrVal(n, attr); v != "yes" {
			t.Fatal(output)
		}
		if _, ok := testutil.AttrVal(frame, attr); ok {
			t.Fatal("nested attribute leaked onto frame")
		}
	}
	for _, want := range []string{`class="widget-action"`, "Caller content", `href="#more"`} {
		if !strings.Contains(output, want) {
			t.Fatal(output)
		}
	}
	for _, forbidden := range []string{`role=`, `aria-live=`, `class="stack`, `class="box`} {
		if strings.Contains(output, forbidden) {
			t.Fatal(output)
		}
	}
}

func TestWidgetNamingOptionalHeadersAndStableIDs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		props  WidgetProps
		tag    string
		header bool
	}{
		{"body", WidgetProps{}, "div", false},
		{"blank accessible name", WidgetProps{Title: "  ", HeadingID: "blank", Attrs: templ.Attributes{"aria-label": "  "}}, "div", true},
		{"title without ID", WidgetProps{Title: "Summary"}, "div", true},
		{"no invented heading", WidgetProps{HeadingID: "unused"}, "div", false},
		{"named body", WidgetProps{Attrs: templ.Attributes{"aria-label": "Summary"}}, "section", false},
		{"custom header", WidgetProps{Header: commentFixture(`<h4 id="custom">Custom</h4>`), Attrs: templ.Attributes{"aria-labelledby": "custom"}}, "section", true},
		{"custom replaces title", WidgetProps{Title: "Not emitted", HeadingID: "not-emitted", Header: commentFixture(`<p>Custom</p>`)}, "div", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output := testutil.RenderToString(t, WidgetComponent(tc.props))
			nodes := testutil.ParseFragment(t, output)
			if widgetElement(t, output, "widget").Data != tc.tag || (testutil.FindElement(nodes, "header") != nil) != tc.header {
				t.Fatal(output)
			}
			if testutil.FindElement(nodes, "footer") != nil || strings.Contains(output, `role="status"`) || strings.Contains(output, "Not emitted") {
				t.Fatal(output)
			}
		})
	}
	for _, level := range []string{"h1", "h2", "h3", "h4", "h5", "h6", "", "script"} {
		out := testutil.RenderToString(t, WidgetComponent(WidgetProps{Title: "Summary", HeadingID: "title", HeadingLevel: level}))
		want := level
		if level == "" || level == "script" {
			want = "h2"
		}
		if testutil.FindElementByID(testutil.ParseFragment(t, out), "title").Data != want {
			t.Fatal(out)
		}
	}
	combined := templ.Join(Widget("First", "first-heading"), Widget("Second", "second-heading"))
	one := testutil.RenderToString(t, combined)
	two := testutil.RenderToString(t, combined)
	if one != two || strings.Count(one, `id="first-heading"`) != 1 || strings.Count(one, `id="second-heading"`) != 1 || strings.Count(one, `aria-labelledby=`) != 2 {
		t.Fatal(one, two)
	}
}

func TestWidgetMetricsDefinitionListsAndStatus(t *testing.T) {
	for _, value := range []string{"0", "−12.5%", "€ 9.876.543.210.987,65", "+1,234 people", `<script>& value</script>`} {
		out := testutil.RenderToString(t, WidgetMetricComponent(WidgetMetricProps{Label: "Balance", Value: value, Detail: "As reported"}))
		root := widgetElement(t, out, "widget-metric")
		if root.Data != "dl" {
			t.Fatal(out)
		}
		var tags []string
		for c := root.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				tags = append(tags, c.Data)
			}
		}
		if strings.Join(tags, ",") != "dt,dd,dd" {
			t.Fatal(out)
		}
		if got := testutil.NodeText(widgetElement(t, out, "widget-metric-value")); got != value {
			t.Fatal(got, value)
		}
		if testutil.FindElement(testutil.ParseFragment(t, out), "script") != nil {
			t.Fatal(out)
		}
	}
	group := withChildren(WidgetMetrics("caller"), WidgetMetricItem("Active", "0"), WidgetMetricItemComponent(WidgetMetricProps{Label: "Private", Status: "Permission denied", Value: "99", ValueAttrs: templ.Attributes{"style": "--widget-metric-size: 99rem;"}}))
	out := testutil.RenderToString(t, group)
	root := widgetElement(t, out, "widget-metrics")
	if root.Data != "dl" || strings.Count(out, "<dl") != 1 || strings.Count(out, `class="widget-metric"`) != 2 {
		t.Fatal(out)
	}
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data != "div" {
			t.Fatal(out)
		}
	}
	status := widgetElement(t, out, "widget-metric-status")
	if testutil.NodeText(status) != "Permission denied" || strings.Contains(out, "99") || strings.Count(out, "font-metric") != 1 {
		t.Fatal(out)
	}
	for _, statusText := range []string{"Unavailable", "Not reported", "Permission denied", ""} {
		s := testutil.RenderToString(t, WidgetMetricComponent(WidgetMetricProps{Label: "Observation", Status: statusText}))
		if strings.Contains(s, "font-metric") || strings.Contains(s, "widget-metric-value") {
			t.Fatal(s)
		}
	}
}

func TestWidgetMetricTypographyStaysOnValue(t *testing.T) {
	for _, variant := range []string{"", "short", "long", "unexpected"} {
		out := testutil.RenderToString(t, WidgetMetricComponent(WidgetMetricProps{Label: "Label", Value: "42", Detail: "Detail", Variant: variant, Attrs: templ.Attributes{"data-metric": "outer"}, ValueClass: "caller-value", ValueAttrs: templ.Attributes{"data-value": "only", "style": "--widget-metric-size: 2rem; --widget-metric-weight: 650; --widget-metric-width: 80%;"}}))
		value := widgetElement(t, out, "widget-metric-value")
		class, _ := testutil.AttrVal(value, "class")
		if !strings.Contains(class, "font-metric tnum caller-value") || strings.Count(out, "font-metric") != 1 {
			t.Fatal(out)
		}
		if (variant == "short" || variant == "long") && !strings.Contains(class, "widget-metric-value--"+variant) {
			t.Fatal(out)
		}
		if variant == "unexpected" && strings.Contains(out, variant) {
			t.Fatal(out)
		}
		if v, _ := testutil.AttrVal(value, "data-value"); v != "only" {
			t.Fatal(out)
		}
		if _, ok := testutil.AttrVal(widgetElement(t, out, "widget-metric"), "style"); ok {
			t.Fatal(out)
		}
	}
	out := testutil.RenderToString(t, WidgetMetric("Active", "42"))
	if strings.Contains(out, "style=") || strings.Contains(out, "text-") || strings.Contains(out, "font-bold") {
		t.Fatal(out)
	}
}

func TestWidgetMessageQuietEscapedAndOptional(t *testing.T) {
	out := testutil.RenderToString(t, WidgetMessageComponent(WidgetMessageProps{Message: "<b>Unavailable</b>", Detail: "Try <later>", Icon: commentFixture("✓"), Class: "caller", Attrs: templ.Attributes{"data-message": "quiet"}}))
	icon := widgetElement(t, out, "widget-message-icon")
	if v, _ := testutil.AttrVal(icon, "aria-hidden"); v != "true" {
		t.Fatal(out)
	}
	if testutil.NodeText(widgetElement(t, out, "widget-message-detail")) != "Try <later>" {
		t.Fatal(out)
	}
	for _, bad := range []string{"<b>", "<later>", "role=", "aria-live=", `class="card`} {
		if strings.Contains(out, bad) {
			t.Fatal(out)
		}
	}
	simple := testutil.RenderToString(t, WidgetMessage("Nothing yet"))
	if strings.Contains(simple, "widget-message-icon") || strings.Contains(simple, "widget-message-detail") {
		t.Fatal(simple)
	}
}

func TestWidgetLegacyCompatibility(t *testing.T) {
	card := testutil.RenderToString(t, CardComponent(CardProps{Title: "Kept", Class: "caller", Attrs: templ.Attributes{"data-old": "yes"}}))
	if !strings.Contains(card, `class="card box caller"`) || !strings.Contains(card, `data-old="yes"`) || strings.Contains(card, "widget") {
		t.Fatal(card)
	}
	empty := testutil.RenderToString(t, EmptyState("Nothing yet"))
	if !strings.Contains(empty, "empty-state") || strings.Contains(empty, "widget-message") || strings.Contains(empty, "font-metric") {
		t.Fatal(empty)
	}
	stat := testutil.RenderToString(t, StatCardComponent(StatCardProps{Label: "Members", Value: "0", Detail: "Detail 27", SubDetail: "Version v1.0.0", ValueSize: "4xl", Status: "success", Progress: 50, Href: "/members"}))
	if strings.Count(stat, "font-metric") != 1 || !strings.Contains(stat, `class="value font-metric text-4xl"`) || !strings.Contains(stat, `href="/members"`) || !strings.Contains(stat, `aria-valuenow="50"`) || !strings.Contains(stat, "Version v1.0.0") {
		t.Fatal(stat)
	}
}
