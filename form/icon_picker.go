package form

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	lw "github.com/Design-Machines-Studio/livewires-templ"
	"github.com/a-h/templ"
)

// Icon picker defaults match the Live Wires v2.29.0 reference.
const (
	iconPickerDefaultSprite       = "/icons/phosphor-bold.svg"
	iconPickerDefaultPlaceholder  = "selection-plus"
	iconPickerDefaultHeadingLevel = 3
)

// IconPickerIcon is one choosable icon. Name, Label and Tags are
// developer-controlled: render them from a vetted icon list, never from
// request input.
type IconPickerIcon struct {
	Name  string   // Phosphor kebab-case name, such as "users-three". Submitted as the radio value and used as the sprite fragment. Passed through unchanged; templ escapes it.
	Label string   // Human name for the tile title and hidden text. Defaults to Name with hyphens as spaces and the first letter capitalised ("users-three" → "Users three").
	Tags  []string // Search terms for the optional enhancer, joined with spaces into data-tags. Omitted when empty.
}

// IconPickerSection is one headed group of tiles, such as a Phosphor category.
type IconPickerSection struct {
	ID      string           // Id suffix for the section heading, rendered as "<picker id>-<ID>". Defaults to "section-<n>" (1-based).
	Heading string           // Visible section heading; it labels the section.
	Icons   []IconPickerIcon // Tiles in display order.
}

// IconPickerProps configures a render-only Live Wires icon picker: a trigger
// button with a native popover panel of radio tiles. The library emits markup
// only. Live Wires CSS styles it and the optional /dist/icon-picker.js
// enhancer adds search, Recent, keyboard movement and positioning.
//
// The server must validate the submitted value against the same icon list it
// rendered. Sections, Name, Value and every icon name are developer-controlled;
// never build them from untrusted input.
type IconPickerProps struct {
	ID              string              // Id prefix for the root ("<ID>-picker"), panel, heading and search. Defaults to an id-safe form of Name. Sanitized for id use.
	Name            string              // Submitted field name shared by every radio, including Remove. Required.
	Value           string              // Current icon name; empty means no icon and checks Remove. A value outside Sections checks no radio, so the field is omitted on submit.
	Label           string              // Panel heading, such as "Circle icon". Defaults to "Icon". Also drives the default trigger labels.
	ChooseLabel     string              // Trigger label without a value. Defaults to "Choose <label>" ("Choose icon" when Label is empty).
	ChangeLabel     string              // Trigger label prefix with a value; the icon label follows after a comma. Defaults to "Change <label>".
	SpriteHref      string              // Sprite URL without a fragment. Defaults to "/icons/phosphor-bold.svg".
	PlaceholderIcon string              // Sprite name shown when there is no value. Defaults to "selection-plus".
	Ghost           bool                // Borderless tiles (.icon-picker--ghost) and a ghost trigger (.button--ghost instead of .button--subtle).
	HeadingLevel    int                 // Panel heading level, 1–6. Defaults to 3. Section headings use the next level, capped at 6.
	Sections        []IconPickerSection // The rendered subset. The server validates submissions against this same list.
	Class           string              // Additional classes on the root.
	Attrs           templ.Attributes    // Developer-controlled attributes on the root. Never populate from untrusted input.
	InputAttrs      templ.Attributes    // Developer-controlled attributes on every radio, including Remove. Rendered after component-generated attributes; on a name collision the component's attribute wins (HTML first-occurrence). Never populate from untrusted input.
}

// iconPickerBaseID is the sanitized id prefix shared by every id the picker emits.
func iconPickerBaseID(props IconPickerProps) string {
	return controlID(lw.DefaultStr(props.ID, props.Name))
}

func iconPickerRootID(props IconPickerProps) string    { return iconPickerBaseID(props) + "-picker" }
func iconPickerPanelID(props IconPickerProps) string   { return iconPickerBaseID(props) + "-panel" }
func iconPickerHeadingID(props IconPickerProps) string { return iconPickerBaseID(props) + "-heading" }
func iconPickerSearchID(props IconPickerProps) string  { return iconPickerBaseID(props) + "-search" }

func iconPickerSectionID(props IconPickerProps, index int, section IconPickerSection) string {
	return iconPickerBaseID(props) + "-" + sanitizeIDPartOr(section.ID, "section-"+strconv.Itoa(index+1))
}

