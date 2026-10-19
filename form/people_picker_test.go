package form

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"golang.org/x/net/html"
)

func peoplePickerTestTrailing() templ.Component {
	return peoplePickerStaticHTML(`<label class="switch switch--small"><input type="checkbox" role="switch" aria-label="Make Ada an organizer"><span></span></label>`)
}

func peoplePickerStaticHTML(markup string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, writer io.Writer) error {
		_, err := io.WriteString(writer, markup)
		return err
	})
}

func peoplePickerWithChildren(parent templ.Component, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, writer io.Writer) error {
		return parent.Render(templ.WithChildren(ctx, templ.Join(children...)), writer)
	})
}

func parsePeoplePicker(t *testing.T, markup string) *html.Node {
	t.Helper()
	node, err := html.Parse(strings.NewReader(markup))
	if err != nil {
		t.Fatalf("parse rendered HTML: %v", err)
	}
	return node
}

func findPeoplePickerNodes(node *html.Node, tag string, match func(*html.Node) bool) []*html.Node {
	var found []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag && (match == nil || match(n)) {
			found = append(found, n)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return found
}

func TestPeoplePickerMatchesReferenceAnatomyAndEscapesText(t *testing.T) {
	row := PeoplePickerRow(PeoplePickerRowProps{ID: "space-ada", Name: "member_ids", Value: "ada", PersonName: "Ada <Okafor>", Detail: "Board & officers", Checked: true, TrailingLabel: "Organizer", Trailing: peoplePickerTestTrailing(), IdentityClass: "flex gap-05 min-w-0", ShowAvatar: true})
	table := peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{Class: "table-extra", Attrs: templ.Attributes{"data-table": "members"}}),
		peoplePickerStaticHTML(`<thead><tr><th scope="col"><span class="visually-hidden">Select</span></th><th scope="col">Member</th><th scope="col">Organizer</th></tr></thead><tbody>`),
		row,
		peoplePickerStaticHTML(`</tbody>`),
	)
	groups := peoplePickerWithChildren(PeoplePickerGroups(PeoplePickerGroupsProps{Class: "cluster gap-025", Attrs: templ.Attributes{"data-groups": "circles"}}), peoplePickerStaticHTML(`<button type="submit" name="group_action" value="board" formnovalidate>Board</button>`))
	actions := peoplePickerWithChildren(PeoplePickerActions(PeoplePickerActionsProps{Class: "cluster cluster-compact", Attrs: templ.Attributes{"data-actions": "picker"}}), peoplePickerStaticHTML(`<button type="submit" name="picker_action" value="clear" formnovalidate>Clear</button>`))
	markup := testutil.RenderToString(t, peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "Members <today>", Description: "Choose & review.", Class: "picker-extra", Attrs: templ.Attributes{"data-picker": "one"}}), groups, table, PeoplePickerStatus(PeoplePickerStatusProps{}), actions))
	root := parsePeoplePicker(t, markup)
	// Unescaped markup would parse as elements, so decoded text is the evidence that escaping happened.
	legend := findPeoplePickerNodes(root, "legend", nil)
	description := findPeoplePickerNodes(root, "p", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-description") })
	visibleName := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-name") })
	if len(legend) != 1 || len(description) != 1 || len(visibleName) != 1 {
		t.Fatalf("legend/description/name counts = %d/%d/%d, want 1 each", len(legend), len(description), len(visibleName))
	}
	for label, want := range map[string]struct {
		node *html.Node
		text string
	}{
		"legend":      {legend[0], "Members <today>"},
		"description": {description[0], "Choose & review."},
		"name+detail": {visibleName[0], "Ada <Okafor>Board & officers"},
	} {
		if got := testutil.NodeText(want.node); got != want.text {
			t.Fatalf("%s text = %q, want %q (escaped in markup, decoded as text)", label, got, want.text)
		}
	}
	for _, raw := range []string{"<today>", "<Okafor>", "Board & officers", "Choose & review."} {
		if strings.Contains(markup, raw) {
			t.Fatalf("unescaped %q in markup: %s", raw, markup)
		}
	}
	if len(findPeoplePickerNodes(root, "okafor", nil)) != 0 || len(findPeoplePickerNodes(root, "today", nil)) != 0 {
		t.Fatal("caller text produced elements")
	}
	if got := len(findPeoplePickerNodes(root, "fieldset", nil)); got != 1 {
		t.Fatalf("fieldset count = %d, want 1", got)
	}
	if got := len(findPeoplePickerNodes(root, "table", nil)); got != 1 {
		t.Fatalf("table count = %d, want 1", got)
	}
	if got := len(findPeoplePickerNodes(root, "tr", nil)); got != 2 {
		t.Fatalf("table row count = %d, want header and one body row", got)
	}
	if got := len(findPeoplePickerNodes(root, "td", nil)); got != 3 {
		t.Fatalf("body cell count = %d, want selection, identity and trailing cells", got)
	}
	groupNodes := findPeoplePickerNodes(root, "div", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-groups") })
	actionNodes := findPeoplePickerNodes(root, "div", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-actions") })
	if len(groupNodes) != 1 || len(actionNodes) != 1 || !strings.Contains(markup, `data-groups="circles"`) || !strings.Contains(markup, `data-actions="picker"`) {
		t.Fatal("optional group or action slots did not render with their wrapper attributes")
	}
	fieldset := findPeoplePickerNodes(root, "fieldset", nil)[0]
	if got, ok := testutil.AttrVal(fieldset, "data-picker"); !ok || got != "one" {
		t.Fatalf("fieldset attribute data-picker = %q, %v", got, ok)
	}
	if !hasPeoplePickerClass(fieldset, "people-picker") || !hasPeoplePickerClass(fieldset, "picker-extra") {
		t.Fatal("fieldset must carry people-picker and the caller Class")
	}
	if _, ok := testutil.AttrVal(findPeoplePickerNodes(root, "legend", nil)[0], "data-picker"); ok {
		t.Fatal("fieldset attributes leaked onto the legend")
	}
	tableNode := findPeoplePickerNodes(root, "table", nil)[0]
	if got, _ := testutil.AttrVal(tableNode, "data-table"); got != "members" || !hasPeoplePickerClass(tableNode, "people-picker-table") || !hasPeoplePickerClass(tableNode, "table--lined") || !hasPeoplePickerClass(tableNode, "table-extra") {
		t.Fatalf("table Class/Attrs not applied: %v", tableNode.Attr)
	}
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 2 {
		t.Fatalf("input count = %d, want selection and trailing controls", len(inputs))
	}
	selectionID, _ := testutil.AttrVal(inputs[0], "id")
	labelledBy, _ := testutil.AttrVal(inputs[0], "aria-labelledby")
	describedBy, ok := testutil.AttrVal(inputs[0], "aria-describedby")
	if selectionID != "space-ada" || labelledBy != "space-ada-selection-label" || describedBy != "space-ada-detail" || !ok {
		t.Fatalf("selection relationships are incorrect: id=%q labelledby=%q describedby=%q", selectionID, labelledBy, describedBy)
	}
	if got, ok := testutil.AttrVal(inputs[0], "checked"); !ok || got != "" {
		t.Fatalf("checked state missing: %q, %v", got, ok)
	}
	if _, ok := testutil.AttrVal(inputs[0], "disabled"); ok {
		t.Fatal("enabled selection input was rendered disabled")
	}
	statusRegions := findPeoplePickerNodes(root, "p", func(n *html.Node) bool { value, ok := testutil.AttrVal(n, "role"); return ok && value == "status" })
	if len(statusRegions) != 1 {
		t.Fatalf("status region count = %d, want one persistent region", len(statusRegions))
	}
	if atomic, ok := testutil.AttrVal(statusRegions[0], "aria-atomic"); !ok || atomic != "true" {
		t.Fatalf("status aria-atomic = %q, %v", atomic, ok)
	}
	if got, ok := testutil.AttrVal(findPeoplePickerNodes(root, "td", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-trailing") })[0], "data-th"); !ok || got != "Organizer" {
		t.Fatalf("trailing visible label = %q, %v", got, ok)
	}
	avatarContainers := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "avatar") })
	if len(avatarContainers) != 1 {
		t.Fatalf("avatar count = %d, want 1", len(avatarContainers))
	}
	for _, name := range []string{"role", "aria-label", "aria-hidden"} {
		if _, ok := testutil.AttrVal(avatarContainers[0], name); ok {
			t.Fatalf("decorative avatar must not carry %s; the identity label names the person", name)
		}
	}
	initials := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "initials") })
	if len(initials) != 1 {
		t.Fatalf("initials count = %d, want 1", len(initials))
	}
	if got, ok := testutil.AttrVal(initials[0], "aria-hidden"); !ok || got != "true" {
		t.Fatalf("initials should be hidden from assistive technology: aria-hidden=%q, %v", got, ok)
	}
}

