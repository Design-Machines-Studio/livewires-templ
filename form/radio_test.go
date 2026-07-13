package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestRadioGroupRenders(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, `class="check-list"`) {
		t.Error("expected check-list class")
	}
	if !strings.Contains(html, `type="radio"`) {
		t.Error("expected radio input type")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("expected list items")
	}
	if !strings.Contains(html, "Option A") {
		t.Error("expected option label text")
	}
}

func TestRadioGroupSelectedChecked(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, "checked") {
		t.Error("expected checked attribute on selected option")
	}
}

func TestRadioGroupInline(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "A"},
		},
		Inline: true,
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, "check-list--inline") {
		t.Error("expected check-list--inline class")
	}
}

func TestRadioSimple(t *testing.T) {
	html := testutil.RenderToString(t, RadioSimple("color", "red", "Red", false))
	if !strings.Contains(html, `class="radio"`) {
		t.Error("expected radio class")
	}
	if !strings.Contains(html, `name="color"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, `value="red"`) {
		t.Error("expected value attribute")
	}
}

func TestRadioVariantSuccess(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:    "status",
		Value:   "ok",
		Label:   "Success",
		Variant: "success",
	}))
	if !strings.Contains(html, "radio--success") {
		t.Error("expected radio--success class")
	}
}

func TestRadioWithError(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "color",
		Value: "red",
		Label: "Red",
		Error: "Pick a color",
	}))
	if !strings.Contains(html, "radio--error") {
		t.Error("expected radio--error variant when Error is set")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="color-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="color-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, "Pick a color") {
		t.Error("expected error message text")
	}
}

func TestRadioWrapper(t *testing.T) {
	html := testutil.RenderToString(t, RadioSimple("x", "a", "A", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestRadioWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan",
		Value: "pro",
		Label: "Pro",
		Hint:  "Billed annually",
	}))
	if !strings.Contains(html, `<span id="plan-pro-hint" class="hint block">Billed annually</span>`) {
		t.Errorf("expected hint span, got %s", html)
	}
	if !strings.Contains(html, `<span id="plan-pro-label" class="block">Pro</span>`) {
		t.Error("expected label span")
	}
	// The hint lives inside the label so the whole label stays a click target.
	if strings.Contains(html, "</label>") && !strings.Contains(html, `class="hint block">Billed annually</span></span></label>`) {
		t.Error("expected hint nested inside the label element")
	}
}

func TestRadioWithoutHint(t *testing.T) {
	html := testutil.RenderToString(t, RadioSimple("plan", "pro", "Pro", false))
	for _, unwanted := range []string{`class="hint block"`, "aria-describedby", "aria-labelledby", "<span"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %s without a hint, got %s", unwanted, html)
		}
	}
}

func TestRadioHintAccessibleAssociation(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan",
		Value: "pro",
		Label: "Pro",
		Hint:  "Billed annually",
	}))
	// aria-describedby must point at the hint span's own id, and aria-labelledby
	// at the label span, so the hint text does not leak into the accessible name.
	if !strings.Contains(html, `aria-describedby="plan-pro-hint"`) {
		t.Error("expected aria-describedby referencing the hint id")
	}
	if !strings.Contains(html, `aria-labelledby="plan-pro-label"`) {
		t.Error("expected aria-labelledby referencing the label id")
	}
	if !strings.Contains(html, `id="plan-pro-hint"`) {
		t.Error("expected the referenced hint id to exist")
	}
	if !strings.Contains(html, `id="plan-pro-label"`) {
		t.Error("expected the referenced label id to exist")
	}
}

func TestRadioGroupPropagatesHints(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name: "plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free", Hint: "No card needed"},
			{Value: "pro", Label: "Pro", Hint: "Billed annually"},
			{Value: "ent", Label: "Enterprise"},
		},
	}))
	if !strings.Contains(html, "No card needed") || !strings.Contains(html, "Billed annually") {
		t.Error("expected each option hint to render")
	}
	if !strings.Contains(html, `aria-describedby="plan-free-hint"`) {
		t.Error("expected free option described by its own hint")
	}
	if !strings.Contains(html, `aria-describedby="plan-pro-hint"`) {
		t.Error("expected pro option described by its own hint")
	}
	// An option with no hint must not gain a hint element or description.
	if !strings.Contains(html, `<input type="radio" name="plan" value="ent">Enterprise`) {
		t.Errorf("expected hintless option to render unchanged, got %s", html)
	}
	if strings.Count(html, `class="hint block"`) != 2 {
		t.Errorf("expected exactly 2 hint elements, got %d", strings.Count(html, `class="hint block"`))
	}
}

func TestRadioGroupHintIDsUnique(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name: "plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free", Hint: "A"},
			{Value: "pro", Label: "Pro", Hint: "B"},
			{Value: "ent", Label: "Enterprise", Hint: "C"},
		},
	}))
	ids := map[string]bool{}
	for _, chunk := range strings.Split(html, `<span id="`)[1:] {
		id := chunk[:strings.Index(chunk, `"`)]
		if ids[id] {
			t.Errorf("duplicate id %q within one radio group", id)
		}
		ids[id] = true
	}
	// 3 options x (label span + hint span)
	if len(ids) != 6 {
		t.Errorf("expected 6 unique span ids, got %d: %v", len(ids), ids)
	}
}

