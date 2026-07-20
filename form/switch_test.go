package form

import (
	"maps"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

const switchGoldenDefault = `<div><label class="switch error"><input type="checkbox" role="switch" name="notifications" id="notifications-toggle" value="enabled" checked disabled aria-invalid="true" aria-describedby="notifications-toggle-hint notifications-toggle-error" aria-labelledby="notifications-toggle-label"> <span aria-hidden="true"></span><span><span id="notifications-toggle-label" class="block">Notifications</span> <span id="notifications-toggle-hint" class="hint block">Receive weekly updates</span></span></label> <p id="notifications-toggle-error" class="error" role="alert">Choose a setting</p></div>`

func findSwitchElementByID(nodes []*html.Node, id string) *html.Node {
	for _, node := range nodes {
		if found := findSwitchElementByIDFrom(node, id); found != nil {
			return found
		}
	}
	return nil
}

func findSwitchElementByIDFrom(node *html.Node, id string) *html.Node {
	if value, ok := testutil.AttrVal(node, "id"); node.Type == html.ElementNode && ok && value == id {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findSwitchElementByIDFrom(child, id); found != nil {
			return found
		}
	}
	return nil
}

func switchNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var text strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text.WriteString(switchNodeText(child))
	}
	return text.String()
}

func TestSwitchRenders(t *testing.T) {
	html := testutil.RenderToString(t, Switch("notifications", "Enable notifications", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(html, `role="switch"`) {
		t.Error("expected role=switch on input")
	}
	if !strings.Contains(html, `checked`) {
		t.Error("expected checked attribute")
	}
	if !strings.Contains(html, "Enable notifications") {
		t.Error("expected label text")
	}
}

func TestSwitchUnchecked(t *testing.T) {
	html := testutil.RenderToString(t, Switch("dark-mode", "Dark mode", false))
	if strings.Contains(html, `checked`) {
		t.Error("expected no checked attribute when unchecked")
	}
}

func TestSwitchSmallVariant(t *testing.T) {
	html := testutil.RenderToString(t, SwitchSmall("compact", "Compact", false))
	if !strings.Contains(html, "switch--small") {
		t.Error("expected switch--small class")
	}
}

func TestSwitchDisabled(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:     "feature",
		Label:    "Feature",
		Disabled: true,
	}))
	if !strings.Contains(html, `disabled`) {
		t.Error("expected disabled attribute")
	}
}

func TestSwitchDefaultValue(t *testing.T) {
	html := testutil.RenderToString(t, Switch("toggle", "Toggle", false))
	if !strings.Contains(html, `value="1"`) {
		t.Error("expected default value of 1")
	}
}

func TestSwitchCustomValue(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "toggle",
		Label: "Toggle",
		Value: "on",
	}))
	if !strings.Contains(html, `value="on"`) {
		t.Error("expected custom value")
	}
}

func TestSwitchID(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "x",
		ID:    "my-switch",
		Label: "X",
	}))
	if !strings.Contains(html, `id="my-switch"`) {
		t.Error("expected id attribute")
	}
}

func TestSwitchNoIDWhenEmpty(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "x", Label: "X"}))
	if strings.Contains(html, `id=`) {
		t.Error("expected no id attribute when ID is empty")
	}
}

func TestSwitchCustomClass(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "x",
		Label: "X",
		Class: "my-extra-class",
	}))
	if !strings.Contains(html, "my-extra-class") {
		t.Error("expected custom class in output")
	}
}

func TestSwitchAriaHidden(t *testing.T) {
	html := testutil.RenderToString(t, Switch("x", "X", false))
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden on span")
	}
}

func TestSwitchWithError(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "notify",
		Label: "Notifications",
		Error: "Cannot enable",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class on error message")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="notify-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="notify-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert on error message")
	}
	if !strings.Contains(html, "Cannot enable") {
		t.Error("expected error message text")
	}
}

func TestSwitchWithErrorUsesIDForDescribedby(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "notify",
		ID:    "my-switch",
		Label: "Notifications",
		Error: "Nope",
	}))
	if !strings.Contains(html, `aria-describedby="my-switch-error"`) {
		t.Error("expected aria-describedby to use ID when set")
	}
	if !strings.Contains(html, `id="my-switch-error"`) {
		t.Error("expected error message id to use ID when set")
	}
}

func TestSwitchWrapper(t *testing.T) {
	html := testutil.RenderToString(t, Switch("x", "X", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestSwitchSanitizesErrorID(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "email alerts", Label: "Alerts", Error: "Required"}))
	assertDescribedByResolves(t, html, 1)
}

