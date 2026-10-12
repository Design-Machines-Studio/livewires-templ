package form

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestComboboxRendersEscapedNativeOptionsAndAccessibleMessages(t *testing.T) {
	html := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "zone", RootID: "zone-choice", Name: "interface_timezone",
		Label: "Timezone <override>", Hint: "Use <em>the default</em> if unsure.",
		Error: "Choose a <valid> timezone.", Required: true,
		Options: []SelectOption{{Value: "UTC", Label: "UTC <default>"}},
	}))

	if strings.Contains(html, "<em>the default</em>") || strings.Contains(html, "<valid>") {
		t.Fatalf("component text must be escaped, got %s", html)
	}
	if !strings.Contains(html, `Timezone &lt;override&gt;`) || !strings.Contains(html, `UTC &lt;default&gt;`) || !strings.Contains(html, `&lt;valid&gt;`) {
		t.Fatalf("expected escaped text content, got %s", html)
	}
	if !strings.Contains(html, `data-combobox`) || !strings.Contains(html, `id="zone-choice"`) {
		t.Fatalf("expected stable combobox wrapper, got %s", html)
	}
	if !strings.Contains(html, `for="zone"`) || !strings.Contains(html, `id="zone"`) {
		t.Fatalf("expected label to target the native select, got %s", html)
	}
	if !strings.Contains(html, `aria-describedby="zone-hint zone-error"`) {
		t.Fatalf("expected hint then error associations, got %s", html)
	}
	if !strings.Contains(html, `aria-invalid="true"`) || !strings.Contains(html, `aria-required="true"`) {
		t.Fatalf("expected native error and required state, got %s", html)
	}
	if !strings.Contains(html, `class="text-xs text-muted" aria-hidden="true">Required</span>`) {
		t.Fatalf("expected the existing required indicator, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
}

func TestComboboxPutsAttrsOnWrapperAndInputAttrsOnNativeSelect(t *testing.T) {
	html := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "zone", RootID: "zone-wrapper", Name: "timezone", Label: "Timezone",
		Form: "preferences", Disabled: true,
		Attrs:      templ.Attributes{"data-wrapper": "wrapper-value"},
		InputAttrs: templ.Attributes{"data-control": "select-value"},
		Options:    []SelectOption{{Value: "UTC", Label: "UTC"}},
	}))
	nodes := testutil.ParseFragment(t, html)
	root := testutil.FindElement(nodes, "div")
	nativeSelect := testutil.FindElement(nodes, "select")
	if root == nil || nativeSelect == nil {
		t.Fatalf("expected wrapper and native select, got %s", html)
	}
	if got, ok := testutil.AttrVal(root, "data-wrapper"); !ok || got != "wrapper-value" {
		t.Errorf("wrapper Attrs landed on the wrong element: %q, %v", got, ok)
	}
	if _, ok := testutil.AttrVal(root, "data-control"); ok {
		t.Error("InputAttrs must not land on the wrapper")
	}
	if got, ok := testutil.AttrVal(nativeSelect, "data-control"); !ok || got != "select-value" {
		t.Errorf("InputAttrs did not land on the native select: %q, %v", got, ok)
	}
	if got, ok := testutil.AttrVal(nativeSelect, "form"); !ok || got != "preferences" {
		t.Errorf("form association missing from native select: %q, %v", got, ok)
	}
	if _, ok := testutil.AttrVal(nativeSelect, "disabled"); !ok {
		t.Error("expected disabled on the native select")
	}
	if got := testutil.CountAttr(testutil.RawTag(t, html, "select"), "name"); got != 1 {
		t.Errorf("native select should have exactly one name attribute, got %d", got)
	}
	if strings.Contains(html, `role="combobox"`) {
		t.Error("the component renders a native fallback; the optional producer owns enhancement semantics")
	}
}