func TestRadioHintIDSanitized(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan tier",
		Value: "opt a/b",
		Label: "Pro",
		Hint:  "Billed annually",
	}))
	// An id may not contain whitespace, and aria-describedby is a space-separated
	// IDREF list, so any space would split the reference into bogus tokens.
	start := strings.Index(html, `aria-describedby="`) + len(`aria-describedby="`)
	described := html[start : start+strings.Index(html[start:], `"`)]
	if strings.ContainsAny(described, " /") {
		t.Errorf("expected sanitized describedby id, got %q", described)
	}
	if !strings.Contains(html, `id="`+described+`"`) {
		t.Errorf("aria-describedby %q has no matching id, got %s", described, html)
	}
}

// A space and a slash must not fold into the same id, or two options in one
// group would end up sharing an id and the wrong hint would be announced.
func TestRadioHintIDsDoNotCollideAcrossSimilarValues(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name: "plan",
		Options: []RadioOption{
			{Value: "a/b", Label: "Slash", Hint: "One"},
			{Value: "a b", Label: "Space", Hint: "Two"},
			{Value: "a-b", Label: "Dash", Hint: "Three"},
		},
	}))
	ids := map[string]bool{}
	for _, chunk := range strings.Split(html, `<span id="`)[1:] {
		id := chunk[:strings.Index(chunk, `"`)]
		if ids[id] {
			t.Errorf("id %q collides across distinct option values", id)
		}
		ids[id] = true
	}
	if len(ids) != 6 {
		t.Errorf("expected 6 unique span ids across 3 options, got %d: %v", len(ids), ids)
	}
}

// aria-describedby is a space-separated IDREF list, so an error id built from an
// unsanitized name containing a space would silently split into two bad tokens.
func TestRadioErrorIDSanitized(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan tier",
		Value: "pro",
		Label: "Pro",
		Error: "Pick a plan",
	}))
	start := strings.Index(html, `aria-describedby="`) + len(`aria-describedby="`)
	described := html[start : start+strings.Index(html[start:], `"`)]
	if strings.Contains(described, " ") {
		t.Errorf("error-only describedby must be a single IDREF, got %q", described)
	}
	if !strings.Contains(html, `<p id="`+described+`"`) {
		t.Errorf("aria-describedby %q has no matching error paragraph id, got %s", described, html)
	}
}

func TestRadioGroupErrorIDSanitized(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name:    "plan tier",
		Options: []RadioOption{{Value: "pro", Label: "Pro"}},
		Error:   "Pick a plan",
	}))
	if strings.Contains(html, `id="plan tier-error"`) {
		t.Errorf("group error id must not contain whitespace, got %s", html)
	}
}

// Hint and error together must both be announced: two valid, distinct IDREFs.
func TestRadioHintAndErrorIDsBothResolve(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan tier",
		Value: "opt a/b",
		Label: "Pro",
		Hint:  "Billed annually",
		Error: "Pick a plan",
	}))
	start := strings.Index(html, `aria-describedby="`) + len(`aria-describedby="`)
	described := html[start : start+strings.Index(html[start:], `"`)]
	refs := strings.Fields(described)
	if len(refs) != 2 {
		t.Fatalf("expected exactly 2 IDREFs (hint, error), got %q", described)
	}
	for _, ref := range refs {
		if !strings.Contains(html, `id="`+ref+`"`) {
			t.Errorf("IDREF %q resolves to no element, got %s", ref, html)
		}
	}
}

func TestRadioHintEmptyValueFallsBackToOption(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan",
		Label: "Pro",
		Hint:  "Billed annually",
	}))
	if !strings.Contains(html, `id="plan-option-hint"`) {
		t.Errorf("expected fallback hint id for an empty value, got %s", html)
	}
}

func TestRadioHintWithError(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan",
		Value: "pro",
		Label: "Pro",
		Hint:  "Billed annually",
		Error: "Pick a plan",
	}))
	if !strings.Contains(html, `aria-describedby="plan-pro-hint plan-error"`) {
		t.Errorf("expected aria-describedby listing hint then error, got %s", html)
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, "radio--error") {
		t.Error("expected radio--error variant")
	}
	if !strings.Contains(html, `<p id="plan-error" class="error" role="alert">Pick a plan</p>`) {
		t.Error("expected error paragraph unchanged")
	}
}

func TestRadioHintPreservesState(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:     "plan",
		Value:    "pro",
		Label:    "Pro",
		Hint:     "Billed annually",
		Checked:  true,
		Disabled: true,
		Variant:  "success",
		Class:    "extra",
		Attrs:    templ.Attributes{"data-testid": "plan-pro"},
	}))
	for _, want := range []string{
		"checked", "disabled", "radio--success", "extra", `data-testid="plan-pro"`, `class="hint block"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s alongside a hint, got %s", want, html)
		}
	}
}

func TestRadioGroupInlineWithHints(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name:   "plan",
		Inline: true,
		Options: []RadioOption{
			{Value: "free", Label: "Free", Hint: "No card needed"},
		},
	}))
	if !strings.Contains(html, "check-list--inline") {
		t.Error("expected check-list--inline class")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("expected list structure preserved")
	}
	if !strings.Contains(html, `class="hint block"`) {
		t.Error("expected hint inside an inline group")
	}
}

func TestRadioHintEscaping(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "plan",
		Value: "pro",
		Label: `<b>Pro</b>`,
		Hint:  `a & "b" <script>alert(1)</script>`,
	}))
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Errorf("expected label and hint to be escaped, got %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") || !strings.Contains(html, "&amp;") {
		t.Errorf("expected escaped entities, got %s", html)
	}
}

func TestRadioGroupWithError(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B"},
		},
		Error: "Please select one",
	}))
	if !strings.Contains(html, "radio--error") {
		t.Error("expected radio--error variant on group radios")
	}
	if !strings.Contains(html, `id="choice-error"`) {
		t.Error("expected group error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "Please select one") {
		t.Error("expected group error message text")
	}
}
