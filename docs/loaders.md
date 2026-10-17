# Hardware loaders

Requires Live Wires **v2.24.0** (`3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3`) CSS containing `src/css/6_components/loader.css`. Canonical SVG paths and semantics match `public/reference/components/loaders.html` in the companion CSS change. No JavaScript is required for the visual.

```go
// Standalone status (a single live region with readable text).
@component.Loader("Loading proposals", "ratchet")
@component.LoaderComponent(component.LoaderProps{Label: "Syncing members", Variant: "line", Mixed: true})

// Measured progress: named progressbar and matching visible percentage.
@component.LoaderComponent(component.LoaderProps{Label: "Exporting", Variant: "torque", Progress: 0.4})

// Busy button retains its readable children. The host chooses Disabled.
@component.ButtonComponent(component.ButtonProps{Busy: true, Disabled: true, Variant: "accent"}) { Saving }
@component.ButtonComponent(component.ButtonProps{Busy: true, Indicator: &component.LoaderIconProps{Variant: "torque", Progress: 0.4}}) { Exporting }

// Optional decorative visuals preserve the region's existing semantics.
@component.LoadingRegionComponent(component.LoadingRegionProps{
    ID: "editor", Label: "Editor", Indicator: component.LoaderIcon("ratchet"),
}) { Editor content }
@component.SaveStatusComponent(component.SaveStatusProps{
    State: "saving", Indicator: component.LoaderIcon("ratchet"),
})
```

`LoaderIconComponent` emits only a decorative SVG (`aria-hidden="true" focusable="false"`). Use ratchet/torque inside buttons or save status; Driver Line belongs only to standalone centered loading surfaces; do not nest the standalone `Loader` status/progress wrapper there. Its props support ratchet, line or torque, mixed line heads, progress, presentation state, Class and developer-controlled Attrs. Invalid icon variants fall back to ratchet. Progress clamps to 0–1, including infinities; NaN becomes zero. Visible and accessible percentages round to at most two decimal places; SVG progress retains its precision. Torque State accepts complete/reset, otherwise pending. State only affects animation; it never proves a save. Do not override component-owned SVG class, viewBox, style, data-state or ARIA through Attrs.

Standalone hardware `LoaderComponent` emits status for ratchet, a centered LoaderPanel for line, and progressbar for torque. Empty and legacy variants retain the original skeleton rendering. LoadingRegion's nil Indicator retains skeletons. An optional indicator is hidden initially and owned by the existing loading controller. SaveStatus emits Indicator only while saving, within its single persistent status region. Busy buttons show an indicator only while Busy; Busy does not disable the button or create a result announcement. Place result announcements outside busy content.

The host owns progress updates, request tracking, local timing, cancellation and confirmed save scope. Existing selectors and lifecycle contracts remain unchanged. CSS controls size, color, motion preferences and button spacing. For speed/theming, set custom properties on an ancestor; Torque reserves its SVG style for measured progress.

## Driver Line: page and container loading

Driver Line never belongs in a button or inline save status. `ButtonProps.Indicator` accepts `*LoaderIconProps`; line is normalized to ratchet. CSS also hides any manually nested Driver Line inside native buttons or `.button`.

```go
// Only Driver Line is visible; status text remains available to assistive technology.
@component.LoaderPanel("Loading content")
@component.LoaderOverlayComponent(component.LoaderOverlayProps{
    Label: "Loading page", Mixed: true,
    Attrs: templ.Attributes{"hidden": !loading},
})

// Keep LoadingRegion's existing single announcement and request controller.
@component.LoadingRegionComponent(component.LoadingRegionProps{
    Indicator: component.LoaderIcon("line"), IndicatorClass: "loader-panel",
}) { Content }
```

`.loader-overlay` covers the viewport, uses the inherited background color at 80% opacity and centers the fully opaque Driver Line. Mount at the page root outside transformed/busy/inert ancestors. The host owns visibility and, for blocking operations, background interaction/inert state. `.loader-panel` fills a sized container with a six-line minimum height. Both honor hidden; use text utilities on Class to change hardware scale. Neither component traps focus, schedules work, or infers completion.

The overlay keeps its outer status region mounted. `Attrs` and `Class` apply to the inner visual surface, so `hidden` never hides the outer live region. Its hidden label has `tabindex="-1"` for application-owned focus. For blocking operations, focus that label, preserve prior focus and inert states, and restore only state the application changed when closing. Keep the wrapper outside busy/inert subtrees.