func TestComboboxValueAndEmptyOptionSelection(t *testing.T) {
	valueHTML := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "selected-zone", Name: "zone", Label: "Timezone", Value: "Asia/Makassar",
		Options: []SelectOption{{Value: "UTC", Label: "UTC"}, {Value: "Asia/Makassar", Label: "Makassar"}},
	}))
	valueNodes := testutil.ParseFragment(t, valueHTML)
	nativeSelect := testutil.FindElement(valueNodes, "select")
	if nativeSelect == nil {
		t.Fatal("native select is missing")
	}
	if !strings.Contains(valueHTML, `<option value="Asia/Makassar" selected>`) || strings.Contains(valueHTML, `<option value="UTC" selected>`) {
		t.Fatalf("expected only the requested native option selected, got %s", valueHTML)
	}

	retainedHTML := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "saved-zone", Name: "timezone", Label: "Timezone", Value: "America/Whitehorse",
		SelectedLabel: "Whitehorse", Options: []SelectOption{{Value: "UTC", Label: "UTC"}},
	}))
	if !strings.Contains(retainedHTML, `<option value="America/Whitehorse" selected>Whitehorse</option>`) || strings.Contains(retainedHTML, `<option value="UTC" selected>`) {
		t.Fatalf("expected a saved value and its label to survive when absent from curated options, got %s", retainedHTML)
	}

	emptyHTML := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "optional-zone", Name: "zone", Label: "Timezone",
		Options: []SelectOption{{Value: "", Label: "Use organization default", Selected: true}, {Value: "UTC", Label: "UTC"}},
	}))
	if !strings.Contains(emptyHTML, `<option value="" selected`) {
		t.Fatalf("expected explicit empty default option, got %s", emptyHTML)
	}

	withPlaceholder := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "default-zone", Name: "zone", Label: "Timezone", Placeholder: "Use organization default",
		Options: []SelectOption{{Value: "UTC", Label: "UTC"}, {Value: "", Label: "Use organization default", Selected: true}},
	}))
	if got := strings.Count(withPlaceholder, `<option value=""`); got != 1 {
		t.Fatalf("placeholder must not duplicate an explicit empty option, got %s", withPlaceholder)
	}
	if !strings.Contains(withPlaceholder, `<option value="" selected>Use organization default</option><option value="UTC"`) {
		t.Fatalf("the explicit empty option must be moved before other options, got %s", withPlaceholder)
	}
}

func TestComboboxSupportsRequiredEmptyPlaceholderAndConvenienceModes(t *testing.T) {
	html := testutil.RenderToString(t, ComboboxRemote(ComboboxProps{
		ID: "required-zone", Name: "timezone", Label: "Timezone", Placeholder: "Choose a timezone",
		Required: true, Options: []SelectOption{{Value: "UTC", Label: "UTC"}},
	}))
	if !strings.Contains(html, `data-search`) || !strings.Contains(html, `<option value="">Choose a timezone</option>`) {
		t.Fatalf("expected remote marker and empty placeholder, got %s", html)
	}
	selectTag := testutil.RawTag(t, html, "select")
	if !strings.Contains(selectTag, `required`) {
		t.Fatalf("required must remain on the native select, got %s", selectTag)
	}

	saved := testutil.RenderToString(t, ComboboxComponent(ComboboxProps{
		ID: "required-saved-zone", Name: "timezone", Label: "Timezone", Required: true,
		Placeholder: "Choose a timezone", Value: "America/Whitehorse", SelectedLabel: "Whitehorse",
	}))
	if !strings.Contains(saved, `<option value="">Choose a timezone</option> <option value="America/Whitehorse" selected>Whitehorse</option>`) {
		t.Fatalf("the required empty placeholder must precede the retained selected value: %s", saved)
	}
	if !strings.Contains(testutil.RawTag(t, saved, "select"), `required`) {
		t.Fatalf("required validation must remain on the native select: %s", saved)
	}

	simple := testutil.RenderToString(t, ComboboxSimple("Timezone", "timezone", "Choose a timezone", "UTC", []SelectOption{{Value: "UTC", Label: "UTC"}}, true))
	if !strings.Contains(simple, `name="timezone"`) || !strings.Contains(simple, `required`) || !strings.Contains(simple, `<option value="">Choose a timezone</option>`) || strings.Contains(simple, `data-search`) {
		t.Fatalf("expected simple local convenience wrapper, got %s", simple)
	}
}

func TestComboboxInstancesKeepExplicitIDsDistinct(t *testing.T) {
	first := ComboboxComponent(ComboboxProps{ID: "zone-a", RootID: "choice-a", Name: "zone_a", Label: "First timezone", Hint: "First hint", Error: "First error"})
	second := ComboboxRemote(ComboboxProps{ID: "zone-b", RootID: "choice-b", Name: "zone_b", Label: "Second timezone", Hint: "Second hint", Error: "Second error"})
	combined := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := first.Render(ctx, w); err != nil {
			return err
		}
		return second.Render(ctx, w)
	})
	html := testutil.RenderToString(t, combined)
	for _, id := range []string{"zone-a", "zone-a-hint", "zone-a-error", "choice-a", "zone-b", "zone-b-hint", "zone-b-error", "choice-b"} {
		if strings.Count(html, `id="`+id+`"`) != 1 {
			t.Errorf("expected exactly one id %q in combined render, got %s", id, html)
		}
	}
	if strings.Count(html, `data-combobox`) != 2 || strings.Count(html, `name="zone_a"`) != 1 || strings.Count(html, `name="zone_b"`) != 1 {
		t.Fatalf("expected two independent native controls, got %s", html)
	}
}
