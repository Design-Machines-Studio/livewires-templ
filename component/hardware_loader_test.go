package component

import (
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"math"
	"strings"
	"testing"
)

func TestHardwareLoaderDecorativeVariants(t *testing.T) {
	for _, variant := range []string{"ratchet", "line", "torque", "invalid", ""} {
		t.Run(variant, func(t *testing.T) {
			output := testutil.RenderToString(t, LoaderIconComponent(LoaderIconProps{Variant: variant, Mixed: true, Class: "text-lg", Attrs: templ.Attributes{"data-loader": "test"}}))
			nodes := testutil.ParseFragment(t, output)
			svg := testutil.FindElement(nodes, "svg")
			if svg == nil {
				t.Fatalf("SVG missing: %s", output)
			}
			if hidden, _ := testutil.AttrVal(svg, "aria-hidden"); hidden != "true" {
				t.Fatal("icon must be decorative")
			}
			if focus, _ := testutil.AttrVal(svg, "focusable"); focus != "false" {
				t.Fatal("icon must not be focusable")
			}
			if strings.Contains(output, `role=`) || strings.Contains(output, `aria-live=`) {
				t.Fatal("icon must not create another announcement")
			}
			if !strings.Contains(output, `data-loader="test"`) || !strings.Contains(output, "text-lg") {
				t.Fatal("consumer hooks missing")
			}
			count := 1
			if variant == "line" {
				count = 4
			}
			if strings.Count(output, `class="head"`) != count {
				t.Fatalf("head count mismatch: %s", output)
			}
			if variant == "line" && !strings.Contains(output, "--deg: 180deg") {
				t.Fatal("slot head must turn half a rotation")
			}
		})
	}
}

func TestHardwareProgressClampAndSemantics(t *testing.T) {
	tests := []struct {
		value             float64
		progress, percent string
	}{
		{-1, "0", "0"}, {0.4, "0.4", "40"}, {2, "1", "100"},
		{0.29, "0.29", "29"}, {0.12345678, "0.12345678", "12.35"},
		{math.NaN(), "0", "0"}, {math.Inf(1), "1", "100"}, {math.Inf(-1), "0", "0"},
	}
	for _, tc := range tests {
		output := testutil.RenderToString(t, LoaderComponent(LoaderProps{Label: "Exporting", Variant: "torque", Progress: tc.value}))
		for _, want := range []string{`role="progressbar"`, `aria-label="Exporting"`, `aria-valuenow="` + tc.percent + `"`, `--p: ` + tc.progress, `>` + tc.percent + `%<`} {
			if !strings.Contains(output, want) {
				t.Fatalf("missing %s: %s", want, output)
			}
		}
		if strings.Contains(output, `role="status"`) {
			t.Fatal("determinate progress should not nest a live region")
		}
	}
}

func TestHardwareStandaloneAndLegacy(t *testing.T) {
	for _, variant := range []string{"ratchet", "line"} {
		output := testutil.RenderToString(t, Loader("<Loading>", variant))
		if strings.Count(output, `role="status"`) != 1 || !strings.Contains(output, "&lt;Loading&gt;") {
			t.Fatalf("readable escaped status required: %s", output)
		}
	}
	legacy := testutil.RenderToString(t, Loader("Loading", ""))
	if !strings.Contains(legacy, `class="loading"`) || strings.Count(legacy, "skeleton--text") != 2 {
		t.Fatal("legacy skeleton changed")
	}
}

