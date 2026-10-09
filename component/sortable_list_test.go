package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

// sortable-list.html and sortable-table.html are verbatim host excerpts from
// Design-Machines-Studio/livewires@21cff3f08d156feabd9da4bbd3d366781c463f96,
// public/reference/components/sortable-list.html (MIT):
// https://github.com/Design-Machines-Studio/livewires/blob/21cff3f08d156feabd9da4bbd3d366781c463f96/public/reference/components/sortable-list.html
// The captured full source SHA-256 is
// 2c98ac1aab9e3876d6ea6433ddb9db5eebc8ec692c6cde6aff0115457d033778.
// Only the two reference tests claim verbatim parity. Variants below are adaptations.

func sortableReading(variant, density string, handles, moves bool) templ.Component {
	items := []struct {
		id, label, content, up, down string
	}{
		{"essay-01", "Field notes", `<span data-sortable-content><a href="#field-notes">Field notes</a> <small>— a link stays usable</small></span>`, "essay-01", "essay-03"},
		{"essay-02", "On paper", `<span data-sortable-content><a href="#on-paper">On paper</a> <small>— nested controls remain available</small></span>`, "essay-01", ""},
		{"essay-03", "An index", `<span data-sortable-content><a href="#an-index">An index</a><label>Note <input type="text" aria-label="Note about An index" value="Keep this input usable"></label></span>`, "essay-02", ""},
	}
	children := []templ.Component{staticHTML(`<ol data-sortable-list>`)}
	for i, item := range items {
		actions := []templ.Component{staticHTML(`<span data-sortable-actions>`)}
		if handles {
			actions = append(actions, SortableHandle(item.label))
		}
		if moves {
			actions = append(actions,
				SortableMoveComponent(SortableMoveProps{Label: item.label, Direction: "up", Action: "/items/reorder", ItemID: item.id, Before: item.up, Disabled: i == 0}),
				SortableMoveComponent(SortableMoveProps{Label: item.label, Direction: "down", Action: "/items/reorder", ItemID: item.id, Before: item.down, Disabled: i == len(items)-1}),
			)
		}
		actions = append(actions, staticHTML(`</span>`))
		children = append(children, withChildren(SortableItemComponent(SortableItemProps{ItemID: item.id, Label: item.label}),
			staticHTML(item.content), templ.Join(actions...),
		))
	}
	children = append(children, staticHTML(`</ol>`))
	return withChildren(SortableListComponent(SortableListProps{Label: "Reading list", Variant: variant, Density: density}), children...)
}

func sortableTable() templ.Component {
	children := []templ.Component{staticHTML(`<table class="table--lined"><caption>Member requirements</caption><thead><tr><th scope="col"><span class="visually-hidden">Order</span></th><th scope="col">Member practice</th><th scope="col">Shared purpose</th><th scope="col"><span class="visually-hidden">Position controls</span></th></tr></thead><tbody data-sortable-list>`)}
	for _, row := range []struct{ id, label, purpose string }{
		{"orientation", "Member orientation", "Welcome people into the co-op and how members work together."},
		{"agreements", "Shared agreements", "Review the agreements members use to work together."},
		{"governance", "Governance participation", "Take part in a member meeting or decision round."},
	} {
		children = append(children, withChildren(SortableItemComponent(SortableItemProps{ItemID: row.id, Label: row.label, AsRow: true}),
			staticHTML(`<td><span data-sortable-actions>`), SortableHandle(row.label), staticHTML(`</span></td>`),
			staticHTML(`<th scope="row" data-sortable-content>`+row.label+`</th><td data-th="Shared purpose">`+row.purpose+`</td><td data-th="Position"><span data-sortable-actions>`),
			SortableMove(row.label, "up"), SortableMove(row.label, "down"), staticHTML(`</span></td>`),
		))
	}
	children = append(children, staticHTML(`</tbody></table>`))
	return withChildren(SortableList("Member requirements"), children...)
}

func TestReferenceSortableReading(t *testing.T) {
	assertMatchesReference(t, "sortable-list.html", sortableReading("cards", "", true, true))
}

func TestReferenceSortableTable(t *testing.T) {
	assertMatchesReference(t, "sortable-table.html", sortableTable())
}