func TestSwitchWithHint(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "Notifications", Hint: "Weekly digest",
	}))
	if !strings.Contains(html, `<span id="notify-hint" class="hint block">Weekly digest</span>`) {
		t.Errorf("expected hint span, got %s", html)
	}
	if !strings.Contains(html, `aria-labelledby="notify-label"`) {
		t.Errorf("expected aria-labelledby, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
	// The pill span must stay hidden and ahead of the text.
	if !strings.Contains(html, `<span aria-hidden="true"></span><span>`) {
		t.Errorf("expected hint group after the pill, got %s", html)
	}
}

func TestSwitchWithoutHintRendersUnchanged(t *testing.T) {
	html := testutil.RenderToString(t, Switch("notify", "Notifications", true))
	for _, unwanted := range []string{`class="hint block"`, "aria-describedby", "aria-labelledby"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %s without a hint, got %s", unwanted, html)
		}
	}
	if !strings.Contains(html, `<span aria-hidden="true"></span>Notifications`) {
		t.Errorf("expected unchanged markup, got %s", html)
	}
}

// With no visible label there is no span to name the input from, so
// aria-labelledby must be omitted rather than point at nothing -- the caller
// supplies aria-label via InputAttrs.
func TestSwitchHintWithoutLabelOmitsLabelledby(t *testing.T) {
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Hint: "Weekly digest",
		InputAttrs: templ.Attributes{"aria-label": "Notifications"},
	}))
	if strings.Contains(output, "aria-labelledby") {
		t.Errorf("expected no aria-labelledby without a label, got %s", output)
	}
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, ok := testutil.AttrVal(input, "aria-label"); !ok || got != "Notifications" {
		t.Errorf("input aria-label = %q, %v; want %q, true", got, ok, "Notifications")
	}
	assertDescribedByResolves(t, output, 1)
}

func TestSwitchHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "Notifications", Hint: "Weekly digest", Error: "Required",
	}))
	if !strings.Contains(html, `aria-describedby="notify-hint notify-error"`) {
		t.Errorf("expected hint then error, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
	if !strings.Contains(html, `class="switch error"`) {
		t.Errorf("expected error class on the label, got %s", html)
	}
}

// The explicit ID prop wins over Name for id derivation.
func TestSwitchHintUsesExplicitIDAndSanitizesIt(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", ID: "sw 1", Label: "N", Hint: "H", Error: "E",
	}))
	if !strings.Contains(html, `id="sw_20_1"`) {
		t.Errorf("expected sanitized input id, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
}

func TestSwitchHintPreservesState(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "N", Hint: "H", Checked: true, Disabled: true,
		Variant: "small", Class: "extra", Value: "yes",
		Attrs: templ.Attributes{"data-testid": "sw"},
	}))
	for _, want := range []string{"checked", "disabled", "switch--small", "extra", `value="yes"`, `data-testid="sw"`, `role="switch"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s alongside a hint, got %s", want, html)
		}
	}
}

func TestSwitchHintEscaping(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "N", Hint: `a & <script>alert(1)</script>`,
	}))
	if strings.Contains(html, "<script>") {
		t.Errorf("expected escaping, got %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") || !strings.Contains(html, "&amp;") {
		t.Errorf("expected escaped entities, got %s", html)
	}
}

func TestSwitchInputAttrsRenderOnInput(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:       "autosave",
		Attrs:      templ.Attributes{"data-testid": "wrapper"},
		InputAttrs: inputAttrs,
	}))
	nodes := testutil.ParseFragment(t, output)
	label := testutil.FindElement(nodes, "label")
	input := testutil.FindElement(nodes, "input")

	if got, ok := testutil.AttrVal(label, "data-testid"); !ok || got != "wrapper" {
		t.Errorf("label data-testid = %q, %v; want %q, true", got, ok, "wrapper")
	}
	if _, ok := testutil.AttrVal(input, "data-testid"); ok {
		t.Error("input unexpectedly contains wrapper data-testid")
	}
	for key, want := range inputAttrs {
		key = strings.ToLower(key)
		if got, ok := testutil.AttrVal(input, key); !ok || got != want {
			t.Errorf("input %s = %q, %v; want %q, true", key, got, ok, want)
		}
		if _, ok := testutil.AttrVal(label, key); ok {
			t.Errorf("label unexpectedly contains input attribute %s", key)
		}
	}
}

func TestSwitchAriaLabelOnInputElement(t *testing.T) {
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:       "autosave",
		InputAttrs: templ.Attributes{"aria-label": "Enable autosave"},
	}))
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, ok := testutil.AttrVal(input, "aria-label"); !ok || got != "Enable autosave" {
		t.Errorf("input aria-label = %q, %v; want %q, true", got, ok, "Enable autosave")
	}
}

func TestSwitchInputAttrsExactlyOnce(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:       "autosave",
		Attrs:      templ.Attributes{"data-testid": "wrapper"},
		InputAttrs: inputAttrs,
	}))
	inputTag := testutil.RawTag(t, output, "input")
	labelTag := testutil.RawTag(t, output, "label")
	for key := range inputAttrs {
		if got := testutil.CountAttr(inputTag, key); got != 1 {
			t.Errorf("input %s occurrence count = %d; want 1", key, got)
		}
		if got := testutil.CountAttr(labelTag, key); got != 0 {
			t.Errorf("label %s occurrence count = %d; want 0", key, got)
		}
	}
}

func TestSwitchInputAttrsDoNotMutateCallerMaps(t *testing.T) {
	attrs := templ.Attributes{"data-testid": "wrapper"}
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	attrsBefore := maps.Clone(attrs)
	inputAttrsBefore := maps.Clone(inputAttrs)

	testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "autosave", Attrs: attrs, InputAttrs: inputAttrs}))

	if !maps.Equal(attrs, attrsBefore) {
		t.Errorf("Attrs mutated: got %#v, want %#v", attrs, attrsBefore)
	}
	if !maps.Equal(inputAttrs, inputAttrsBefore) {
		t.Errorf("InputAttrs mutated: got %#v, want %#v", inputAttrs, inputAttrsBefore)
	}
}

func TestSwitchEmptyInputAttrsUnchangedOutput(t *testing.T) {
	omitted := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notifications", ID: "notifications-toggle", Label: "Notifications",
		Hint: "Receive weekly updates", Checked: true, Disabled: true,
		Value: "enabled", Error: "Choose a setting",
	}))
	nilAttrs := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notifications", ID: "notifications-toggle", Label: "Notifications",
		Hint: "Receive weekly updates", Checked: true, Disabled: true,
		Value: "enabled", Error: "Choose a setting", InputAttrs: nil,
	}))
	emptyAttrs := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notifications", ID: "notifications-toggle", Label: "Notifications",
		Hint: "Receive weekly updates", Checked: true, Disabled: true,
		Value: "enabled", Error: "Choose a setting", InputAttrs: templ.Attributes{},
	}))

	for name, got := range map[string]string{"omitted": omitted, "nil": nilAttrs, "empty": emptyAttrs} {
		if got != switchGoldenDefault {
			t.Errorf("%s InputAttrs output changed:\ngot  %q\nwant %q", name, got, switchGoldenDefault)
		}
	}
	if omitted != nilAttrs || omitted != emptyAttrs {
		t.Error("omitted, nil, and empty InputAttrs outputs differ")
	}
}

func TestSwitchHintErrorIDRefsResolveWithInputAttrs(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "autosave", Label: "Autosave", Hint: "Saves changes", Error: "Unavailable", InputAttrs: inputAttrs,
	}))
	nodes := testutil.ParseFragment(t, output)
	input := testutil.FindElement(nodes, "input")
	describedBy, ok := testutil.AttrVal(input, "aria-describedby")
	if !ok {
		t.Fatal("input has no aria-describedby")
	}
	for _, ref := range strings.Fields(describedBy) {
		if findSwitchElementByID(nodes, ref) == nil {
			t.Errorf("aria-describedby IDREF %q resolves to no element", ref)
		}
	}
	if got, ok := testutil.AttrVal(input, "aria-invalid"); !ok || got != "true" {
		t.Errorf("input aria-invalid = %q, %v; want true, true", got, ok)
	}
	errorMessage := testutil.FindElement(nodes, "p")
	if got, ok := testutil.AttrVal(errorMessage, "role"); !ok || got != "alert" {
		t.Errorf("error role = %q, %v; want alert, true", got, ok)
	}
	if got := switchNodeText(errorMessage); got != "Unavailable" {
		t.Errorf("error text = %q; want %q", got, "Unavailable")
	}
}

func TestSwitchInputAttrsCollisionComponentWins(t *testing.T) {
	output := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "autosave", Label: "Autosave", Hint: "Saves changes", Error: "Unavailable",
		InputAttrs: templ.Attributes{"aria-invalid": "false", "aria-describedby": "bogus"},
	}))
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, _ := testutil.AttrVal(input, "aria-invalid"); got != "true" {
		t.Errorf("parsed aria-invalid = %q; want component value true", got)
	}
	if got, _ := testutil.AttrVal(input, "aria-describedby"); got != "autosave-hint autosave-error" {
		t.Errorf("parsed aria-describedby = %q; want component IDREF list", got)
	}
	inputTag := testutil.RawTag(t, output, "input")
	for _, key := range []string{"aria-invalid", "aria-describedby"} {
		if got := testutil.CountAttr(inputTag, key); got != 2 {
			t.Errorf("raw input %s occurrence count = %d; want 2", key, got)
		}
	}
}
