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
	return templ.ComponentFunc(func(_ context.Context, writer io.Writer) error {
		_, err := io.WriteString(writer, `<label class="switch switch--small"><input type="checkbox" role="switch" aria-label="Make Ada an organizer"><span></span></label>`)
		return err
	})
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

func peoplePickerAttr(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func TestPeoplePickerMatchesReferenceAnatomyAndEscapesText(t *testing.T) {
	row := PeoplePickerRow(PeoplePickerRowProps{ID: "space-ada", InputName: "member_ids", Value: "ada", Name: "Ada <Okafor>", Detail: "Board & officers", Checked: true, TrailingLabel: "Organizer", Trailing: peoplePickerTestTrailing(), IdentityClass: "flex gap-05 min-w-0", ShowAvatar: true})
	table := peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}),
		peoplePickerStaticHTML(`<thead><tr><th scope="col"><span class="visually-hidden">Select</span></th><th scope="col">Member</th><th scope="col">Organizer</th></tr></thead><tbody>`),
		row,
		peoplePickerStaticHTML(`</tbody>`),
	)
	groups := peoplePickerWithChildren(PeoplePickerGroups(PeoplePickerGroupsProps{Class: "cluster gap-025", Attrs: templ.Attributes{"data-groups": "circles"}}), peoplePickerStaticHTML(`<button type="submit" name="group_action" value="board" formnovalidate>Board</button>`))
	actions := peoplePickerWithChildren(PeoplePickerActions(PeoplePickerActionsProps{Class: "cluster cluster-compact", Attrs: templ.Attributes{"data-actions": "picker"}}), peoplePickerStaticHTML(`<button type="submit" name="picker_action" value="clear" formnovalidate>Clear</button>`))
	markup := testutil.RenderToString(t, peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "Members <today>", Description: "Choose & review.", Attrs: templ.Attributes{"data-picker": "one"}}), groups, table, PeoplePickerStatus(PeoplePickerStatusProps{}), actions))
	if !strings.Contains(markup, "Members &lt;today&gt;") || !strings.Contains(markup, "Ada &lt;Okafor&gt;") || !strings.Contains(markup, "Board &amp; officers") {
		t.Fatalf("text content was not escaped: %s", markup)
	}
	root := parsePeoplePicker(t, markup)
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
	if !strings.Contains(markup, `data-picker="one"`) {
		t.Fatal("fieldset attributes did not stay on the fieldset")
	}
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 2 {
		t.Fatalf("input count = %d, want selection and trailing controls", len(inputs))
	}
	selectionID, _ := peoplePickerAttr(inputs[0], "id")
	labelledBy, _ := peoplePickerAttr(inputs[0], "aria-labelledby")
	describedBy, ok := peoplePickerAttr(inputs[0], "aria-describedby")
	if selectionID != "space-ada" || labelledBy != "space-ada-selection-label" || describedBy != "space-ada-detail" || !ok {
		t.Fatalf("selection relationships are incorrect: id=%q labelledby=%q describedby=%q", selectionID, labelledBy, describedBy)
	}
	if got, ok := peoplePickerAttr(inputs[0], "checked"); !ok || got != "" {
		t.Fatalf("checked state missing: %q, %v", got, ok)
	}
	if _, ok := peoplePickerAttr(inputs[0], "disabled"); ok {
		t.Fatal("enabled selection input was rendered disabled")
	}
	if got, ok := peoplePickerAttr(inputs[1], "aria-hidden"); ok || got != "" {
		t.Fatal("unexpected attribute on trailing input")
	}
	statusRegions := findPeoplePickerNodes(root, "p", func(n *html.Node) bool { value, ok := peoplePickerAttr(n, "role"); return ok && value == "status" })
	if len(statusRegions) != 1 {
		t.Fatalf("status region count = %d, want one persistent region", len(statusRegions))
	}
	if atomic, ok := peoplePickerAttr(statusRegions[0], "aria-atomic"); !ok || atomic != "true" {
		t.Fatalf("status aria-atomic = %q, %v", atomic, ok)
	}
	if got, ok := peoplePickerAttr(findPeoplePickerNodes(root, "td", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker-trailing") })[0], "data-th"); !ok || got != "Organizer" {
		t.Fatalf("trailing visible label = %q, %v", got, ok)
	}
	if !strings.Contains(markup, "<th scope=\"col\">Organizer</th>") || !strings.Contains(markup, "data-th=\"Organizer\"") {
		t.Fatal("trailing column heading and narrow row label must match")
	}
	avatarContainers := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { return hasPeoplePickerClass(n, "avatar") })
	if len(avatarContainers) != 1 {
		t.Fatalf("avatar count = %d, want 1", len(avatarContainers))
	}
	if got, ok := peoplePickerAttr(avatarContainers[0], "aria-hidden"); !ok || got != "true" {
		t.Fatalf("avatar should be decorative: aria-hidden=%q, %v", got, ok)
	}
}