func TestHardwareButtonAndRegionComposition(t *testing.T) {
	output := testutil.RenderToString(t, ButtonComponent(ButtonProps{Busy: true}))
	if !strings.Contains(output, `aria-busy="true"`) || !strings.Contains(output, "loader--ratchet") || strings.Contains(output, `disabled`) {
		t.Fatalf("busy must show hardware without disabling: %s", output)
	}
	idle := testutil.RenderToString(t, ButtonComponent(ButtonProps{Indicator: &LoaderIconProps{Variant: "torque", Progress: 0.4}}))
	if strings.Contains(idle, "<svg") || strings.Contains(idle, "aria-busy") {
		t.Fatal("idle button must not show indicator")
	}
	busy := testutil.RenderToString(t, ButtonComponent(ButtonProps{Busy: true, Disabled: true, Indicator: &LoaderIconProps{Variant: "torque", Progress: 0.4}}))
	if !strings.Contains(busy, "loader--torque") || !strings.Contains(busy, "disabled") {
		t.Fatal("custom busy indicator missing")
	}
	region := testutil.RenderToString(t, LoadingRegionComponent(LoadingRegionProps{Indicator: LoaderIcon("line"), IndicatorClass: "loader-panel"}))
	if strings.Count(region, `role="status"`) != 1 || strings.Contains(region, "skeleton") || !strings.Contains(region, `data-loading-indicator hidden aria-hidden="true"`) {
		t.Fatalf("loading region must keep its single announcement and hidden indicator: %s", region)
	}
	fallback := testutil.RenderToString(t, LoadingRegion("Editor"))
	if !strings.Contains(fallback, "skeleton--text") {
		t.Fatal("region skeleton fallback missing")
	}
	for _, state := range []string{"ready", "dirty", "saving", "saved", "failed", "unknown", "invalid"} {
		output := testutil.RenderToString(t, SaveStatusComponent(SaveStatusProps{State: state, Indicator: LoaderIcon("ratchet")}))
		if strings.Count(output, `role="status"`) != 1 {
			t.Fatal("save region must have exactly one announcement")
		}
		if strings.Contains(output, "<svg") != (state == "saving") {
			t.Fatalf("hardware should appear only while saving: %s", output)
		}
	}
}

func TestHardwareLoaderReferenceGeometry(t *testing.T) {
	for _, variant := range []string{"ratchet", "line", "torque"} {
		props := LoaderIconProps{Variant: variant, Progress: 0.4}
		assertMatchesReference(t, "loader-"+variant+".html", LoaderIconComponent(props))
	}
	assertMatchesReference(t, "loader-mixed.html", LoaderIconComponent(LoaderIconProps{Variant: "line", Mixed: true}))
}

func TestDriverLineNeverInBusyButton(t *testing.T) {
	output := testutil.RenderToString(t, ButtonComponent(ButtonProps{Busy: true, Indicator: &LoaderIconProps{Variant: "line"}}))
	if strings.Contains(output, "loader--line") || !strings.Contains(output, "loader--ratchet") {
		t.Fatalf("Driver Line must fall back to ratchet in a button: %s", output)
	}
}

func TestDriverLineSurfaces(t *testing.T) {
	for _, component := range []templ.Component{
		LoaderPanel("Loading content"), LoaderOverlay("Loading page"), Loader("Loading content", "line"),
	} {
		output := testutil.RenderToString(t, component)
		if strings.Count(output, `role="status"`) != 1 || !strings.Contains(output, "visually-hidden") || !strings.Contains(output, "loader--line") {
			t.Fatalf("surface must show Driver Line with one accessible, visually hidden status: %s", output)
		}
		if strings.Contains(output, "<button") {
			t.Fatal("standalone surfaces must not introduce controls")
		}
	}
	output := testutil.RenderToString(t, LoaderOverlayComponent(LoaderOverlayProps{Label: "Loading page", Mixed: true, Attrs: templ.Attributes{"hidden": true, "id": "page-loading"}}))
	if !strings.Contains(output, "hidden") || !strings.Contains(output, `id="page-loading"`) || !strings.Contains(output, "--deg: 180deg") {
		t.Fatal("overlay visibility, ID or mixed hardware seam missing")
	}
}

func TestDriverLineSurfaceReferenceParity(t *testing.T) {
	assertMatchesReference(t, "loader-panel.html", LoaderPanel("Syncing members"))
	assertMatchesReference(t, "loader-overlay.html", LoaderOverlayComponent(LoaderOverlayProps{
		Label: "Loading page", Mixed: true,
		Attrs: templ.Attributes{"hidden": true, "data-loader-overlay": true},
	}))
}

func TestOverlayKeepsStatusOutsideHiddenSurface(t *testing.T) {
	output := testutil.RenderToString(t, LoaderOverlayComponent(LoaderOverlayProps{Attrs: templ.Attributes{"hidden": true}}))
	nodes := testutil.ParseFragment(t, output)
	status := testutil.FindElement(nodes, "div")
	if role, _ := testutil.AttrVal(status, "role"); role != "status" {
		t.Fatal("outer status missing")
	}
	if _, hidden := testutil.AttrVal(status, "hidden"); hidden {
		t.Fatal("status must stay mounted outside hidden visual")
	}
	if !strings.Contains(output, `class="loader-overlay" hidden`) || !strings.Contains(output, `tabindex="-1"`) {
		t.Fatal("visual visibility and label focus hooks required")
	}
}
