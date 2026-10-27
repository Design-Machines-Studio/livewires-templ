package component

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

func mentionAttr(t *testing.T, node *html.Node, key, want string) {
	t.Helper()
	if node == nil {
		t.Fatalf("missing node for %s", key)
	}
	if got, ok := testutil.AttrVal(node, key); !ok || got != want {
		t.Errorf("%s = %q (%v), want %q", key, got, ok, want)
	}
}

func TestMentionStates(t *testing.T) {
	for _, state := range []string{"closed", "loading", "ready", "empty", "error", "", "invalid"} {
		for _, busy := range []bool{false, true} {
			t.Run(state+"/"+boolStr(busy), func(t *testing.T) {
				out := testutil.RenderToString(t, withChildren(MentionSuggestionsComponent(MentionSuggestionsProps{ID: "people", Label: "People", State: state, Busy: busy, Status: "Caller status"}), MentionOption("people-ada", "Ada", "", false)))
				nodes := testutil.ParseFragment(t, out)
				root, list := testutil.FindElement(nodes, "div"), testutil.FindElement(nodes, "ul")
				want := state
				if state == "" || state == "invalid" {
					want = "closed"
				}
				mentionAttr(t, root, "data-state", want)
				mentionAttr(t, root, "class", "rounded scheme-default")
				mentionAttr(t, root, "dir", "ltr")
				mentionAttr(t, list, "role", "listbox")
				mentionAttr(t, list, "id", "people")
				mentionAttr(t, list, "aria-label", "People")
				_, isBusy := testutil.AttrVal(list, "aria-busy")
				if isBusy != busy {
					t.Error("aria-busy does not follow Busy")
				}
				if busy {
					mentionAttr(t, list, "aria-busy", "true")
				}
				_, hidden := testutil.AttrVal(list, "hidden")
				visible := state == "ready" || (state == "loading" && busy)
				if hidden == visible || (testutil.FindElement(nodes, "li") != nil) != visible {
					t.Error("wrong row/list visibility")
				}
				if strings.Count(out, `role="status"`) != 1 {
					t.Fatal("expected exactly one status region")
				}
				mentionAttr(t, testutil.FindElement(nodes, "p"), "aria-atomic", "true")
				if testutil.NodeText(testutil.FindElement(nodes, "p")) != "Caller status" {
					t.Error("status must be caller supplied")
				}
				if strings.Contains(out, "aria-expanded") || strings.Contains(out, "tabindex") {
					t.Error("popup must not emit editor semantics/focus")
				}
			})
		}
	}
}

func TestMentionOptionSemantics(t *testing.T) {
	for _, active := range []bool{false, true} {
		out := testutil.RenderToString(t, MentionOptionComponent(MentionOptionProps{ID: "list-person", Name: "Full name", Detail: "Real circle", Active: active, Disabled: true, Truncate: true, Avatar: &AvatarProps{Name: "Full name"}}))
		nodes := testutil.ParseFragment(t, out)
		row := testutil.FindElement(nodes, "li")
		mentionAttr(t, row, "role", "option")
		mentionAttr(t, row, "id", "list-person")
		mentionAttr(t, row, "aria-selected", boolStr(active))
		mentionAttr(t, row, "aria-disabled", "true")
		_, marked := testutil.AttrVal(row, "data-active")
		if marked != active {
			t.Error("data-active mismatch")
		}
		class, _ := testutil.AttrVal(row, "class")
		if active && class != "rounded-sm scheme-subtle" || !active && class != "" {
			t.Errorf("active classes = %q", class)
		}
		mentionAttr(t, testutil.FindElementByClass(nodes, "avatar"), "aria-hidden", "true")
		if !strings.Contains(out, `<span class="initials">FN</span>`) {
			t.Error("expected decorative initials")
		}
		if strings.Count(out, `class="truncate"`) != 2 || !strings.Contains(out, "Full name") || !strings.Contains(out, "Real circle") {
			t.Error("truncate must preserve full name/detail")
		}
	}
	image := testutil.RenderToString(t, MentionOptionComponent(MentionOptionProps{Name: "Ada", Avatar: &AvatarProps{Name: "Ada", Src: "/ada.jpg", ShowName: true, Size: "xl"}}))
	nodes := testutil.ParseFragment(t, image)
	mentionAttr(t, testutil.FindElement(nodes, "img"), "alt", "")
	mentionAttr(t, testutil.FindElementByClass(nodes, "avatar"), "aria-hidden", "true")
	if strings.Contains(image, "mention-detail") || strings.Contains(image, "initials") || strings.Contains(image, "avatar--") || strings.Contains(image, "cluster") {
		t.Error("unexpected detail/avatar/layout markup")
	}
	plain := testutil.RenderToString(t, MentionOption("id", "Ada", "", false))
	if strings.Contains(plain, "avatar") || strings.Contains(plain, "mention-detail") || strings.Contains(plain, "aria-disabled") || strings.Contains(plain, "class=") {
		t.Error("optional markup should be absent")
	}
}

