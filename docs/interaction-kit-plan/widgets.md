# Shared widget presentation

The next useful extraction is a shared widget frame, unboxed metric and inline message. Baseplate, Governance and Atlas already repeat these responsibilities. This work does not need Rocket, a new registry or a new dashboard editor.

## Inspected repositories

Inspected fetched `origin/main`, not the active working copies, on 2026-10-09. No product code, service or data changed. Counts are registered source definitions, not the number visible to any particular member.

| Repository | Exact main | Widget system and extraction opportunity |
| --- | --- | --- |
| Baseplate | `102b46792becc927ba777c640f7278a84fbbd6dc` | 14 host definitions: two fixed-order entries, four library entries and eight account/admin overview adapters. Owns registry, authorized rendering, personal/shared layouts, library and editor. Internal body/message helpers are the starting design. |
| Governance | `e91eeafa1325da82e20fc3bd77af373b0a05e588` | Nine optional widgets. `GovernanceWidgetBody` repeats Baseplate's heading, growing content and bottom footer, with extra wrapping. Lists, metrics and deadline/distribution content belong to Governance. |
| Atlas | `313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8` | Twelve optional widgets. Repeated named sections, headings, content and footers across summary, collection and network templates; local CSS adds container sizing and chart treatment. |
| Fixture Jig | `fbabfbb2e32882abbbf326fc1108e803f514b332` | No registered widget contribution found on main. Needs one teaching example once the shared presentation exists. Its authorization/search guide still says contributions are withheld, contradicting the current host API. |
| Demo | `cfc031eeaf4345a477ad20e5e55ea09583229866` | Distribution/data bundle and pinned stock Baseplate, not another rendering implementation. It receives widgets through later approved host/composition publication. |
| Floor | `265d855b67b884474306f0eea77dbefadfbb68ce` | Separate JavaScript development-run viewer, not a Go Fixture or contributor to the host Dashboard. No Templ migration proposed. |
| Legacy Assembly | `7f1d188c0c6de2cb94d0aa9d207d555aa454f7a8` | Older `WidgetProvider`/`Section` interface; no concrete `ContributeDashboard` implementation found on main. Money and other prototype pages are design references, not additional current public-SDK widget systems. |

The authenticated organization inventory identified Governance and Atlas as the current production Fixture repositories. Future Money, Equity, Documents and Calendar systems should consume the shared presentation when they become actual production consumers; do not invent migrations for repositories that do not exist.

## What should be shared

| Primitive | Responsibility | Evidence |
| --- | --- | --- |
| Widget frame/body | Optional title, description, header action, content, footer; explicit heading/instance IDs; full-height composition, wrapping and inherited theme. | Baseplate `DashboardWidgetBody`, Governance `GovernanceWidgetBody`, twelve Atlas sections. |
| Unboxed metric | Accessible label/value/detail with a formatted string value and valid definition-list semantics. No nested card surface. | Governance statistic widgets, Atlas `networkNumberContent`, Baseplate overview totals. |
| Inline message | Quiet icon/message/detail/action composition for caller-chosen empty/unavailable/informational states. | Baseplate `WidgetMessage`, repeated Governance and Atlas state paragraphs. |

Names are provisional. Follow the kit's Props/Class/Attrs/full-component/convenience-wrapper convention. Take icons and optional rich regions as Templ components. Keep heading/body/footer attribute seams small and justified by real callers. Use instance-qualified IDs so grid/library previews and multiple rendered examples cannot cross-label each other.

Existing Card and StatCard are not exact replacements. StatCard owns an outer card surface and emits dt/dd content intended for a definition-list group; inserting it unchanged inside a widget body risks nested chrome or invalid semantics. Extract/reuse a small metric responsibility without breaking current callers. Timeline, ProgressBar, StatusIndicator and existing table primitives remain useful inside the new frame. No special widget variants of every component are needed.

CSS must land in Live Wires first, followed by matching CSS-free Templ markup. The source implementations use utility-heavy wrappers, while this kit's current contract prohibits baking layout primitives into new components. Settle semantic reference hooks and put spacing/layout decisions in CSS; do not blindly copy the older class lists. Any visible normalization needs before/after acceptance, not a claim of automatic pixel parity.

## What stays with its present owner

- **Baseplate:** `fixture.DashboardWidget`, registry/admission, exact action/resource checks, lifecycle availability, trusted fragment rendering/error isolation, size bounds, widget selection/order, shared defaults, personal overrides and editor controls.
- **Each Fixture:** data queries, permissions, filtering before totals/excerpts, dates, labels, domain states and destination URLs.
- **Atlas/Baseplate chart integration:** ECharts/map configuration, assets, initialization/disposal, resizing and data provenance. The shared frame accepts a chart and its accessible fallback; it does not become a chart engine.

Preserve the existing editor: Add/Remove/move/resize/scheme changes create a draft; only Save/Done editing persists it. Leaving the editor discards the draft. Native forms/details, keyboard controls, inert previews, animation preferences and `data-widget-*` hooks remain host-owned. This is already a user-arranged dashboard, so the earlier generic bento discovery needs to ask about gaps in that existing system, not whether such a consumer exists.

## Specific differences to preserve or deliberately resolve