func TestPeoplePickerOptionalDetailsAndTrailingCell(t *testing.T) {
	row := PeoplePickerRow(PeoplePickerRowProps{ID: "draft-bo", InputName: "collaborator_ids", Value: "bo", Name: "Bo Lindqvist", Checked: true, Disabled: true})
	table := peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), row, peoplePickerStaticHTML(`</tbody>`))
	markup := testutil.RenderToString(t, table)
	root := parsePeoplePicker(t, markup)
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 1 {
		t.Fatalf("input count = %d, want one", len(inputs))
	}
	if _, ok := peoplePickerAttr(inputs[0], "aria-describedby"); ok {
		t.Fatal("aria-describedby must be omitted when detail text is absent")
	}
	if _, ok := peoplePickerAttr(inputs[0], "checked"); !ok {
		t.Fatal("checked state was omitted")
	}
	if _, ok := peoplePickerAttr(inputs[0], "disabled"); !ok {
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
		peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "First", Attrs: templ.Attributes{"data-picker": "first"}}), peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), PeoplePickerRowComponent(PeoplePickerRowProps{ID: "first-person", InputName: "member_ids", Value: "first", Name: "Person", Class: "row-owner", Attrs: templ.Attributes{"data-row": "first"}, InputAttrs: templ.Attributes{"data-input": "first"}}), peoplePickerStaticHTML(`</tbody>`))),
		peoplePickerWithChildren(PeoplePicker(PeoplePickerProps{Legend: "Second", Attrs: templ.Attributes{"data-picker": "second"}}), peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), peoplePickerStaticHTML(`<tbody>`), PeoplePickerRow(PeoplePickerRowProps{ID: "second-person", InputName: "member_ids", Value: "second", Name: "Person", Class: "row-owner", Attrs: templ.Attributes{"data-row": "second"}, InputAttrs: templ.Attributes{"data-input": "second"}}), peoplePickerStaticHTML(`</tbody>`))),
	))
	root := parsePeoplePicker(t, markup)
	inputs := findPeoplePickerNodes(root, "input", nil)
	if len(inputs) != 2 {
		t.Fatalf("input count = %d, want two independent selections", len(inputs))
	}
	ids := map[string]bool{}
	for _, input := range inputs {
		id, _ := peoplePickerAttr(input, "id")
		if id == "" || ids[id] {
			t.Fatalf("duplicate or empty ID %q", id)
		}
		ids[id] = true
		if got, ok := peoplePickerAttr(input, "data-input"); !ok || got == "" {
			t.Fatalf("input attribute missing from input: %q, %v", got, ok)
		}
		if _, ok := peoplePickerAttr(input, "data-row"); ok {
			t.Fatal("row wrapper attribute leaked onto the input")
		}
		labelledBy, _ := peoplePickerAttr(input, "aria-labelledby")
		labels := findPeoplePickerNodes(root, "span", func(n *html.Node) bool { value, ok := peoplePickerAttr(n, "id"); return ok && value == labelledBy })
		if len(labels) != 1 {
			t.Fatalf("accessible label %q resolves to %d elements", labelledBy, len(labels))
		}
	}
	rows := findPeoplePickerNodes(root, "tr", func(n *html.Node) bool { _, ok := peoplePickerAttr(n, "data-row"); return ok })
	if len(rows) != 2 {
		t.Fatalf("row wrapper attributes landed on %d rows", len(rows))
	}
	for _, row := range rows {
		if _, ok := peoplePickerAttr(row, "data-input"); ok {
			t.Fatal("input attribute leaked onto the row wrapper")
		}
	}
}

func hasPeoplePickerClass(node *html.Node, class string) bool {
	value, ok := peoplePickerAttr(node, "class")
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
