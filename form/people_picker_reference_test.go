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

// The fixture is copied verbatim from Live Wires v2.25.0 (aad7e3843898745a2f29833393b20896a204b90e),
// public/reference/components/people-picker.html. Templ ids the hidden selection label
// "<id>-selection-label" so it cannot collide with the "<id>-label" span that form.Checkbox and
// form.Radio emit (labelSpanID in form/ids.go); the producer uses "<id>-label". Parity normalizes only
// that suffix.
const peoplePickerReferenceFixture = "testdata/livewires/people-picker-reference.html"

func peoplePickerReferenceFieldsets(t *testing.T) []string {
	t.Helper()
	source, err := os.ReadFile(peoplePickerReferenceFixture)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	var fieldsets []string
	for _, node := range findPeoplePickerNodes(doc, "fieldset", func(n *html.Node) bool { return hasPeoplePickerClass(n, "people-picker") }) {
		var rendered bytes.Buffer
		if err := html.Render(&rendered, node); err != nil {
			t.Fatal(err)
		}
		fieldsets = append(fieldsets, rendered.String())
	}
	return fieldsets
}

func assertPeoplePickerMatchesReference(t *testing.T, index int, component templ.Component) {
	t.Helper()
	fieldsets := peoplePickerReferenceFieldsets(t)
	if index >= len(fieldsets) {
		t.Fatalf("reference has %d pickers, want index %d", len(fieldsets), index)
	}
	got := strings.ReplaceAll(testutil.CanonicalHTML(t, testutil.RenderToString(t, component)), "-selection-label", "-label")
	want := testutil.CanonicalHTML(t, fieldsets[index])
	if got != want {
		t.Fatalf("picker %d differs from the published reference\ngot:\n%s\nwant:\n%s", index, got, want)
	}
}

func peoplePickerReferenceTable(header string, rows ...templ.Component) templ.Component {
	children := []templ.Component{peoplePickerStaticHTML(`<thead><tr><th scope="col"><span class="visually-hidden">Select</span></th><th scope="col">Member</th>` + header + `</tr></thead><tbody>`)}
	children = append(children, rows...)
	children = append(children, peoplePickerStaticHTML(`</tbody>`))
	return peoplePickerWithChildren(PeoplePickerTable(PeoplePickerTableProps{}), children...)
}

func TestPeoplePickerReferenceCollaboratorsWithoutGroupsOrTrailingColumn(t *testing.T) {
	assertPeoplePickerMatchesReference(t, 2, peoplePickerWithChildren(
		PeoplePicker(PeoplePickerProps{Legend: "Draft collaborators", Description: "Collaborators can read, discuss and edit this draft."}),
		peoplePickerReferenceTable("",
			PeoplePickerRow(PeoplePickerRowProps{ID: "draft-ada", Name: "collaborator_ids", Value: "ada", PersonName: "Ada Okafor", Checked: true, IdentityClass: "flex gap-05 min-w-0"}),
			PeoplePickerRow(PeoplePickerRowProps{ID: "draft-bo", Name: "collaborator_ids", Value: "bo", PersonName: "Bo Lindqvist", IdentityClass: "flex gap-05 min-w-0"}),
		),
	))
}

func TestPeoplePickerReferenceRoleFieldInTrailingColumn(t *testing.T) {
	role := peoplePickerStaticHTML(`<label class="visually-hidden" for="circle-ada-role">Role in this circle for Ada Okafor</label><input type="text" id="circle-ada-role" name="role_ada" value="Chair" placeholder="Optional role">`)
	assertPeoplePickerMatchesReference(t, 3, peoplePickerWithChildren(
		PeoplePicker(PeoplePickerProps{Legend: "Circle members"}),
		peoplePickerReferenceTable(`<th scope="col">Role in this circle</th>`,
			PeoplePickerRow(PeoplePickerRowProps{ID: "circle-ada", Name: "member_ids", Value: "ada", PersonName: "Ada Okafor", Checked: true, IdentityClass: "flex gap-05 min-w-0", TrailingLabel: "Role in this circle", Trailing: role}),
		),
	))
}

func TestPeoplePickerReferenceEmptyState(t *testing.T) {
	assertPeoplePickerMatchesReference(t, 5, PeoplePicker(PeoplePickerProps{Legend: "Members", Description: "No members are available."}))
}

func peoplePickerReferenceNodes(t *testing.T, fieldset int, tag, class string) []string {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(peoplePickerReferenceFieldsets(t)[fieldset]))
	if err != nil {
		t.Fatal(err)
	}
	var rendered []string
	for _, node := range findPeoplePickerNodes(doc, tag, func(n *html.Node) bool { return class == "" || hasPeoplePickerClass(n, class) }) {
		var buffer bytes.Buffer
		if err := html.Render(&buffer, node); err != nil {
			t.Fatal(err)
		}
		rendered = append(rendered, buffer.String())
	}
	return rendered
}

func peoplePickerOrganizerSwitch(id, name string, checked, disabled bool) templ.Component {
	attrs := ""
	if checked {
		attrs += " checked"
	}
	if disabled {
		attrs += " disabled"
	}
	return peoplePickerStaticHTML(`<label class="switch switch--small"><input type="checkbox" role="switch" name="organizer_ids" value="` + id + `" aria-label="Make ` + name + ` an organizer"` + attrs + `><span></span></label>`)
}

