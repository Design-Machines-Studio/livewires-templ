package form

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

// The fixture is copied verbatim from Live Wires v2.29.0 (d407d3edf652d6f4f754987a0ad9f48aaeb1faa4),
// public/reference/components/icon-picker.html. The surrounding form.stack, div.cluster and submit
// button belong to the page, so parity compares only each div.icon-picker root.
const iconPickerReferenceFixture = "testdata/livewires/icon-picker-reference.html"

// iconPickerReferenceIcons is the reference's 16-icon subset with its Phosphor tags.
var iconPickerReferenceIcons = []IconPickerIcon{
	{Name: "users-three", Tags: strings.Fields("user group team community people profiles accounts contacts")},
	{Name: "gavel", Tags: strings.Fields("judge justice legal law court hammer government")},
	{Name: "scales", Tags: strings.Fields("measure balance law justice government")},
	{Name: "handshake", Tags: strings.Fields("emoji deal agreement")},
	{Name: "calendar-check", Tags: strings.Fields("dates times events schedule todo checklist")},
	{Name: "megaphone", Tags: strings.Fields("bullhorn announcements loudspeaker broadcast")},
	{Name: "chat-circle-text", Tags: strings.Fields("send sent messages messaging sms texting comment round bubble")},
	{Name: "file-text", Tags: strings.Fields("documents files save write")},
	{Name: "coins", Tags: strings.Fields("cents change money currency payment paying purchase price sell")},
	{Name: "piggy-bank", Tags: strings.Fields("savings deposit money coins")},
	{Name: "shield-check", Tags: strings.Fields("badge security defense authentication authenticated guarded locked encrypted encryption")},
	{Name: "flag", Tags: strings.Fields("country countries finished completed flags")},
	{Name: "house-line", Tags: strings.Fields("homes buildings places locations")},
	{Name: "leaf", Tags: strings.Fields("plants trees branches leaves nodes green environmental")},
	{Name: "lightbulb", Tags: strings.Fields("flashlight bulbs lighting led energy idea")},
	{Name: "star", Tags: strings.Fields("rate ratings favorites favorited")},
}

func iconPickerReferenceSections() []IconPickerSection {
	return []IconPickerSection{{ID: "all", Heading: "Circle icons", Icons: iconPickerReferenceIcons}}
}

func iconPickerReferenceRoots(t *testing.T) []string {
	t.Helper()
	source, err := os.ReadFile(iconPickerReferenceFixture)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	var roots []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			if _, ok := testutil.AttrVal(n, "data-icon-picker"); ok {
				var rendered bytes.Buffer
				if err := html.Render(&rendered, n); err != nil {
					t.Fatal(err)
				}
				roots = append(roots, rendered.String())
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return roots
}

func TestIconPickerMatchesReference(t *testing.T) {
	roots := iconPickerReferenceRoots(t)
	if len(roots) != 2 {
		t.Fatalf("reference has %d pickers, want 2", len(roots))
	}
	cases := []struct {
		name  string
		index int
		props IconPickerProps
	}{
		{"bordered", 0, IconPickerProps{ID: "ref", Name: "ref_icon", Label: "Circle icon", Sections: iconPickerReferenceSections()}},
		{"ghost", 1, IconPickerProps{ID: "ghost", Name: "ghost_icon", Label: "Circle icon", Ghost: true, Sections: iconPickerReferenceSections()}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := testutil.CanonicalHTML(t, testutil.RenderToString(t, IconPickerComponent(tc.props)))
			want := testutil.CanonicalHTML(t, roots[tc.index])
			if got != want {
				t.Fatalf("%s picker differs from the published reference\ngot:\n%s\nwant:\n%s", tc.name, got, want)
			}
		})
	}
}

func renderIconPicker(t *testing.T, props IconPickerProps) (string, []*html.Node) {
	t.Helper()
	out := testutil.RenderToString(t, IconPickerComponent(props))
	return out, testutil.ParseFragment(t, out)
}