func TestPeoplePickerOptionalDetailsAndTrailingCell(t *testing.T) {
	row := PeoplePickerRow(PeoplePickerRowProps{ID: "draft-bo", Name: "collaborator_ids", Value: "bo", PersonName: "Bo Lindqvist", Checked: true, Disabled: true})
	table := peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), row, peoplePickerStaticHTML(`</tbody>`))
	markup := testutil.RenderToString(t, table)
	root := parsePeoplePicker(t, markup)
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 1 {
		t.Fatalf("input count = %d, want one", len(inputs))
	}
	if _, ok := testutil.AttrVal(inputs[0], "aria-describedby"); ok {
		t.Fatal("aria-describedby must be omitted when detail text is absent")
	}
	if _, ok := testutil.AttrVal(inputs[0], "checked"); !ok {
		t.Fatal("checked state was omitted")
	}
	if _, ok := testutil.AttrVal(inputs[0], "disabled"); !ok {
		t.Fatal("disabled state was omitted")
	}
	if got := len(findPeoplePickerNodes(root, "td", nil)); got != 2 {
		t.Fatalf("cell count = %d, want selection and identity only", got)
	}
	if got := len(findPeoplePickerNodes(root, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "avatar") })); got != 0 {
		t.Fatalf("avatar count = %d, want omitted avatar", got)
	}
}

