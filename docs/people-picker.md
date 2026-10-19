# People picker

`form.PeoplePicker` and its companion components provide the semantic HTML
presentation for a people-selection field. They do not store people, resolve
group membership, calculate selection state, authorize roles, make requests, or
save changes.

## Components

| Component | Purpose |
| --- | --- |
| `PeoplePicker` | Fieldset, legend, optional description, and content slot |
| `PeoplePickerGroups` | Optional group-control slot |
| `PeoplePickerTable` | Lined table used for member-selection rows |
| `PeoplePickerRow` | Selection checkbox, identity label, optional detail/avatar, and optional trailing cell |
| `PeoplePickerStatus` | Persistent caller-owned `role="status"` region with `aria-atomic="true"` |
| `PeoplePickerActions` | Optional action slot |

Each component has a `*Component(Props)` entry point and a convenience wrapper.
Props components expose `Class` and developer-controlled `Attrs` fields. The
group, table, and action slots accept children. `PeoplePickerRowProps.Trailing`
accepts a caller-owned `templ.Component`, such as a switch or role field.

The row `ID` is required in use, must be unique in the whole document, and must be a
valid id token with no whitespace, because `aria-labelledby` and `aria-describedby`
split on spaces. The
same stable ID links the selection input to the identity label; the component
adds `-selection-label` and, when detail text exists, `-detail` for the label
and description references. Do not derive IDs from list positions. Use distinct
IDs for each picker instance.

`Name`, `Value`, `Checked`, and `Disabled` map to native input attributes. `Name` is the
submitted field name, as on `form.Checkbox`. `PersonName` is the person's display name.
`InputAttrs` are developer-controlled attributes on the selection input.
Generated `type`, `id`, `name`, `value`, `checked`, `disabled`, `aria-labelledby`, and
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

Browsers submit a form from Enter in a text field using the first submit button in
source order. When the form starts with group or Clear buttons, a role field's Enter
key triggers that button, not Save. Put the primary Save submit button first in
source order and place it visually with CSS, or keep role fields out of forms whose
first submit button is a group action. Another option is to give the form a hidden
default submit button first. Name each group and Clear button so the server can
tell them apart.

```templ
<form method="post" action="/spaces/update-members">
	@form.PeoplePicker(form.PeoplePickerProps{Legend: "Members"}) {
		@form.PeoplePickerGroups(form.PeoplePickerGroupsProps{Class: "cluster gap-025", Attrs: templ.Attributes{"role": "group", "aria-label": "Circles"}}) {
			<button class="button button--small" type="submit" name="group_action" value="board" aria-pressed="true" formnovalidate>Board</button>
			<button class="button button--small" type="submit" name="group_action" value="officers" aria-pressed="false" formnovalidate>Officers</button>
		}
		<p class="people-picker-description">Choose circles to add everyone in them, then adjust individuals.</p>
		@form.PeoplePickerTable(form.PeoplePickerTableProps{}) {
			<thead><tr>
				<th scope="col"><span class="visually-hidden">Select</span></th>
				<th scope="col">Member</th>
				<th scope="col">Organizer</th>
			</tr></thead>
			<tbody>
				@form.PeoplePickerRow(form.PeoplePickerRowProps{
					ID: "space-ada", Name: "member_ids", Value: "ada",
					PersonName: "Ada Okafor", Detail: "Board organizer", Checked: true,
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

The example has no text fields. If you add a role field, put the primary Save submit
button first in source order, as the warning above describes.

Use the same trailing heading text in each row's `TrailingLabel`; it becomes
the narrow-layout label in `data-th`. A `Trailing` control without a
`TrailingLabel` still renders, with no narrow-layout label. Leave out `TrailingLabel` and `Trailing`
for a two-column picker. Omit the table for an empty state and provide a short
description.

The `Description` prop renders directly after the legend. The reference places a
description that explains the groups, such as why a group is unavailable, after the
group controls. When groups are present, write that
`<p class="people-picker-description">` yourself after the groups slot, as the
example does, and link a disabled group to it with `aria-describedby`.

Set `Detail` on a locked row to say why it is locked, for example "Always included:
super admin". Otherwise a screen reader announces only that the checkbox is
unavailable.

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
- Status: `text-sm text-muted`
- Optional name truncation: `NameClass: "block truncate"`
- Optional detail truncation: `DetailClass: "truncate"`

Truncation is only the `NameClass` and `DetailClass` parameters. When either is set,
the row also puts the full text in a `title` for mouse users. The full name stays in
the DOM and in the hidden "Select [full name]" label. Prefer wrapping for names, which
is the default, and truncate `Detail` instead. Truncated text is hidden from sighted
users who zoom or have low vision, and `title` is only a mouse hint, so never rely on
truncation for text people must read to choose the right person.

The gallery demonstrates these combinations authored against Live Wires v2.25.0
(`aad7e3843898745a2f29833393b20896a204b90e`) and now served from the gallery lock at
v2.26.0 (`d4c197b29b1675bf6cada14d051b7a8e1425fe2f`), which leaves the People Picker CSS
unchanged. Use a version that
contains the People Picker component stylesheet. No CSS or production browser
controller ships from this Go package.

The maintained gallery uses local simulation code for group toggles and Clear,
including the selected-circle detail text. It does not persist changes. Its role and
selection controls remain native form inputs; its group/Clear buttons stay disabled
without JavaScript. The narrow-container examples are static: one box sets
`NameClass` and `DetailClass`, and the other wraps.