func iconPickerClass(props IconPickerProps) string {
	return lw.ClassNames("icon-picker", lw.If(props.Ghost, "icon-picker--ghost"), props.Class)
}

func iconPickerTriggerClass(props IconPickerProps) string {
	variant := "button--subtle"
	if props.Ghost {
		variant = "button--ghost"
	}
	return lw.ClassNames("button", variant, "button--icon", "icon-picker-trigger")
}

// iconPickerIconLabel returns the human label, deriving it from the kebab name
// exactly as the reference does.
func iconPickerIconLabel(icon IconPickerIcon) string {
	if icon.Label != "" {
		return icon.Label
	}
	return iconPickerNameLabel(icon.Name)
}

func iconPickerNameLabel(name string) string {
	return lw.CapitalizeFirst(strings.ReplaceAll(name, "-", " "))
}

// iconPickerLowerFirst lowercases the first letter so a heading can follow a
// verb ("Circle icon" → "Choose circle icon"). It leaves acronyms such as "URL"
// unchanged.
func iconPickerLowerFirst(s string) string {
	first, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return s
	}
	if next, _ := utf8.DecodeRuneInString(s[size:]); unicode.IsUpper(next) {
		return s
	}
	return string(unicode.ToLower(first)) + s[size:]
}

func iconPickerHeading(props IconPickerProps) string {
	return lw.DefaultStr(props.Label, "Icon")
}

func iconPickerVerbLabel(override, verb string, props IconPickerProps) string {
	if override != "" {
		return override
	}
	return verb + " " + iconPickerLowerFirst(lw.DefaultStr(props.Label, "icon"))
}

func iconPickerChooseLabel(props IconPickerProps) string {
	return iconPickerVerbLabel(props.ChooseLabel, "Choose", props)
}

func iconPickerChangeLabel(props IconPickerProps) string {
	return iconPickerVerbLabel(props.ChangeLabel, "Change", props)
}

// iconPickerTriggerLabel mirrors the enhancer's trigger name so the server and
// enhanced renders agree: "Choose circle icon" or "Change circle icon, Gavel".
func iconPickerTriggerLabel(props IconPickerProps) string {
	if props.Value == "" {
		return iconPickerChooseLabel(props)
	}
	label := iconPickerNameLabel(props.Value)
	if section, icon, ok := iconPickerCheckedIndex(props); ok {
		label = iconPickerIconLabel(props.Sections[section].Icons[icon])
	}
	return iconPickerChangeLabel(props) + ", " + label
}

// iconPickerCheckedIndex finds the first tile matching Value. Only that tile is
// checked, so a name repeated across sections still checks one radio.
func iconPickerCheckedIndex(props IconPickerProps) (section, icon int, ok bool) {
	if props.Value == "" {
		return 0, 0, false
	}
	for s, sec := range props.Sections {
		for i, ic := range sec.Icons {
			if ic.Name == props.Value {
				return s, i, true
			}
		}
	}
	return 0, 0, false
}

func iconPickerTileChecked(props IconPickerProps, section, icon int) bool {
	s, i, ok := iconPickerCheckedIndex(props)
	return ok && s == section && i == icon
}

func iconPickerHref(props IconPickerProps, name string) string {
	return lw.DefaultStr(props.SpriteHref, iconPickerDefaultSprite) + "#" + name
}

func iconPickerPlaceholderHref(props IconPickerProps) string {
	return iconPickerHref(props, lw.DefaultStr(props.PlaceholderIcon, iconPickerDefaultPlaceholder))
}

func iconPickerTileAttrs(icon IconPickerIcon) templ.Attributes {
	attrs := templ.Attributes{"title": iconPickerIconLabel(icon)}
	if tags := strings.Join(icon.Tags, " "); tags != "" {
		attrs["data-tags"] = tags
	}
	return attrs
}

// iconPickerLevels returns the panel and section heading levels.
func iconPickerLevels(props IconPickerProps) (panel, section int) {
	panel = props.HeadingLevel
	if panel < 1 || panel > 6 {
		panel = iconPickerDefaultHeadingLevel
	}
	return panel, min(panel+1, 6)
}

func iconPickerPanelLevel(props IconPickerProps) int {
	panel, _ := iconPickerLevels(props)
	return panel
}

func iconPickerSectionLevel(props IconPickerProps) int {
	_, section := iconPickerLevels(props)
	return section
}