func TestPeoplePickerIDsAndInputAttrsStayScoped(t *testing.T) {
	markup := testutil.RenderToString(t, templ.Join(
		peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "First", Attrs: templ.Attributes{"data-picker": "first"}}), peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), PeoplePickerRowComponent(PeoplePickerRowProps{ID: "first-person", Name: "member_ids", Value: "first", PersonName: "Person", Class: "row-owner", Attrs: templ.Attributes{"data-row": "first"}, InputAttrs: templ.Attributes{"data-input": "first"}}), peoplePickerStaticHTML(`</tbody>`))),
		peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "Second", Attrs: templ.Attributes{"data-picker": "second"}}), peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), PeoplePickerRow(PeoplePickerRowProps{ID: "second-person", Name: "member_ids", Value: "second", PersonName: "Person", Class: "row-owner", Attrs: templ.Attributes{"data-row": "second"}, InputAttrs: templ.Attributes{"data-input": "second"}}), peoplePickerStaticHTML(`</tbody>`))),
	))
	root := parsePeoplePicker(t, markup)
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 2 {
		t.Fatalf("input count = %d, want two independent selections", len(inputs))
	}
	ids := map[string]bool{}
	for _, input := range inputs {
		id, _ := testutil.AttrVal(input, "id")
		if id == "" || ids[id] {
			t.Fatalf("duplicate or empty ID %q", id)
		}
		ids[id] = true
		if got, ok := testutil.AttrVal(input, "data-input"); !ok || got == "" {
			t.Fatalf("input attribute missing from input: %q, %v", got, ok)
		}
		if _, ok := testutil.AttrVal(input, "data-row"); ok {
			t.Fatal("row wrapper attribute leaked onto the input")
		}
		labelledBy, _ := testutil.AttrVal(input, "aria-labelledby")
		labels := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { value, ok := testutil.AttrVal(n, "id"); return ok && value == labelledBy })
		if len(labels) != 1 {
			t.Fatalf("accessible label %q resolves to %d elements", labelledBy, len(labels))
		}
	}
	rows := findPeoplePickerNodes(root, "tr", func(n *html.Node) bool { _, ok := testutil.AttrVal(n, "data-row"); return ok })
	if len(rows) != 2 {
		t.Fatalf("row wrapper attributes landed on %d rows", len(rows))
	}
	for _, row := range rows {
		if _, ok := testutil.AttrVal(row, "data-input"); ok {
			t.Fatal("input attribute leaked onto the row wrapper")
		}
		if !hasPeoplePickerClass(row, "row-owner") {
			t.Fatalf("row Class must land on the <tr>: %v", row.Attr)
		}
	}
}

