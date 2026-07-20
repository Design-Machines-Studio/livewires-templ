package form

import (
	"maps"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

const checkboxGoldenDefault = `<div><label class="checkbox checkbox--error"><input type="checkbox" name="terms" aria-invalid="true" aria-describedby="terms-accepted-hint terms-error" aria-labelledby="terms-accepted-label" checked disabled value="accepted"><span><span id="terms-accepted-label" class="block">Accept terms</span> <span id="terms-accepted-hint" class="hint block">Required to continue</span></span></label> <p id="terms-error" class="error" role="alert">Acceptance is required</p></div>`

func TestCheckboxRenders(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", false))
	if !strings.Contains(html, `class="checkbox"`) {
		t.Error("expected checkbox class")
	}
	if !strings.Contains(html, `type="checkbox"`) {
		t.Error("expected checkbox input type")
	}
	if !strings.Contains(html, `name="agree"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "I agree") {
		t.Error("expected label text")
	}
}

func TestCheckboxChecked(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", true))
	if !strings.Contains(html, "checked") {
		t.Error("expected checked attribute")
	}
}

func TestCheckboxVariantSuccess(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "ok",
		Label:   "Approved",
		Variant: "success",
	}))
	if !strings.Contains(html, "checkbox--success") {
		t.Error("expected checkbox--success class")
	}
}

func TestCheckboxVariantError(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "err",
		Label:   "Rejected",
		Variant: "error",
	}))
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected checkbox--error class")
	}
}

func TestCheckboxWithValue(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxWithValue("color", "red", "Red", false))
	if !strings.Contains(html, `value="red"`) {
		t.Error("expected value attribute")
	}
}

func TestCheckboxWithError(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:  "agree",
		Label: "I agree",
		Error: "You must agree",
	}))
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected checkbox--error variant when Error is set")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="agree-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="agree-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "You must agree") {
		t.Error("expected error message text")
	}
}

func TestCheckboxErrorPreservesExplicitVariant(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "ok",
		Label:   "OK",
		Variant: "warning",
		Error:   "Problem",
	}))
	if !strings.Contains(html, "checkbox--warning") {
		t.Error("expected explicit variant preserved when Error is set")
	}
}

func TestCheckboxWrapper(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("x", "X", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestCheckboxSanitizesErrorID(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "terms accepted", Label: "Agree", Error: "Required"}))
	assertDescribedByResolves(t, html, 1)
}

func TestCheckboxWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "You can opt out later",
	}))
	if !strings.Contains(html, `<span id="agree-hint" class="hint block">You can opt out later</span>`) {
		t.Errorf("expected hint span, got %s", html)
	}
	// The hint sits inside the label, so the whole label stays a click target.
	if !strings.Contains(html, `class="hint block">You can opt out later</span></span></label>`) {
		t.Errorf("expected hint nested inside the label, got %s", html)
	}
}

// A wrapping label folds the hint into the accessible name unless the input
// names itself from the label span alone.
func TestCheckboxHintAccessibleAssociation(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "You can opt out later",
	}))
	if !strings.Contains(html, `aria-labelledby="agree-label"`) {
		t.Errorf("expected aria-labelledby, got %s", html)
	}
	if !strings.Contains(html, `aria-describedby="agree-hint"`) {
		t.Errorf("expected aria-describedby, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
	if !strings.Contains(html, `id="agree-label"`) {
		t.Error("aria-labelledby resolves to no element")
	}
}

func TestCheckboxWithoutHintRendersUnchanged(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", false))
	for _, unwanted := range []string{`class="hint block"`, "aria-describedby", "aria-labelledby", "<span"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %s without a hint, got %s", unwanted, html)
		}
	}
	if !strings.Contains(html, `<input type="checkbox" name="agree">I agree`) {
		t.Errorf("expected unchanged markup, got %s", html)
	}
}

func TestCheckboxHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "Opt out later", Error: "Required",
	}))
	if !strings.Contains(html, `aria-describedby="agree-hint agree-error"`) {
		t.Errorf("expected hint then error, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected error variant")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
}

// Checkboxes in a set share a name, so their hint ids must be disambiguated
// by value or the wrong hint gets announced.
func TestCheckboxHintIDsUniqueAcrossASharedName(t *testing.T) {
	a := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "opts", Value: "a", Label: "A", Hint: "First"}))
	b := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "opts", Value: "b", Label: "B", Hint: "Second"}))
	if !strings.Contains(a, `id="opts-a-hint"`) || !strings.Contains(b, `id="opts-b-hint"`) {
		t.Errorf("expected value-disambiguated hint ids, got %s and %s", a, b)
	}
}

func TestCheckboxHintPreservesState(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "opts", Value: "a", Label: "A", Hint: "H",
		Checked: true, Disabled: true, Variant: "large", Class: "extra",
		Attrs: templ.Attributes{"data-testid": "cb"},
	}))
	for _, want := range []string{"checked", "disabled", "checkbox--large", "extra", `data-testid="cb"`, `value="a"`, `class="hint block"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s alongside a hint, got %s", want, html)
		}
	}
}