func findAll(nodes []*html.Node, match func(*html.Node) bool) []*html.Node {
	var found []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && match(n) {
			found = append(found, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, n := range nodes {
		walk(n)
	}
	return found
}

func iconPickerRadios(nodes []*html.Node) []*html.Node {
	return findAll(nodes, func(n *html.Node) bool {
		v, _ := testutil.AttrVal(n, "type")
		return n.Data == "input" && v == "radio"
	})
}

func attr(t *testing.T, n *html.Node, key string) string {
	t.Helper()
	if n == nil {
		t.Fatalf("missing element for attribute %q", key)
	}
	v, ok := testutil.AttrVal(n, key)
	if !ok {
		t.Fatalf("<%s> missing %q", n.Data, key)
	}
	return v
}

func hasAttr(n *html.Node, key string) bool {
	_, ok := testutil.AttrVal(n, key)
	return ok
}

func checkedValues(nodes []*html.Node) []string {
	var values []string
	for _, r := range iconPickerRadios(nodes) {
		if hasAttr(r, "checked") {
			v, _ := testutil.AttrVal(r, "value")
			values = append(values, v)
		}
	}
	return values
}

func twoSections() []IconPickerSection {
	return []IconPickerSection{
		{ID: "people", Heading: "People", Icons: []IconPickerIcon{{Name: "users-three"}, {Name: "handshake", Label: "Agreement"}}},
		{Heading: "Office", Icons: []IconPickerIcon{{Name: "file-text"}, {Name: "users-three"}}},
	}
}

func TestIconPickerIDWiring(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{ID: "space", Name: "icon", Label: "Space icon", Sections: twoSections()})

	root := testutil.FindElementByClass(nodes, "icon-picker")
	if got := attr(t, root, "id"); got != "space-picker" {
		t.Errorf("root id = %q", got)
	}
	panel := testutil.FindElementByClass(nodes, "icon-picker-panel")
	trigger := testutil.FindElementByClass(nodes, "icon-picker-trigger")
	if attr(t, trigger, "popovertarget") != attr(t, panel, "id") || attr(t, panel, "id") != "space-panel" {
		t.Errorf("popovertarget %q does not match panel id %q", attr(t, trigger, "popovertarget"), attr(t, panel, "id"))
	}
	if !hasAttr(panel, "popover") {
		t.Error("panel missing popover")
	}
	if hasAttr(trigger, "aria-expanded") {
		t.Error("trigger must not emit aria-expanded; the enhancer owns it")
	}
	if attr(t, trigger, "type") != "button" {
		t.Error("trigger must be type=button")
	}

	fieldset := testutil.FindElementByClass(nodes, "icon-picker-fieldset")
	heading := testutil.FindElementByClass(nodes, "icon-picker-heading")
	if attr(t, fieldset, "aria-labelledby") != attr(t, heading, "id") || attr(t, heading, "id") != "space-heading" {
		t.Errorf("fieldset aria-labelledby %q does not match heading id", attr(t, fieldset, "aria-labelledby"))
	}

	search := testutil.FindElementByID(nodes, "space-search")
	if search == nil || attr(t, search, "type") != "search" || attr(t, search, "autocomplete") != "off" || attr(t, search, "class") != "input--small" {
		t.Fatal("search input wiring is wrong")
	}
	searchLabel := findAll(nodes, func(n *html.Node) bool {
		v, _ := testutil.AttrVal(n, "for")
		return n.Data == "label" && v == "space-search"
	})
	if len(searchLabel) != 1 {
		t.Errorf("want one label for the search input, got %d", len(searchLabel))
	}

	sections := findAll(nodes, func(n *html.Node) bool { return n.Data == "section" && hasAttr(n, "aria-labelledby") })
	wantIDs := []string{"space-people", "space-section-2"}
	if len(sections) != len(wantIDs) {
		t.Fatalf("want %d labelled sections, got %d", len(wantIDs), len(sections))
	}
	for i, section := range sections {
		ref := attr(t, section, "aria-labelledby")
		if ref != wantIDs[i] {
			t.Errorf("section %d aria-labelledby = %q, want %q", i, ref, wantIDs[i])
		}
		if h := testutil.FindElementByID([]*html.Node{section}, ref); h == nil || h.Data != "h4" {
			t.Errorf("section %d heading %q missing or not h4", i, ref)
		}
	}
}

func TestIconPickerIDDefaultsFromNameAndSanitizes(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{Name: "space icon", Sections: []IconPickerSection{{ID: "a b", Heading: "All"}}})
	root := testutil.FindElementByClass(nodes, "icon-picker")
	if got := attr(t, root, "id"); got != "space_20_icon-picker" {
		t.Errorf("root id = %q", got)
	}
	if testutil.FindElementByID(nodes, "space_20_icon-a_20_b") == nil {
		t.Error("section id was not sanitized")
	}
}