// sortableElements inspects the actual parsed tree, including table repair.
func sortableElements(nodes []*html.Node, tag, attr string) []*html.Node {
	var found []*html.Node
	var visit func(*html.Node)
	visit = func(node *html.Node) {
		_, hasAttr := testutil.AttrVal(node, attr)
		if node.Type == html.ElementNode && (tag == "" || node.Data == tag) && (attr == "" || hasAttr) {
			found = append(found, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	for _, node := range nodes {
		visit(node)
	}
	return found
}

func sortableAttr(t *testing.T, node *html.Node, name, want string) {
	t.Helper()
	if node == nil {
		t.Fatalf("missing element for %s", name)
	}
	if got, ok := testutil.AttrVal(node, name); !ok || got != want {
		t.Errorf("%s = %q (present %v), want %q", name, got, ok, want)
	}
}

func TestSortableAdaptedVariants(t *testing.T) {
	// The producer describes these adaptations; it does not supply verbatim hosts.
	for _, variant := range []string{"", "cards"} {
		for _, density := range []string{"", "compact"} {
			for _, controls := range []struct {
				name           string
				handles, moves bool
			}{{"handle-only", true, false}, {"arrow-only", false, true}, {"combined", true, true}} {
				t.Run(variant+"/"+density+"/"+controls.name, func(t *testing.T) {
					output := testutil.RenderToString(t, sortableReading(variant, density, controls.handles, controls.moves))
					nodes := testutil.ParseFragment(t, output)
					host := testutil.FindElement(nodes, "lw-sortable-list")
					for name, want := range map[string]string{"data-sortable-variant": variant, "data-sortable-density": density} {
						got, present := testutil.AttrVal(host, name)
						if got != want || present != (want != "") {
							t.Errorf("%s = %q (present %v), want %q", name, got, present, want)
						}
					}
					wantHandles, wantMoves := 0, 0
					if controls.handles {
						wantHandles = 3
					}
					if controls.moves {
						wantMoves = 6
					}
					if len(sortableElements(nodes, "button", "data-sortable-handle")) != wantHandles || len(sortableElements(nodes, "button", "data-sortable-move")) != wantMoves {
						t.Error("control selection did not preserve composition")
					}
					if len(sortableElements(nodes, "li", "data-sortable-item")) != 3 || len(sortableElements(nodes, "a", "")) != 3 || len(sortableElements(nodes, "input", "aria-label")) != 1 {
						t.Error("list items or ordinary interactive content lost")
					}
					for _, node := range sortableElements(nodes, "", "") {
						if _, present := testutil.AttrVal(node, "class"); present {
							t.Error("component introduced a class without caller classes")
						}
						if node.Data == "script" || node.Data == "style" {
							t.Error("component introduced an asset")
						}
					}
				})
			}
		}
	}
}

func TestSortableNativeReferenceValues(t *testing.T) {
	nodes := testutil.ParseFragment(t, testutil.RenderToString(t, sortableReading("cards", "", true, true)))
	forms := sortableElements(nodes, "form", "")
	if len(forms) != 6 {
		t.Fatalf("got %d forms, want 6", len(forms))
	}
	for i, want := range []struct {
		itemID, before, direction string
		disabled                  bool
	}{
		{"essay-01", "essay-01", "up", true},
		{"essay-01", "essay-03", "down", false},
		{"essay-02", "essay-01", "up", false},
		{"essay-02", "", "down", false},
		{"essay-03", "essay-02", "up", false},
		{"essay-03", "", "down", true},
	} {
		form := forms[i]
		sortableAttr(t, form, "method", "post")
		sortableAttr(t, form, "action", "/items/reorder")
		inputs := sortableElements([]*html.Node{form}, "input", "")
		if len(inputs) != 2 {
			t.Fatalf("form %d has %d inputs, want 2", i, len(inputs))
		}
		sortableAttr(t, inputs[0], "name", "itemId")
		sortableAttr(t, inputs[0], "value", want.itemID)
		sortableAttr(t, inputs[1], "name", "before")
		sortableAttr(t, inputs[1], "value", want.before)
		for _, input := range inputs {
			sortableAttr(t, input, "type", "hidden")
		}
		button := testutil.FindElement([]*html.Node{form}, "button")
		sortableAttr(t, button, "type", "submit")
		sortableAttr(t, button, "data-sortable-move", want.direction)
		if _, disabled := testutil.AttrVal(button, "disabled"); disabled != want.disabled {
			t.Errorf("form %d disabled = %v, want %v", i, disabled, want.disabled)
		}
	}
}

func TestSortableMetadataAndAttributeOwnership(t *testing.T) {
	props := SortableMoveProps{
		Label: "Field notes", Direction: "up", Action: "/reorder?scope=reading&revision=2", Method: "get", ItemID: "essay-01", Before: "essay-01", Disabled: true,
		Metadata: staticHTML(`<input type="hidden" name="csrf" value="development-token"><input type="hidden" name="revision" value="2">`),
		Class:    "button button--small", Attrs: templ.Attributes{"data-button": "move", "type": "reset", "data-sortable-move": "down", "aria-label": "wrong", "disabled": "false"},
		FormAttrs: templ.Attributes{"data-form": "native", "method": "post", "action": "/wrong"},
	}
	output := testutil.RenderToString(t, SortableMoveComponent(props))
	nodes := testutil.ParseFragment(t, output)
	form := testutil.FindElement(nodes, "form")
	button := testutil.FindElement(nodes, "button")
	sortableAttr(t, form, "method", "get")
	sortableAttr(t, form, "action", props.Action)
	sortableAttr(t, form, "data-form", "native")
	sortableAttr(t, button, "type", "submit")
	sortableAttr(t, button, "data-sortable-move", "up")
	sortableAttr(t, button, "aria-label", "Move Field notes up")
	sortableAttr(t, button, "disabled", "")
	sortableAttr(t, button, "class", props.Class)
	sortableAttr(t, button, "data-button", "move")
	if _, present := testutil.AttrVal(form, "class"); present {
		t.Error("button Class leaked onto form")
	}
	if _, present := testutil.AttrVal(button, "data-form"); present {
		t.Error("form attributes leaked onto button")
	}
	for tag, attrs := range map[string][]string{"form": {"method", "action"}, "button": {"type", "data-sortable-move", "aria-label", "disabled"}} {
		for _, attr := range attrs {
			// Like the existing wrapped controls, trusted spreads retain duplicates;
			// the component owns the first occurrence in parsed HTML.
			if testutil.CountAttr(testutil.RawTag(t, output, tag), attr) != 2 {
				t.Errorf("expected owned %s before the caller's duplicate on %s", attr, tag)
			}
		}
	}
	inputs := sortableElements(nodes, "input", "")
	if len(inputs) != 4 {
		t.Fatalf("got %d hidden fields, want 4", len(inputs))
	}
	sortableAttr(t, inputs[2], "name", "csrf")
	sortableAttr(t, inputs[2], "value", "development-token")
	sortableAttr(t, inputs[3], "name", "revision")
	sortableAttr(t, inputs[3], "value", "2")
	props.Action = ""
	client := testutil.RenderToString(t, SortableMoveComponent(props))
	clientNodes := testutil.ParseFragment(t, client)
	if len(sortableElements(clientNodes, "form", "")) != 0 || len(sortableElements(clientNodes, "input", "")) != 0 || strings.Contains(client, "data-form") {
		t.Error("client mode rendered native form seams")
	}
	clientButton := testutil.FindElement(clientNodes, "button")
	sortableAttr(t, clientButton, "type", "button")
	sortableAttr(t, clientButton, "disabled", "")
	sortableAttr(t, clientButton, "class", props.Class)
	sortableAttr(t, clientButton, "data-button", "move")
}

func TestSortableInstancesAndStatus(t *testing.T) {
	first := withChildren(SortableListComponent(SortableListProps{ID: "reading", Label: "Reading list", Class: "caller", Attrs: templ.Attributes{"data-on:sortable-move-request": "@post('/reorder')"}}),
		staticHTML(`<ol data-sortable-list>`), withChildren(SortableItem("essay-01", "Field notes"), SortableHandle("Field notes")), staticHTML(`</ol>`))
	second := sortableTable()
	nodes := testutil.ParseFragment(t, testutil.RenderToString(t, templ.Join(first, second)))
	hosts := sortableElements(nodes, "lw-sortable-list", "")
	if len(hosts) != 2 {
		t.Fatalf("got %d hosts, want 2", len(hosts))
	}
	sortableAttr(t, hosts[0], "id", "reading")
	sortableAttr(t, hosts[0], "class", "caller")
	sortableAttr(t, hosts[0], "data-on:sortable-move-request", "@post('/reorder')")
	if _, present := testutil.AttrVal(hosts[1], "id"); present {
		t.Error("optional ID rendered without a caller ID")
	}
	for _, host := range hosts {
		statuses := sortableElements([]*html.Node{host}, "", "data-sortable-status")
		if len(statuses) != 1 {
			t.Fatalf("host has %d status nodes, want 1", len(statuses))
		}
		status := statuses[0]
		if status.Data != "p" || status.Parent != host || status.NextSibling != nil || len(status.Attr) != 5 || status.FirstChild != nil {
			t.Error("status must be the exact empty paragraph after children, directly inside host")
		}
		for name, want := range map[string]string{"data-sortable-status": "", "role": "status", "aria-live": "polite", "aria-atomic": "true", "tabindex": "-1"} {
			sortableAttr(t, status, name, want)
		}
	}
	for _, row := range sortableElements(nodes, "tr", "data-sortable-item") {
		if row.Parent.Data != "tbody" || row.Parent.Parent.Data != "table" {
			t.Error("sortable row lost its semantic table ancestry")
		}
		if row.FirstChild.Data != "td" || row.LastChild.Data != "td" {
			t.Error("table handle and arrows must remain in left and right cells")
		}
	}
}

func TestSortableEscaping(t *testing.T) {
	label, id := `Notes <script> & "quotes"`, `item"<&`
	move := SortableMoveComponent(SortableMoveProps{Label: label, Direction: "down", Action: `/reorder?x="<&`, ItemID: id, Before: `target"<&`, Class: `caller"<&`, Attrs: templ.Attributes{"data-note": label}})
	item := withChildren(SortableItemComponent(SortableItemProps{ItemID: id, Label: label, Class: "item-caller", Attrs: templ.Attributes{"data-item-note": label}}),
		SortableHandleComponent(SortableHandleProps{Label: label, Class: "handle-caller", Attrs: templ.Attributes{"data-handle-note": label}}), move)
	output := testutil.RenderToString(t, withChildren(SortableListComponent(SortableListProps{ID: id, Label: label, Class: "host-caller", Attrs: templ.Attributes{"data-host-note": label}}), staticHTML(`<ol data-sortable-list>`), item, staticHTML(`</ol>`)))
	nodes := testutil.ParseFragment(t, output)
	if len(sortableElements(nodes, "script", "")) != 0 {
		t.Error("unescaped label introduced a script")
	}
	host := testutil.FindElement(nodes, "lw-sortable-list")
	sortableAttr(t, host, "id", id)
	sortableAttr(t, host, "aria-label", label)
	sortableAttr(t, host, "data-host-note", label)
	sortableAttr(t, host, "class", "host-caller")
	li := testutil.FindElement(nodes, "li")
	sortableAttr(t, li, "data-item-id", id)
	sortableAttr(t, li, "data-item-label", label)
	sortableAttr(t, li, "class", "item-caller")
	sortableAttr(t, li, "data-item-note", label)
	buttons := sortableElements(nodes, "button", "")
	if len(buttons) != 2 {
		t.Fatalf("got %d buttons, want 2", len(buttons))
	}
	sortableAttr(t, buttons[0], "aria-label", "Move "+label)
	sortableAttr(t, buttons[0], "class", "handle-caller")
	sortableAttr(t, buttons[0], "data-handle-note", label)
	sortableAttr(t, buttons[1], "aria-label", "Move "+label+" down")
	sortableAttr(t, buttons[1], "data-note", label)
	sortableAttr(t, buttons[1], "class", `caller"<&`)
	sortableAttr(t, testutil.FindElement(nodes, "form"), "action", `/reorder?x="<&`)
	inputs := sortableElements(nodes, "input", "")
	sortableAttr(t, inputs[0], "value", id)
	sortableAttr(t, inputs[1], "value", `target"<&`)
	unsafe := testutil.RenderToString(t, SortableMoveComponent(SortableMoveProps{Label: label, Direction: "up", Action: "javascript:alert(1)"}))
	sortableAttr(t, testutil.FindElement(testutil.ParseFragment(t, unsafe), "form"), "action", "about:invalid#TemplFailedSanitizationURL")
}