func hasPeoplePickerClass(node *html.Node, class string) bool {
	value, ok := testutil.AttrVal(node, "class")
	if !ok {
		return false
	}
	for _, item := range strings.Fields(value) {
		if item == class {
			return true
		}
	}
	return false
}

func TestPeoplePickerUncheckedEnabledRowOmitsStateAttributes(t *testing.T) {
	markup := testutil.RenderToString(t, PeoplePickerRow(PeoplePickerRowProps{ID: "p-cy", Name: "member_ids", Value: "cy", PersonName: "Cy Moreno"}))
	input := findPeoplePickerNodes(parsePeoplePicker(t, "<table><tbody>"+markup+"</tbody></table>"), "input", nil)[0]
	for _, name := range []string{"checked", "disabled"} {
		if _, ok := testutil.AttrVal(input, name); ok {
			t.Fatalf("%s must be omitted for an unchecked, enabled row: %s", name, markup)
		}
	}
	if got, _ := testutil.AttrVal(input, "value"); got != "cy" {
		t.Fatalf("value = %q", got)
	}
}

func TestPeoplePickerInputAttrsCannotOverrideComponentAttributes(t *testing.T) {
	markup := testutil.RenderToString(t, PeoplePickerRow(PeoplePickerRowProps{
		ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", Detail: "Board",
		InputAttrs: templ.Attributes{"id": "other", "name": "other", "value": "other", "aria-labelledby": "other", "aria-describedby": "other", "data-demo": "kept"},
	}))
	input := findPeoplePickerNodes(parsePeoplePicker(t, "<table><tbody>"+markup+"</tbody></table>"), "input", nil)[0]
	for name, want := range map[string]string{"id": "p-ada", "name": "member_ids", "value": "ada", "aria-labelledby": "p-ada-selection-label", "aria-describedby": "p-ada-detail", "data-demo": "kept"} {
		if got, _ := testutil.AttrVal(input, name); got != want {
			t.Fatalf("%s = %q, want %q (component attributes win on collision, first occurrence)", name, got, want)
		}
	}
}