func TestIconPickerTwoPickersDoNotCollide(t *testing.T) {
	a := testutil.RenderToString(t, IconPickerComponent(IconPickerProps{ID: "a", Name: "a_icon", Sections: twoSections()}))
	b := testutil.RenderToString(t, IconPickerComponent(IconPickerProps{ID: "b", Name: "b_icon", Sections: twoSections()}))
	nodes := testutil.ParseFragment(t, a+b)
	seen := map[string]bool{}
	for _, n := range findAll(nodes, func(n *html.Node) bool { return hasAttr(n, "id") }) {
		id := attr(t, n, "id")
		if seen[id] {
			t.Errorf("duplicate id %q", id)
		}
		seen[id] = true
	}
}

func TestIconPickerCheckedFollowsValue(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  []string
	}{
		{"empty checks Remove", "", []string{""}},
		{"value checks its tile", "handshake", []string{"handshake"}},
		{"repeated name checks first only", "users-three", []string{"users-three"}},
		{"value outside subset checks nothing", "rocket", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Value: tc.value, Sections: twoSections()})
			got := checkedValues(nodes)
			if strings.Join(got, ",") != strings.Join(tc.want, ",") || len(got) != len(tc.want) {
				t.Errorf("checked = %q, want %q", got, tc.want)
			}
		})
	}

	_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Value: "users-three", Sections: twoSections()})
	first := testutil.FindElementByID(nodes, "p-people")
	radio := iconPickerRadios([]*html.Node{first.Parent})[0]
	if !hasAttr(radio, "checked") {
		t.Error("first users-three tile should be checked")
	}
}

func TestIconPickerRadiosShareName(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "space_icon", Sections: twoSections()})
	radios := iconPickerRadios(nodes)
	if len(radios) != 5 {
		t.Fatalf("want 5 radios (Remove + 4 tiles), got %d", len(radios))
	}
	for _, r := range radios {
		if attr(t, r, "name") != "space_icon" {
			t.Errorf("radio name = %q", attr(t, r, "name"))
		}
		if attr(t, r, "class") != "visually-hidden" {
			t.Errorf("radio class = %q", attr(t, r, "class"))
		}
	}
	remove := testutil.FindElementByClass(nodes, "icon-picker-remove")
	if remove.Data != "label" || attr(t, remove, "class") != "icon-picker-remove button button--ghost button--small" {
		t.Errorf("remove label = <%s class=%q>", remove.Data, attr(t, remove, "class"))
	}
	if r := iconPickerRadios([]*html.Node{remove}); len(r) != 1 || attr(t, r[0], "value") != "" {
		t.Error("Remove must wrap one empty-value radio")
	}
	if strings.TrimSpace(testutil.NodeText(remove)) != "Remove" {
		t.Errorf("Remove text = %q", testutil.NodeText(remove))
	}
}

// Names pass through unchanged: the library does not validate kebab case. The
// server validates submissions against the rendered list.
func TestIconPickerNamePassThroughAndLabelDerivation(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Sections: []IconPickerSection{{Heading: "All", Icons: []IconPickerIcon{
		{Name: "chat-circle-text"},
		{Name: "handshake", Label: "Agreement"},
		{Name: "Not_Kebab"},
	}}}})
	tiles := findAll(nodes, func(n *html.Node) bool { v, _ := testutil.AttrVal(n, "class"); return v == "icon-picker-tile" })
	want := []struct{ value, label string }{{"chat-circle-text", "Chat circle text"}, {"handshake", "Agreement"}, {"Not_Kebab", "Not_Kebab"}}
	for i, tile := range tiles {
		if attr(t, tile, "title") != want[i].label {
			t.Errorf("tile %d title = %q, want %q", i, attr(t, tile, "title"), want[i].label)
		}
		if attr(t, iconPickerRadios([]*html.Node{tile})[0], "value") != want[i].value {
			t.Errorf("tile %d value changed", i)
		}
		hidden := testutil.FindElementByClass([]*html.Node{tile}, "visually-hidden")
		var span *html.Node
		for _, n := range findAll([]*html.Node{tile}, func(n *html.Node) bool { return n.Data == "span" }) {
			span = n
		}
		if span == nil || testutil.NodeText(span) != want[i].label || hidden == nil {
			t.Errorf("tile %d hidden label wrong", i)
		}
	}
}