- Baseplate's generic body is a `div`; Atlas uses named `section` elements; Governance adds wrapping and always renders its footer. A shared API must preserve named sections where useful and allow an absent footer.
- Baseplate places generic body styling in `dashboard-widget-editor.css`, including header separators and inherited border colors. Atlas carries separate widget CSS and repeated stylesheet links. Move only common presentation into the shared stylesheet and remove redundant links/styles only after compatible host assets are available.
- Baseplate has a route-sensitive header exception for the updates button. Keep that application-specific decision outside the shared library; use an explicit presentation option if necessary.
- Account/Admin source cards are reused in dashboards through an `OverviewReader` adapter and a dashboard presentation flag. Test both their original pages and dashboard placement when changing them.
- Atlas's confirmed zero, unavailable, overflow, missing join dates and provisional membership observations must stay distinct. A generic metric must not infer these meanings from a numeric value.
- Governance's deadline wording, response eligibility, private drafts and organization-year totals remain domain logic. Do not turn them into generic component policy.

## Source evidence

- [Baseplate body/message helpers](https://github.com/Design-Machines-Studio/assembly-baseplate/blob/102b46792becc927ba777c640f7278a84fbbd6dc/internal/components/dashboard_widget_body.templ), [editor wrapper](https://github.com/Design-Machines-Studio/assembly-baseplate/blob/102b46792becc927ba777c640f7278a84fbbd6dc/internal/components/dashboard_widget.templ), [ADR-011](https://github.com/Design-Machines-Studio/assembly-baseplate/blob/102b46792becc927ba777c640f7278a84fbbd6dc/docs/adr-011-dashboard-widget-registry.md).
- [Governance widget bodies](https://github.com/Design-Machines-Studio/assembly-governance/blob/e91eeafa1325da82e20fc3bd77af373b0a05e588/pages/workspace/proposals/widgets.templ) and [nine registrations](https://github.com/Design-Machines-Studio/assembly-governance/blob/e91eeafa1325da82e20fc3bd77af373b0a05e588/internal/module/module.go).
- [Atlas registrations](https://github.com/Design-Machines-Studio/assembly-atlas/blob/313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8/internal/module/dashboard_settings.go), [summary bodies](https://github.com/Design-Machines-Studio/assembly-atlas/blob/313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8/pages/workspace/network/dashboard_widgets.templ), [collection bodies](https://github.com/Design-Machines-Studio/assembly-atlas/blob/313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8/pages/workspace/network/dashboard_collection.templ), [network bodies](https://github.com/Design-Machines-Studio/assembly-atlas/blob/313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8/pages/workspace/network/dashboard_network.templ) and [chart slots](https://github.com/Design-Machines-Studio/assembly-atlas/blob/313a1c6b5ccfd5d47d26f36a7323c10f4008c7d8/pages/workspace/network/shared_charts.templ).
- [Jig module](https://github.com/Design-Machines-Studio/assembly-fixture-jig/blob/fbabfbb2e32882abbbf326fc1108e803f514b332/internal/module/module.go) and [stale guide](https://github.com/Design-Machines-Studio/assembly-fixture-jig/blob/fbabfbb2e32882abbbf326fc1108e803f514b332/docs/wiki/Authorization-and-Search.md).

## Coordination and verification limits

Governance widget PRs #85, #86 and #88 are merged, but issue #84 remains open. That contradiction calls for checking its remaining release/browser acceptance, not closing it from source inspection. Atlas widget PRs #41 and #42 are merged. Baseplate's registry/editor contributions are merged through #1010, #1081, #1086 and #1095. These are source facts, not verification of a running composed install.

Current open collision boundaries include Governance #111 (proposal first use) and Atlas #58/#59 (network/chart refresh work). Canonical Baseplate, Governance and Atlas checkouts are on task heads rather than the inspected main. Jig has unrelated dirty files and an old feature checkout. None was changed. The board contains unrelated active/release coordination; no message was treated as authority to clear an instance or dependency.

This pass inspected source and native GitHub state. It did not run builds, render widgets, audit screen readers or verify published dependency assets. Later migration acceptance must include matched populated/empty/unavailable/denied states, light/dark/accent schemes, 320/375/desktop widths, keyboard/zoom, chart fallbacks, editor preview/resize and save/cancel/reload. Pin exact producer references and compatible CSS/component artifacts. Do not equate local replacements, merged source, publication and composed consumer proof.


## Owner-approved metric typography

The owner additionally requested a global `font-metric` utility and optional per-widget size, weight and font width for variable fonts. This extends the initial presentation-only scope deliberately.

- Baseplate `--font-metric` is a validated semantic font-family token, configurable in the existing theme editor and applied through `.font-metric` on the value. Labels and supporting content retain their existing font. Match the utility in Live Wires references before shared Templ emits it automatically on StatCard/WidgetMetric values.
- Per-widget size, weight and font width are optional overrides, separate from grid size. Use ordinary CSS properties for registered axes; preserve a readable static/missing-font fallback and actual @font-face ranges. Optical sizing can stay automatic where supported. Do not assume every variable font has a width axis.
- Baseplate owns preview, save, reset and shared/member persistence. `LayoutSetting` currently stores only identity, enabled state, position, spans and scheme; adding optional typography requires a small validated storage extension with existing-row compatibility. W07 owns that change; the earlier no-migrations restriction still applies to presentation-only tasks.
- Keep defaults inherited and resettable. Metric values may use different shapes: long currency amounts can use a narrower face or smaller size; a short percentage can use a larger size. No number truncation, glyph transforms or continuous auto-fit engine. Demonstrate negative/zero/decimal amounts, currencies, counts and percentages at narrow widths and zoom.

Technical reference: [MDN variable-font guidance](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Fonts/Variable_fonts) maps weight/width/optical sizing to their standard CSS properties and describes font-face ranges. The controls must work with the actual selected font; this research is not proof that every installed font supports every axis.