func TestPeoplePickerTrailingControlRendersWithoutLabel(t *testing.T) {
	markup := testutil.RenderToString(t, PeoplePickerRow(PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", Trailing: peoplePickerTestTrailing()}))
	cells := findPeoplePickerNodes(parsePeoplePicker(t, "<table><tbody>"+markup+"</tbody></table>"), "td", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-trailing") })
	if len(cells) != 1 {
		t.Fatalf("a supplied trailing control must not be dropped: %s", markup)
	}
	if _, ok := testutil.AttrVal(cells[0], "data-th"); ok {
		t.Fatal("data-th must be omitted when no trailing label is supplied")
	}
}

func TestPeoplePickerComposesWithExistingSwitchAndCheckbox(t *testing.T) {
	organizer := SwitchComponent(SwitchProps{Name: "organizer_ids", ID: "organizer-ada", Value: "ada", Variant: "small", InputAttrs: templ.Attributes{"aria-label": "Make Ada an organizer"}})
	notify := Checkbox(CheckboxProps{Name: "notify_ids", Value: "bo", Label: "Notify Bo"})
	markup := testutil.RenderToString(t, peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "Members"}), peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`),
		PeoplePickerRow(PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", Checked: true, TrailingLabel: "Organizer", Trailing: organizer}),
		PeoplePickerRow(PeoplePickerRowProps{ID: "p-bo", Name: "member_ids", Value: "bo", PersonName: "Bo", TrailingLabel: "Notify", Trailing: notify}),
		peoplePickerStaticHTML(`</tbody>`))))
	root := parsePeoplePicker(t, markup)
	seen := map[string]bool{}
	for _, node := range findPeoplePickerNodes(root, "span", func(n *html.Node) bool { _, ok := testutil.AttrVal(n, "id"); return ok }) {
		id, _ := testutil.AttrVal(node, "id")
		if seen[id] {
			t.Fatalf("duplicate id %q across row and existing component", id)
		}
		seen[id] = true
	}
	for _, input := range findPeoplePickerNodes(root, "input", nil) {
		id, _ := testutil.AttrVal(input, "id")
		if id != "" && seen[id] {
			t.Fatalf("duplicate id %q across row and existing component", id)
		}
		seen[id] = true
	}
	trailing := findPeoplePickerNodes(root, "td", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-trailing") })
	if len(trailing) != 2 {
		t.Fatalf("trailing cells = %d, want 2", len(trailing))
	}
	for _, cell := range trailing {
		if len(findPeoplePickerNodes(cell, "input", nil)) != 1 {
			t.Fatal("each trailing cell must hold exactly its own control")
		}
	}
	// Nothing interactive may live inside the selection or identity labels besides the selection input.
	for _, label := range findPeoplePickerNodes(root, "label", func(n *html.Node) bool {
		return hasPeoplePickerClass(n, "checkbox--icon") || hasPeoplePickerClass(n, "people-picker-person")
	}) {
		for _, tag := range []string{"button", "a", "select", "textarea"} {
			if len(findPeoplePickerNodes(label, tag, nil)) != 0 {
				t.Fatalf("<%s> nested in a selection label", tag)
			}
		}
		if hasPeoplePickerClass(label, "people-picker-person") && len(findPeoplePickerNodes(label, "input", nil)) != 0 {
			t.Fatal("identity label must not contain inputs")
		}
	}
}

func peoplePickerRow(t *testing.T, props PeoplePickerRowProps) (markup string, root *html.Node) {
	t.Helper()
	markup = testutil.RenderToString(t, PeoplePickerRow(props))
	return markup, parsePeoplePicker(t, "<table><tbody>"+markup+"</tbody></table>")
}

func TestPeoplePickerStatusRendersCallerTextAndAttributes(t *testing.T) {
	markup := testutil.RenderToString(t, PeoplePickerStatus(PeoplePickerStatusProps{Text: "Board <selected> & 2 people", Class: "text-sm text-muted", Attrs: templ.Attributes{"data-status": "picker"}}))
	statuses := findPeoplePickerNodes(parsePeoplePicker(t, markup), "p", func(n *html.Node) bool { role, _ := testutil.AttrVal(n, "role"); return role == "status" })
	if len(statuses) != 1 {
		t.Fatalf("status region count = %d: %s", len(statuses), markup)
	}
	status := statuses[0]
	if got := testutil.NodeText(status); got != "Board <selected> & 2 people" {
		t.Fatalf("status text = %q", got)
	}
	if strings.Contains(markup, "<selected>") || !strings.Contains(markup, "Board &lt;selected&gt; &amp; 2 people") {
		t.Fatalf("status text was not escaped: %s", markup)
	}
	if atomic, _ := testutil.AttrVal(status, "aria-atomic"); atomic != "true" {
		t.Fatalf("aria-atomic = %q", atomic)
	}
	if !hasPeoplePickerClass(status, "people-picker-status") || !hasPeoplePickerClass(status, "text-sm") {
		t.Fatalf("status classes = %v", status.Attr)
	}
	if got, _ := testutil.AttrVal(status, "data-status"); got != "picker" {
		t.Fatalf("status Attrs not applied: %v", status.Attr)
	}
}

func TestPeoplePickerTruncationKeepsFullTextAndAddsTitles(t *testing.T) {
	const name = "Maximiliana Featherstonehaugh-Abernathy"
	const detail = "Organizer, Festival submissions working group"
	_, root := peoplePickerRow(t, PeoplePickerRowProps{ID: "narrow-long", Name: "member_ids", Value: "long", PersonName: name, Detail: detail, NameClass: "block truncate", DetailClass: "truncate"})
	nameSpans := findPeoplePickerNodes(root, "span", func(n *html.Node) bool {
		return hasPeoplePickerClass(n, "block") && hasPeoplePickerClass(n, "truncate")
	})
	detailSpans := findPeoplePickerNodes(root, "span", func(n *html.Node) bool {
		return hasPeoplePickerClass(n, "people-picker-detail") && hasPeoplePickerClass(n, "truncate")
	})
	if len(nameSpans) != 1 || len(detailSpans) != 1 {
		t.Fatalf("truncating spans = %d name, %d detail; want 1 each", len(nameSpans), len(detailSpans))
	}
	for label, want := range map[string]struct {
		node *html.Node
		text string
	}{"name": {nameSpans[0], name}, "detail": {detailSpans[0], detail}} {
		if got := testutil.NodeText(want.node); got != want.text {
			t.Fatalf("%s text = %q, want the full text %q", label, got, want.text)
		}
		if got, ok := testutil.AttrVal(want.node, "title"); !ok || got != want.text {
			t.Fatalf("%s title = %q, %v; want the full text for a mouse hint", label, got, ok)
		}
	}
	if got, _ := testutil.AttrVal(detailSpans[0], "id"); got != "narrow-long-detail" {
		t.Fatalf("detail id = %q", got)
	}
	hidden := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { id, _ := testutil.AttrVal(n, "id"); return id == "narrow-long-selection-label" })
	if len(hidden) != 1 || testutil.NodeText(hidden[0]) != "Select "+name {
		t.Fatalf("the accessible label must keep the full name when truncated")
	}

	_, plain := peoplePickerRow(t, PeoplePickerRowProps{ID: "wrap", Name: "member_ids", Value: "wrap", PersonName: name, Detail: detail})
	for _, span := range findPeoplePickerNodes(plain, "span", func(n *html.Node) bool { _, ok := testutil.AttrVal(n, "title"); return ok }) {
		t.Fatalf("title must only appear when truncation classes are requested: %v", span.Attr)
	}
}

func TestPeoplePickerTrailingLabelWithoutControlRendersEmptyLabelledCell(t *testing.T) {
	_, root := peoplePickerRow(t, PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", TrailingLabel: "Organizer"})
	cells := findPeoplePickerNodes(root, "td", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-trailing") })
	if len(cells) != 1 {
		t.Fatalf("trailing cells = %d, want 1 so the column stays aligned", len(cells))
	}
	if got, _ := testutil.AttrVal(cells[0], "data-th"); got != "Organizer" {
		t.Fatalf("data-th = %q", got)
	}
	if cells[0].FirstChild != nil {
		t.Fatal("trailing cell must be empty when the caller supplies no control")
	}
}

func TestPeoplePickerAvatarRequiresShowAvatar(t *testing.T) {
	_, hidden := peoplePickerRow(t, PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", AvatarSrc: "/img/ada.jpg"})
	if len(findPeoplePickerNodes(hidden, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "avatar") })) != 0 || len(findPeoplePickerNodes(hidden, "img", nil)) != 0 {
		t.Fatal("AvatarSrc without ShowAvatar must render no avatar")
	}
	_, shown := peoplePickerRow(t, PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", AvatarSrc: "/img/ada.jpg", ShowAvatar: true})
	images := findPeoplePickerNodes(shown, "img", nil)
	if len(images) != 1 {
		t.Fatalf("avatar images = %d, want 1", len(images))
	}
	if alt, ok := testutil.AttrVal(images[0], "alt"); !ok || alt != "" {
		t.Fatalf("decorative avatar image alt = %q, %v; want present and empty", alt, ok)
	}
	if got, _ := testutil.AttrVal(images[0], "src"); got != "/img/ada.jpg" {
		t.Fatalf("avatar src = %q", got)
	}
}

// A row derives "<id>-selection-label" and "<id>-detail". Field derives "<name>-hint" and "<name>-error".
// Sharing one base between a row and a Field is the caller's uniqueness duty, so the base id is the only
// id the document may repeat: no derived suffix may collide.
func TestPeoplePickerDerivedIDsDoNotCollideWithFieldIDs(t *testing.T) {
	markup := testutil.RenderToString(t, templ.Join(
		peoplePickerStaticHTML(`<table><tbody>`),
		PeoplePickerRow(PeoplePickerRowProps{ID: "p-ada", Name: "member_ids", Value: "ada", PersonName: "Ada", Detail: "Board"}),
		peoplePickerStaticHTML(`</tbody></table>`),
		Field(FieldProps{Label: "Role", Name: "p-ada", Hint: "Optional", Error: "Too long"}),
	))
	counts := map[string]int{}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if id, ok := testutil.AttrVal(n, "id"); ok {
				counts[id]++
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(parsePeoplePicker(t, markup))
	for _, id := range []string{"p-ada-selection-label", "p-ada-detail", "p-ada-hint", "p-ada-error"} {
		if counts[id] != 1 {
			t.Fatalf("derived id %q appears %d times, want exactly 1", id, counts[id])
		}
	}
	if counts["p-ada"] != 2 {
		t.Fatalf("the shared base id should be the only repeated id, got %d uses: %v", counts["p-ada"], counts)
	}
}