// Ed's reference row links aria-describedby to an empty detail span. Templ omits both when
// there is no detail, so that row is intentionally not compared.
func TestPeoplePickerReferenceRowsWithDetailsAvatarsAndOrganizers(t *testing.T) {
	rows := []PeoplePickerRowProps{
		{ID: "space-ada", Name: "member_ids", Value: "ada", PersonName: "Ada Okafor", Detail: "Organizer, Board", AvatarSrc: "/img/guide/ludd-avatar.jpg", Checked: true, Trailing: peoplePickerOrganizerSwitch("ada", "Ada Okafor", true, false)},
		{ID: "space-bo", Name: "member_ids", Value: "bo", PersonName: "Bo Lindqvist", Detail: "Board", Checked: true, Trailing: peoplePickerOrganizerSwitch("bo", "Bo Lindqvist", false, false)},
		{ID: "space-cy", Name: "member_ids", Value: "cy", PersonName: "Cy Moreno", Detail: "Removed", Trailing: peoplePickerOrganizerSwitch("cy", "Cy Moreno", false, true)},
		{ID: "space-di", Name: "member_ids", Value: "di", PersonName: "Di Ng", Detail: "Super admin", Checked: true, Disabled: true, Trailing: peoplePickerOrganizerSwitch("di", "Di Ng", false, true)},
	}
	want := peoplePickerReferenceNodes(t, 0, "tr", "")
	if len(want) < len(rows)+2 {
		t.Fatalf("reference has %d rows", len(want))
	}
	for i, row := range rows {
		row.ShowAvatar, row.TrailingLabel, row.IdentityClass = true, "Organizer", "flex gap-05 min-w-0"
		got := strings.ReplaceAll(testutil.CanonicalHTML(t, testutil.RenderToString(t, templ.Join(peoplePickerStaticHTML("<table><tbody>"), PeoplePickerRow(row), peoplePickerStaticHTML("</tbody></table>")))), "-selection-label", "-label")
		if exp := testutil.CanonicalHTML(t, "<table><tbody>"+want[i+1]+"</tbody></table>"); got != exp {
			t.Errorf("%s differs from the published row\ngot:\n%s\nwant:\n%s", row.ID, got, exp)
		}
	}
}

func TestPeoplePickerReferenceGroupsStatusAndActions(t *testing.T) {
	buttons := []templ.Component{
		peoplePickerStaticHTML(`<button type="button" class="button button--small" aria-pressed="false" data-demo-group="all" disabled>All</button>`),
		peoplePickerStaticHTML(`<button type="button" class="button button--small" aria-pressed="true" data-demo-group="board" disabled>Board</button>`),
		peoplePickerStaticHTML(`<button type="button" class="button button--small" aria-pressed="false" data-demo-group="officers" disabled>Officers</button>`),
		peoplePickerStaticHTML(`<button type="button" class="button button--small" aria-pressed="false" disabled>Members</button>`),
	}
	clear := peoplePickerStaticHTML(`<button type="button" class="button button--small button--ghost" name="picker_action" value="clear" disabled>Clear selection</button>`)
	slots := []struct {
		name, tag, class string
		got              templ.Component
	}{
		{"groups", "div", "people-picker-groups", peoplePickerWithChildren(PeoplePickerGroups(PeoplePickerGroupsProps{Class: "cluster gap-025", Attrs: templ.Attributes{"role": "group", "aria-label": "Circles"}}), buttons...)},
		{"status", "p", "people-picker-status", PeoplePickerStatus(PeoplePickerStatusProps{Class: "text-sm text-muted"})},
		{"actions", "div", "people-picker-actions", peoplePickerWithChildren(PeoplePickerActions(PeoplePickerActionsProps{Class: "cluster cluster-compact"}), clear)},
	}
	for _, slot := range slots {
		want := peoplePickerReferenceNodes(t, 0, slot.tag, slot.class)
		if len(want) != 1 {
			t.Fatalf("%s: reference has %d nodes", slot.name, len(want))
		}
		if got, exp := testutil.CanonicalHTML(t, testutil.RenderToString(t, slot.got)), testutil.CanonicalHTML(t, want[0]); got != exp {
			t.Errorf("%s differs from the published reference\ngot:\n%s\nwant:\n%s", slot.name, got, exp)
		}
	}
}

// The producer's truncated spans also carry its demo hook data-demo-truncate-text. That attribute belongs to
// the producer's manual script, not to the shared markup, so it is the one normalized difference.
func TestPeoplePickerReferenceTruncatedRow(t *testing.T) {
	want := peoplePickerReferenceNodes(t, 4, "tr", "")
	if len(want) != 2 {
		t.Fatalf("reference narrow picker has %d rows, want header and one body row", len(want))
	}
	const name = "Maximiliana Featherstonehaugh-Abernathy"
	row := PeoplePickerRowProps{
		ID: "narrow-long", Name: "member_ids", Value: "long", PersonName: name, Detail: "Organizer, Festival submissions working group",
		ShowAvatar: true, Checked: true, NameClass: "block truncate", DetailClass: "truncate", TrailingLabel: "Organizer",
		IdentityClass: "flex gap-05 min-w-0", Trailing: peoplePickerOrganizerSwitch("long", name, true, false),
	}
	got := strings.ReplaceAll(testutil.CanonicalHTML(t, testutil.RenderToString(t, templ.Join(peoplePickerStaticHTML("<table><tbody>"), PeoplePickerRow(row), peoplePickerStaticHTML("</tbody></table>")))), "-selection-label", "-label")
	exp := testutil.CanonicalHTML(t, "<table><tbody>"+strings.ReplaceAll(want[1], " data-demo-truncate-text=\"\"", "")+"</tbody></table>")
	if got != exp {
		t.Errorf("truncated row differs from the published reference\ngot:\n%s\nwant:\n%s", got, exp)
	}
}