func TestIconPickerTagsAndSprite(t *testing.T) {
	out, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", SpriteHref: "/static/sprite.svg", PlaceholderIcon: "plus", Sections: []IconPickerSection{{Heading: "All", Icons: []IconPickerIcon{
		{Name: "house", Tags: []string{"home", "building"}},
		{Name: "leaf"},
	}}}})
	tiles := findAll(nodes, func(n *html.Node) bool { v, _ := testutil.AttrVal(n, "class"); return v == "icon-picker-tile" })
	if attr(t, tiles[0], "data-tags") != "home building" {
		t.Errorf("data-tags = %q", attr(t, tiles[0], "data-tags"))
	}
	if hasAttr(tiles[1], "data-tags") {
		t.Error("empty tags must omit data-tags")
	}
	uses := findAll(nodes, func(n *html.Node) bool { return n.Data == "use" })
	var hrefs []string
	for _, u := range uses {
		hrefs = append(hrefs, attr(t, u, "href"))
	}
	if strings.Join(hrefs, " ") != "/static/sprite.svg#plus /static/sprite.svg#house /static/sprite.svg#leaf" {
		t.Errorf("hrefs = %q", hrefs)
	}
	for _, svg := range findAll(nodes, func(n *html.Node) bool { return n.Data == "svg" }) {
		if attr(t, svg, "aria-hidden") != "true" || attr(t, svg, "focusable") != "false" || attr(t, svg, "fill") != "currentColor" || attr(t, svg, "class") != "icon" {
			t.Errorf("svg attributes wrong: %v", svg.Attr)
		}
	}
	if !strings.Contains(out, `href="/static/sprite.svg#house"`) {
		t.Error("sprite href not composed")
	}

	_, defaults := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon"})
	if got := attr(t, findAll(defaults, func(n *html.Node) bool { return n.Data == "use" })[0], "href"); got != "/icons/phosphor-bold.svg#selection-plus" {
		t.Errorf("default placeholder href = %q", got)
	}
}

func TestIconPickerTriggerState(t *testing.T) {
	_, empty := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Label: "Circle icon", Sections: twoSections()})
	trigger := testutil.FindElementByClass(empty, "icon-picker-trigger")
	if attr(t, trigger, "aria-label") != "Choose circle icon" || attr(t, trigger, "data-label-choose") != "Choose circle icon" || attr(t, trigger, "data-label-change") != "Change circle icon" {
		t.Errorf("empty trigger labels: %v", trigger.Attr)
	}
	current := testutil.FindElementByClass(empty, "icon-picker-current")
	placeholder := testutil.FindElementByClass(empty, "icon-picker-placeholder")
	if !hasAttr(current, "hidden") || current.FirstChild != nil || hasAttr(placeholder, "hidden") {
		t.Error("empty value: current must be hidden and empty, placeholder visible")
	}

	_, set := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Label: "Circle icon", Value: "handshake", Sections: twoSections()})
	trigger = testutil.FindElementByClass(set, "icon-picker-trigger")
	if got := attr(t, trigger, "aria-label"); got != "Change circle icon, Agreement" {
		t.Errorf("set trigger aria-label = %q", got)
	}
	current = testutil.FindElementByClass(set, "icon-picker-current")
	placeholder = testutil.FindElementByClass(set, "icon-picker-placeholder")
	if hasAttr(current, "hidden") || !hasAttr(placeholder, "hidden") {
		t.Error("set value: current visible, placeholder hidden")
	}
	if use := findAll([]*html.Node{current}, func(n *html.Node) bool { return n.Data == "use" }); len(use) != 1 || attr(t, use[0], "href") != "/icons/phosphor-bold.svg#handshake" {
		t.Error("current icon must reference the value")
	}

	_, outside := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Value: "paper-plane-tilt", Sections: twoSections()})
	if got := attr(t, testutil.FindElementByClass(outside, "icon-picker-trigger"), "aria-label"); got != "Change icon, Paper plane tilt" {
		t.Errorf("outside-subset aria-label = %q", got)
	}

	_, overrides := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Label: "URL icon", ChooseLabel: "Pick", ChangeLabel: "Swap", Value: "leaf"})
	trigger = testutil.FindElementByClass(overrides, "icon-picker-trigger")
	if attr(t, trigger, "aria-label") != "Swap, Leaf" || attr(t, trigger, "data-label-choose") != "Pick" || attr(t, trigger, "data-label-change") != "Swap" {
		t.Errorf("override labels: %v", trigger.Attr)
	}
	_, acronym := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Label: "URL icon"})
	if got := attr(t, testutil.FindElementByClass(acronym, "icon-picker-trigger"), "aria-label"); got != "Choose URL icon" {
		t.Errorf("acronym label = %q", got)
	}
	_, unlabelled := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon"})
	if testutil.NodeText(testutil.FindElementByClass(unlabelled, "icon-picker-heading")) != "Icon" ||
		attr(t, testutil.FindElementByClass(unlabelled, "icon-picker-trigger"), "aria-label") != "Choose icon" {
		t.Error("unlabelled defaults wrong")
	}
}