func TestMentionAttrsEscapingAndIDs(t *testing.T) {
	if MentionOptionID("mentions-list", "sam-2") != "mentions-list-sam-2" || MentionOptionID("rtl", "person_123") != "rtl-person_123" {
		t.Error("stable id scheme changed")
	}
	if MentionOptionID("a", "same") == MentionOptionID("b", "same") {
		t.Error("list ids must scope option ids")
	}
	name, detail, status := `<img src=x onerror="x"> & Ada`, `<b>Kitchen & finance</b>`, `<script>alert('x')</script>`
	row := MentionOptionComponent(MentionOptionProps{ID: "list-person", Name: name, Detail: detail, Active: true, Class: "caller-row", Attrs: templ.Attributes{"role": "button", "id": "wrong", "aria-selected": "false", "class": "wrong", "data-owner": "row"}})
	out := testutil.RenderToString(t, withChildren(MentionSuggestionsComponent(MentionSuggestionsProps{ID: "list", Label: "People", State: "ready", Status: status, Dir: "rtl", Class: "caller-root", Attrs: templ.Attributes{"class": "wrong", "data-state": "error", "dir": "ltr", "data-owner": "root"}}), row))
	nodes := testutil.ParseFragment(t, out)
	root := testutil.FindElement(nodes, "div")
	li := testutil.FindElement(nodes, "li")
	mentionAttr(t, root, "class", "rounded scheme-default caller-root")
	mentionAttr(t, root, "dir", "rtl")
	mentionAttr(t, root, "data-state", "ready")
	mentionAttr(t, root, "data-owner", "root")
	mentionAttr(t, li, "class", "rounded-sm scheme-subtle caller-row")
	mentionAttr(t, li, "role", "option")
	mentionAttr(t, li, "id", "list-person")
	mentionAttr(t, li, "aria-selected", "true")
	mentionAttr(t, li, "data-owner", "row")
	for _, tag := range []string{"div", "li"} {
		raw := testutil.RawTag(t, out, tag)
		if strings.Index(raw, `class="wrong"`) < strings.Index(raw, `class="rounded`) {
			t.Error("caller class precedes component class")
		}
	}
	if testutil.FindElement(nodes, "img") != nil || testutil.FindElement(nodes, "script") != nil || testutil.FindElement(nodes, "b") != nil {
		t.Error("unescaped text")
	}
	if (!strings.Contains(testutil.NodeText(root), name) || !strings.Contains(testutil.NodeText(root), detail)) || testutil.NodeText(testutil.FindElement(nodes, "p")) != status {
		t.Error("escaped text changed content")
	}
}

// Golden source is the unmodified producer v2.27.0 reference, not generated
// from this component. Compare the ready subtree byte-for-byte after canonical
// serialization: only source indentation, attribute order and boolean spelling
// are normalized; every tag, attribute, class and text value is compared.
func TestMentionReadyProducerGolden(t *testing.T) {
	b, err := os.ReadFile("testdata/livewires/mention-suggestions-reference.html")
	if err != nil {
		t.Fatal(err)
	}
	ref := testutil.FindElementByID(testutil.ParseFragment(t, string(b)), "r-ready")
	if ref == nil {
		t.Fatal("producer ready root missing")
	}
	var expected bytes.Buffer
	if err := html.Render(&expected, ref); err != nil {
		t.Fatal(err)
	}
	rows := []templ.Component{}
	for i, p := range [][3]string{{"ada", "Ada Okafor", "Finance circle"}, {"bo", "Bo Lindqvist", "Kitchen"}, {"sam-k", "Sam Rivera", "Kitchen"}, {"sam-f", "Sam Rivera", "Finance circle"}, {"di", "Di Ng", ""}} {
		props := MentionOptionProps{ID: MentionOptionID("r-ready", p[0]), Name: p[1], Detail: p[2], Active: i == 0, Disabled: i == 3}
		if i < 4 {
			props.Avatar = &AvatarProps{Name: p[1]}
		}
		rows = append(rows, MentionOptionComponent(props))
	}
	out := testutil.RenderToString(t, withChildren(MentionSuggestionsComponent(MentionSuggestionsProps{ID: "r-ready-list", Label: "People", State: "ready", Status: "5 people found", Attrs: templ.Attributes{"id": "r-ready", "style": "position: static; --mention-suggestions-max-block: none"}}), rows...))
	// Reference omits dir (implicit ltr); the API explicitly emits it as required
	// by the attribute contract. Supply that same declared direction to the fixture.
	ref.Attr = append(ref.Attr, html.Attribute{Key: "dir", Val: "ltr"})
	expected.Reset()
	html.Render(&expected, ref)
	if got, want := testutil.CanonicalHTML(t, out), testutil.CanonicalHTML(t, expected.String()); got != want {
		t.Fatalf("producer ready mismatch\nGOT\n%s\nWANT\n%s", got, want)
	}
}
