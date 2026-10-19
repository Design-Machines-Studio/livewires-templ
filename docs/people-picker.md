# People picker

`form.PeoplePicker` and its companion components provide the semantic HTML
presentation for a people-selection field. They do not store people, resolve
group membership, calculate selection state, authorize roles, make requests, or
save changes.

## Components

| Component | Purpose |
| --- | --- |
| `PeoplePicker` | Fieldset, legend, optional description, and content slot |
| `PeoplePickerFieldset` | Same fieldset from just a legend string, for pickers that need no description or attributes |
| `PeoplePickerGroups` | Optional group-control slot |
| `PeoplePickerTable` | Lined table used for member-selection rows |
| `PeoplePickerRow` | Selection checkbox, identity label, optional detail/avatar, and optional trailing cell |
| `PeoplePickerStatus` | Persistent caller-owned `role="status"` region with `aria-atomic="true"` |
| `PeoplePickerActions` | Optional action slot |

Each component has a `*Component(Props)` entry point and a convenience wrapper.
Props components expose `Class` and developer-controlled `Attrs` fields. The
group, table, and action slots accept children. `PeoplePickerRowProps.Trailing`
accepts a caller-owned `templ.Component`, such as a switch or role field.

The row `ID` is required in use and must be unique in the whole document. The
same stable ID links the selection input to the identity label; the component
adds `-selection-label` and, when detail text exists, `-detail` for the label
and description references. Do not derive IDs from list positions. Use distinct
IDs for each picker instance.

`InputName`, `Value`, `Checked`, and `Disabled` map to native input attributes.
`InputAttrs` are developer-controlled attributes on the selection input.
Generated `id`, `name`, `value`, `checked`, `disabled`, `aria-labelledby`, and
optional `aria-describedby` attributes come first; do not repeat them in
`InputAttrs`. The row `Attrs` map lands on `<tr>`. Other `Attrs` maps land on
their documented wrapper elements. Attribute maps must never contain untrusted
input.

The identity label uses a `for` relationship to the checkbox, so clicking its
visible name or optional avatar operates that checkbox. The accessible name is
the single `aria-labelledby` relationship to the hidden “Select [full name]”
label. Optional detail text has a separate `aria-describedby` relationship,
omitted when there is no detail. The complete name remains in the DOM even when
the caller requests visual truncation. Avatars are hidden from assistive
technology in these already-labelled identity rows: image avatars use
`alt=""` and initials are `aria-hidden`. The row renders the reference avatar
markup itself because `component.Avatar` adds `role="img"` and an `aria-label`
to initials avatars, which would name the person a second time.

The selection cell is not `form.Checkbox`. The reference needs the checkbox `id`
(so the identity label in the next cell can target it), a hidden name label, and
no wrapping `<div>`. `form.Switch`, `form.Checkbox` and other controls work in the
`Trailing` slot, and the render tests cover that composition. `form.Switch`
adds a wrapping `<div>` and `aria-hidden` on its track `<span>`, so the gallery
uses the reference organizer markup instead.

## Native form composition

Applications choose the form method, endpoint, people, IDs, names, values,
selection state, group rules, organizer state, role controls, and status text.
Named group and Clear actions can be ordinary submit buttons with
`formnovalidate`; their application handler owns the meaning of those actions.
Keep unrelated role controls outside the selection label. When group buttons
are toggles, set `aria-pressed` to the strings `"true"` or `"false"`. A bare
boolean attribute renders as an empty value, and `false` omits it.

```templ
<form method="post" action="/spaces/update-members">
	@form.PeoplePicker(form.PeoplePickerProps{Legend: "Members", Description: "Choose people for this space."}) {
		@form.PeoplePickerGroups(form.PeoplePickerGroupsProps{Class: "cluster gap-025", Attrs: templ.Attributes{"role": "group", "aria-label": "Circles"}}) {
			<button class="button button--small" type="submit" name="group_action" value="board" formnovalidate>Board</button>
			<button class="button button--small" type="submit" name="group_action" value="officers" formnovalidate>Officers</button>
		}
		@form.PeoplePickerTable(form.PeoplePickerTableProps{}) {
			<thead><tr>
				<th scope="col"><span class="visually-hidden">Select</span></th>
				<th scope="col">Member</th>
				<th scope="col">Organizer</th>
			</tr></thead>
			<tbody>
				@form.PeoplePickerRow(form.PeoplePickerRowProps{
					ID: "space-ada", InputName: "member_ids", Value: "ada",
					Name: "Ada Okafor", Detail: "Board organizer", Checked: true,
					ShowAvatar: true, AvatarSrc: "/members/ada/avatar",
					IdentityClass: "flex gap-05 min-w-0", TrailingLabel: "Organizer",
					Trailing: organizerSwitch("organizer-ada", "Ada Okafor", true, false),
				})
			</tbody>
		}
		@form.PeoplePickerStatus(form.PeoplePickerStatusProps{Class: "text-sm text-muted"})
		@form.PeoplePickerActions(form.PeoplePickerActionsProps{Class: "cluster cluster-compact"}) {
			<button class="button button--small button--ghost" type="submit" name="picker_action" value="clear" formnovalidate>Clear selection</button>
			<button class="button button--small button--accent" type="submit">Save members</button>
		}
	}
</form>
```

Use the same trailing heading text in each row's `TrailingLabel`; it becomes
the narrow-layout label in `data-th`. A `Trailing` control without a
`TrailingLabel` still renders, with no narrow-layout label. Leave out `TrailingLabel` and `Trailing`
for a two-column picker. Omit the table for an empty state and provide a short
description. Group explanations can be regular content in the group slot.

Disabled controls are omitted from browser form submission. A locked selection
or role must be preserved and checked by the server from authoritative data;
rendering a disabled control does not authorize or protect that value. Use
separate application state for overlapping groups, direct selections,
exclusions, and organizer or role values. Removing one group must not erase a
person selected directly or through another group.

The status region is persistent and atomic. Supply text only for meaningful
caller-owned status changes. Do not put a live region on every row or announce
each render.

## Caller-supplied Live Wires utilities

This library contains no CSS. Callers supply layout utilities through the
component seams:

- Groups: `cluster gap-025`
- Identity label: `flex gap-05 min-w-0`
- Actions: `cluster cluster-compact`
- Optional name truncation: `NameClass: "block truncate"`
- Optional detail truncation: `DetailClass: "truncate"`

The gallery demonstrates these combinations using Live Wires v2.25.0, pinned
to producer commit `aad7e3843898745a2f29833393b20896a204b90e`. Use a version that
contains the People Picker component stylesheet. No CSS or production browser
controller ships from this Go package.

The maintained gallery uses local simulation code for group toggles and Clear.
It does not persist changes. Its role and selection controls remain native
form inputs; its group/Clear buttons stay disabled without JavaScript.