func TestIconPickerGhostVariant(t *testing.T) {
	_, bordered := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon"})
	_, ghost := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Ghost: true})
	if got := attr(t, testutil.FindElementByClass(bordered, "icon-picker"), "class"); got != "icon-picker" {
		t.Errorf("bordered root class = %q", got)
	}
	if got := attr(t, testutil.FindElementByClass(bordered, "icon-picker-trigger"), "class"); got != "button button--subtle button--icon icon-picker-trigger" {
		t.Errorf("bordered trigger class = %q", got)
	}
	if got := attr(t, testutil.FindElementByClass(ghost, "icon-picker"), "class"); got != "icon-picker icon-picker--ghost" {
		t.Errorf("ghost root class = %q", got)
	}
	if got := attr(t, testutil.FindElementByClass(ghost, "icon-picker-trigger"), "class"); got != "button button--ghost button--icon icon-picker-trigger" {
		t.Errorf("ghost trigger class = %q", got)
	}
}

func TestIconPickerHeadingLevels(t *testing.T) {
	cases := []struct {
		level          int
		panel, section string
	}{
		{0, "h3", "h4"},
		{2, "h2", "h3"},
		{5, "h5", "h6"},
		{6, "h6", "h6"},
		{9, "h3", "h4"},
		{-1, "h3", "h4"},
	}
	for _, tc := range cases {
		_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", HeadingLevel: tc.level, Sections: twoSections()})
		if got := testutil.FindElementByClass(nodes, "icon-picker-heading").Data; got != tc.panel {
			t.Errorf("level %d panel heading = %s, want %s", tc.level, got, tc.panel)
		}
		for _, h := range findAll(nodes, func(n *html.Node) bool {
			v, _ := testutil.AttrVal(n, "class")
			return v == "icon-picker-section-heading"
		}) {
			if h.Data != tc.section {
				t.Errorf("level %d section heading = %s, want %s", tc.level, h.Data, tc.section)
			}
		}
	}
}

func TestIconPickerHiddenRegions(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Sections: twoSections()})
	for _, class := range []string{"icon-picker-search", "icon-picker-recent", "icon-picker-status"} {
		if !hasAttr(testutil.FindElementByClass(nodes, class), "hidden") {
			t.Errorf(".%s must be server-rendered hidden", class)
		}
	}
	recent := testutil.FindElementByClass(nodes, "icon-picker-recent")
	if grid := testutil.FindElementByClass([]*html.Node{recent}, "icon-picker-grid"); grid == nil || grid.FirstChild != nil {
		t.Error("Recent grid must be present and empty")
	}
	if testutil.NodeText(testutil.FindElementByClass([]*html.Node{recent}, "icon-picker-section-heading")) != "Recent" {
		t.Error("Recent heading text missing")
	}
	statuses := findAll(nodes, func(n *html.Node) bool { v, _ := testutil.AttrVal(n, "role"); return v == "status" })
	if len(statuses) != 1 {
		t.Fatalf("want exactly one status region, got %d", len(statuses))
	}
	status := statuses[0]
	if attr(t, status, "class") != "icon-picker-status" || attr(t, status, "aria-atomic") != "true" || status.Parent.Data != "fieldset" {
		t.Error("status must be p.icon-picker-status[aria-atomic=true] inside the fieldset")
	}
}

