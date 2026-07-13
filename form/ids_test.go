package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func renderRadioHint(t *testing.T) string {
	t.Helper()
	return testutil.RenderToString(t, Radio(RadioProps{Name: "plan", Value: "pro", Label: "Pro", Hint: "Billed annually"}))
}

func renderCheckboxHint(t *testing.T) string {
	t.Helper()
	return testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "agree", Label: "I agree", Hint: "Opt out later"}))
}

func renderSwitchHint(t *testing.T) string {
	t.Helper()
	return testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "notify", Label: "Notifications", Hint: "Weekly digest"}))
}

// assertDescribedByResolves checks that every IDREF in the control's
// aria-describedby names an element that actually exists in the markup. A
// dangling or whitespace-split IDREF means the text is never announced.
func assertDescribedByResolves(t *testing.T, html string, wantRefs int) []string {
	t.Helper()

	const attr = `aria-describedby="`
	i := strings.Index(html, attr)
	if i == -1 {
		if wantRefs == 0 {
			return nil
		}
		t.Fatalf("expected aria-describedby, got %s", html)
	}
	rest := html[i+len(attr):]
	value := rest[:strings.Index(rest, `"`)]

	refs := strings.Fields(value)
	if len(refs) != wantRefs {
		t.Fatalf("expected %d IDREFs in aria-describedby, got %q", wantRefs, value)
	}
	for _, ref := range refs {
		if !strings.Contains(html, `id="`+ref+`"`) {
			t.Errorf("IDREF %q resolves to no element, got %s", ref, html)
		}
	}
	return refs
}

// assertLabelPointsAtControl checks the `for` / `id` pair still matches after
// sanitization, otherwise the label stops labelling its control.
func assertLabelPointsAtControl(t *testing.T, html, name string) {
	t.Helper()

	id := controlID(name)
	if strings.ContainsAny(id, " \t\n") {
		t.Fatalf("control id %q contains whitespace", id)
	}
	if !strings.Contains(html, `for="`+id+`"`) {
		t.Errorf("expected label for=%q, got %s", id, html)
	}
	if !strings.Contains(html, `id="`+id+`"`) {
		t.Errorf("expected control id=%q, got %s", id, html)
	}
}

func TestSanitizeIDPart(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"alphanumeric passes through", "email", "email"},
		{"digits and dashes pass through", "field-2", "field-2"},
		{"space is escaped", "plan tier", "plan_20_tier"},
		{"slash is escaped", "a/b", "a_2f_b"},
		{"brackets are escaped", "user[email]", "user_5b_email_5d_"},
		{"dot is escaped", "a.b", "a_2e_b"},
		{"underscore is escaped so it cannot forge an escape", "a_b", "a_5f_b"},
		{"empty stays empty", "", ""},
		{"unicode is escaped", "café", "caf_e9_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeIDPart(tt.in); got != tt.want {
				t.Errorf("sanitizeIDPart(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// The escape must be injective: inputs that differ must produce ids that
// differ, or two controls on one page end up sharing an id.
func TestSanitizeIDPartIsInjective(t *testing.T) {
	inputs := []string{"a/b", "a b", "a-b", "a_b", "a.b", "ab", "a[b]", "a_2f_b"}
	seen := map[string]string{}
	for _, in := range inputs {
		got := sanitizeIDPart(in)
		if prev, dup := seen[got]; dup {
			t.Errorf("%q and %q both sanitize to %q", prev, in, got)
		}
		seen[got] = in
	}
}

func TestControlIDFallsBackWhenNameIsEmpty(t *testing.T) {
	if got := controlID(""); got != nameFallback {
		t.Errorf("controlID(\"\") = %q, want %q", got, nameFallback)
	}
	// An all-symbol name still sanitizes to something non-empty, so no fallback.
	if got := controlID("!!"); got != "_21__21_" {
		t.Errorf("controlID(%q) = %q", "!!", got)
	}
}

func TestErrorAndHintIDsShareTheSanitizedControlID(t *testing.T) {
	const name = "plan tier"
	base := controlID(name)
	if got, want := errorID(name), base+"-error"; got != want {
		t.Errorf("errorID = %q, want %q", got, want)
	}
	if got, want := hintID(name), base+"-hint"; got != want {
		t.Errorf("hintID = %q, want %q", got, want)
	}
}

func TestDescribedByAttrs(t *testing.T) {
	tests := []struct {
		name string
		hint string
		err  string
		want any // nil means: no attribute at all
	}{
		{"neither emits nothing", "", "", nil},
		{"hint only", "h", "", "email-hint"},
		{"error only", "", "e", "email-error"},
		{"hint precedes error", "h", "e", "email-hint email-error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := describedByAttrs("email", tt.hint, tt.err)
			if tt.want == nil {
				if attrs != nil {
					t.Errorf("expected no attributes, got %v", attrs)
				}
				return
			}
			if got := attrs["aria-describedby"]; got != tt.want {
				t.Errorf("aria-describedby = %v, want %v", got, tt.want)
			}
		})
	}
}

// The hint group is a bare grouping element. Live Wires CSS owns the stacking
// and spacing of the label and hint lines, so this library must not bake layout
// classes into the markup.
func TestWrappedLabelGroupCarriesNoLayoutClasses(t *testing.T) {
	cases := map[string]string{
		"radio":    renderRadioHint(t),
		"checkbox": renderCheckboxHint(t),
		"switch":   renderSwitchHint(t),
	}
	for name, html := range cases {
		t.Run(name, func(t *testing.T) {
			for _, unwanted := range []string{"stack", "stack-compact", "cluster"} {
				if strings.Contains(html, unwanted) {
					t.Errorf("hint group must carry no %q class, got %s", unwanted, html)
				}
			}
			// The group itself stays unstyled; only the hint keeps its class.
			if !strings.Contains(html, "<span><span id=") {
				t.Errorf("expected a bare grouping span, got %s", html)
			}
			if !strings.Contains(html, `class="hint"`) {
				t.Errorf("expected the hint to keep its class, got %s", html)
			}
		})
	}
}