func TestCheckboxHintEscaping(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: `<b>L</b>`, Hint: `a & "b" <script>alert(1)</script>`,
	}))
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Errorf("expected escaping, got %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") || !strings.Contains(html, "&amp;") {
		t.Errorf("expected escaped entities, got %s", html)
	}
}

func TestCheckboxHintIDSanitized(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "opt group", Value: "a b", Label: "A", Hint: "H",
	}))
	assertDescribedByResolves(t, html, 1)
	if strings.Contains(html, `id="opt group-a b-hint"`) {
		t.Error("hint id must not contain whitespace")
	}
}

// With no visible label there is no span to name the input from, so
// aria-labelledby must be omitted rather than point at nothing.
func TestCheckboxHintWithoutLabelOmitsLabelledby(t *testing.T) {
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Hint: "You can opt out later",
		InputAttrs: templ.Attributes{"aria-label": "I agree"},
	}))
	if strings.Contains(output, "aria-labelledby") {
		t.Errorf("expected no aria-labelledby without a label, got %s", output)
	}
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, ok := testutil.AttrVal(input, "aria-label"); !ok || got != "I agree" {
		t.Errorf("input aria-label = %q, %v; want %q, true", got, ok, "I agree")
	}
	assertDescribedByResolves(t, output, 1)
}

func TestCheckboxInputAttrsRenderOnInput(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
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

func TestCheckboxAriaLabelOnInputElement(t *testing.T) {
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:       "autosave",
		InputAttrs: templ.Attributes{"aria-label": "Enable autosave"},
	}))
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, ok := testutil.AttrVal(input, "aria-label"); !ok || got != "Enable autosave" {
		t.Errorf("input aria-label = %q, %v; want %q, true", got, ok, "Enable autosave")
	}
}

func TestCheckboxInputAttrsExactlyOnce(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
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

func TestCheckboxInputAttrsDoNotMutateCallerMaps(t *testing.T) {
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

	testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "autosave", Attrs: attrs, InputAttrs: inputAttrs}))

	if !maps.Equal(attrs, attrsBefore) {
		t.Errorf("Attrs mutated: got %#v, want %#v", attrs, attrsBefore)
	}
	if !maps.Equal(inputAttrs, inputAttrsBefore) {
		t.Errorf("InputAttrs mutated: got %#v, want %#v", inputAttrs, inputAttrsBefore)
	}
}

func TestCheckboxEmptyInputAttrsUnchangedOutput(t *testing.T) {
	omitted := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "terms", Value: "accepted", Label: "Accept terms", Hint: "Required to continue",
		Checked: true, Disabled: true, Error: "Acceptance is required",
	}))
	nilAttrs := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "terms", Value: "accepted", Label: "Accept terms", Hint: "Required to continue",
		Checked: true, Disabled: true, Error: "Acceptance is required", InputAttrs: nil,
	}))
	emptyAttrs := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "terms", Value: "accepted", Label: "Accept terms", Hint: "Required to continue",
		Checked: true, Disabled: true, Error: "Acceptance is required", InputAttrs: templ.Attributes{},
	}))

	for name, got := range map[string]string{"omitted": omitted, "nil": nilAttrs, "empty": emptyAttrs} {
		if got != checkboxGoldenDefault {
			t.Errorf("%s InputAttrs output changed:\ngot  %q\nwant %q", name, got, checkboxGoldenDefault)
		}
	}
	if omitted != nilAttrs || omitted != emptyAttrs {
		t.Error("omitted, nil, and empty InputAttrs outputs differ")
	}
}

func TestCheckboxHintErrorIDRefsResolveWithInputAttrs(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "autosave", Label: "Autosave", Hint: "Saves changes", Error: "Unavailable", InputAttrs: inputAttrs,
	}))
	nodes := testutil.ParseFragment(t, output)
	input := testutil.FindElement(nodes, "input")
	describedBy, ok := testutil.AttrVal(input, "aria-describedby")
	if !ok {
		t.Fatal("input has no aria-describedby")
	}
	for _, ref := range strings.Fields(describedBy) {
		if testutil.FindElementByID(nodes, ref) == nil {
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
	if got := testutil.NodeText(errorMessage); got != "Unavailable" {
		t.Errorf("error text = %q; want %q", got, "Unavailable")
	}
}

func TestCheckboxInputAttrsCollisionComponentWins(t *testing.T) {
	output := testutil.RenderToString(t, Checkbox(CheckboxProps{
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