func TestIconPickerEscaping(t *testing.T) {
	out := testutil.RenderToString(t, IconPickerComponent(IconPickerProps{
		ID: `x"><b>`, Name: `n"><i>`, Label: `<script>alert(1)</script>`,
		Sections: []IconPickerSection{{ID: `s"`, Heading: `<em>Head</em>`, Icons: []IconPickerIcon{{Name: `a"b`, Label: `<b>L</b>`, Tags: []string{`"t"`, `<u>`}}}}},
	}))
	for _, raw := range []string{"<script>", "<b>", "<i>", "<em>", "<u>"} {
		if strings.Contains(out, raw) {
			t.Errorf("unescaped %q in output", raw)
		}
	}
	nodes := testutil.ParseFragment(t, out)
	tile := testutil.FindElementByClass(nodes, "icon-picker-tile")
	if attr(t, tile, "title") != "<b>L</b>" || attr(t, tile, "data-tags") != `"t" <u>` {
		t.Errorf("escaped tile attributes did not round-trip: %v", tile.Attr)
	}
	if testutil.NodeText(testutil.FindElementByClass(nodes, "icon-picker-heading")) != "<script>alert(1)</script>" {
		t.Error("heading text did not round-trip")
	}
	if got := attr(t, testutil.FindElementByClass(nodes, "icon-picker"), "id"); strings.ContainsAny(got, `"<> `) {
		t.Errorf("root id not sanitized: %q", got)
	}
}

func TestIconPickerClassAndAttrsOnRoot(t *testing.T) {
	_, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Class: "extra", Attrs: templ.Attributes{"data-testid": "picker", "id": "override"}})
	root := testutil.FindElementByClass(nodes, "icon-picker")
	if attr(t, root, "class") != "icon-picker extra" || attr(t, root, "data-testid") != "picker" {
		t.Errorf("root attrs: %v", root.Attr)
	}
	if attr(t, root, "id") != "p-picker" {
		t.Error("component id must win over Attrs id")
	}
	if !hasAttr(root, "data-icon-picker") {
		t.Error("root missing data-icon-picker")
	}
	trigger := testutil.FindElementByClass(nodes, "icon-picker-trigger")
	if hasAttr(trigger, "data-testid") {
		t.Error("Attrs must stay on the root")
	}
}

func TestIconPickerInputAttrsOnEveryRadio(t *testing.T) {
	out, nodes := renderIconPicker(t, IconPickerProps{ID: "p", Name: "icon", Value: "handshake", Sections: twoSections(),
		InputAttrs: templ.Attributes{"data-on-change": "$dirty = true", "name": "hijack", "value": "rocket"}})
	radios := iconPickerRadios(nodes)
	for _, r := range radios {
		if attr(t, r, "data-on-change") != "$dirty = true" {
			t.Errorf("radio %q missing InputAttrs", attr(t, r, "value"))
		}
		if attr(t, r, "name") != "icon" || attr(t, r, "value") == "rocket" {
			t.Error("component attributes must win on collision")
		}
	}
	for _, raw := range strings.Split(out, "<input")[1:] {
		tag := "<input" + raw[:strings.Index(raw, ">")+1]
		if !strings.Contains(tag, `type="radio"`) {
			continue
		}
		if testutil.CountAttr(tag, "data-on-change") != 1 {
			t.Errorf("InputAttrs not spread once: %s", tag)
		}
		if idx := strings.Index(tag, "data-on-change"); idx < strings.Index(tag, `class="visually-hidden"`) {
			t.Errorf("InputAttrs must follow component attributes: %s", tag)
		}
	}
	search := testutil.FindElementByID(nodes, "p-search")
	if hasAttr(search, "data-on-change") {
		t.Error("InputAttrs must not reach the search input")
	}
}

func TestIconPickerWrappers(t *testing.T) {
	bordered := testutil.RenderToString(t, IconPicker("icon", "Circle icon", "leaf", twoSections()))
	full := testutil.RenderToString(t, IconPickerComponent(IconPickerProps{Name: "icon", Label: "Circle icon", Value: "leaf", Sections: twoSections()}))
	if bordered != full {
		t.Error("IconPicker wrapper must match the component")
	}
	ghost := testutil.RenderToString(t, IconPickerGhost("icon", "Circle icon", "", twoSections()))
	if !strings.Contains(ghost, "icon-picker--ghost") || !strings.Contains(ghost, "button--ghost button--icon") {
		t.Error("IconPickerGhost must render the ghost variant")
	}
}

func TestIconPickerNoLayoutClasses(t *testing.T) {
	out := testutil.RenderToString(t, IconPickerComponent(IconPickerProps{ID: "p", Name: "icon", Sections: twoSections()}))
	nodes := testutil.ParseFragment(t, out)
	for _, n := range findAll(nodes, func(n *html.Node) bool { return hasAttr(n, "class") }) {
		for _, c := range strings.Fields(attr(t, n, "class")) {
			switch c {
			case "stack", "stack-compact", "cluster", "grid", "sidebar", "center", "box":
				t.Errorf("layout class %q baked into <%s>", c, n.Data)
			}
		}
	}
}
