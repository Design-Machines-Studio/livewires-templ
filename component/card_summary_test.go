package component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

// Reference fixture card-summary-reference.html is the summary card section of
// livewires public/reference/components/cards.html at v2.26.0
// (d4c197b29b1675bf6cada14d051b7a8e1425fe2f).

const badgeRow = "cluster cluster-between cluster-nowrap items-start"

const (
	refAvatar = "/img/guide/ludd-avatar.jpg"

	chatsIcon  = `<svg class="icon flex-none" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor"><path d="M216,76H188V48a20,20,0,0,0-20-20H40A20,20,0,0,0,20,48V176a12,12,0,0,0,19.54,9.33l28.46-23V184a20,20,0,0,0,20,20h92.17l36.29,29.33A12,12,0,0,0,236,224V96A20,20,0,0,0,216,76ZM44,150.87V52H164v80H71.58A12,12,0,0,0,64,134.67Zm168,48-20-16.2a12,12,0,0,0-7.54-2.67H92V156h76a20,20,0,0,0,20-20V100h24Z"></path></svg>`
	clockIcon  = `<svg class="icon flex-none" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor"><path d="M128,20A108,108,0,1,0,236,128,108.12,108.12,0,0,0,128,20Zm0,192a84,84,0,1,1,84-84A84.09,84.09,0,0,1,128,212Zm68-84a12,12,0,0,1-12,12H128a12,12,0,0,1-12-12V72a12,12,0,0,1,24,0v44h44A12,12,0,0,1,196,128Z"></path></svg>`
	peopleIcon = `<svg class="icon flex-none" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor"><path d="M164.38,181.1a52,52,0,1,0-72.76,0,75.89,75.89,0,0,0-30,28.89,12,12,0,0,0,20.78,12,53,53,0,0,1,91.22,0,12,12,0,1,0,20.78-12A75.89,75.89,0,0,0,164.38,181.1ZM100,144a28,28,0,1,1,28,28A28,28,0,0,1,100,144Zm147.21,9.59a12,12,0,0,1-16.81-2.39c-8.33-11.09-19.85-19.59-29.33-21.64a12,12,0,0,1-1.82-22.91,20,20,0,1,0-24.78-28.3,12,12,0,1,1-21-11.6,44,44,0,1,1,73.28,48.35,92.18,92.18,0,0,1,22.85,21.69A12,12,0,0,1,247.21,153.59Zm-192.28-24c-9.48,2.05-21,10.55-29.33,21.65A12,12,0,0,1,6.41,136.79,92.37,92.37,0,0,1,29.26,115.1a44,44,0,1,1,73.28-48.35,12,12,0,1,1-21,11.6,20,20,0,1,0-24.78,28.3,12,12,0,0,1-1.82,22.91Z"></path></svg>`
)

func summaryCard(props CardSummaryProps, children ...templ.Component) templ.Component {
	props.Scheme = "subtle"
	props.Class = "rounded"
	props.HeaderClass = "stack"
	return withChildren(CardSummaryComponent(props), children...)
}

func summaryFooter(class string, children ...templ.Component) templ.Component {
	return withChildren(CardFooterComponent(CardFooterProps{Class: class}), children...)
}

func people(n int) []CardPerson {
	out := make([]CardPerson, n)
	for i := range out {
		out[i] = CardPerson{Name: "Ludd", Src: refAvatar}
	}
	return out
}

func membersFooter(extra, label string, n int, count templ.Component) templ.Component {
	return summaryFooter("cluster cluster-between items-center text-sm",
		CardPeopleComponent(CardPeopleProps{People: people(n), Extra: extra, Label: label, Class: "cluster cluster-compact items-center"}),
		count,
	)
}

func conversations(value, label, class string) templ.Component {
	return CardCountComponent(CardCountProps{Value: value, Label: label, Icon: staticHTML(chatsIcon), Class: class})
}

func authorFooter(name, datetime, label string) templ.Component {
	return summaryFooter("cluster cluster-compact cluster-nowrap items-center text-sm",
		withChildren(CardAuthorComponent(CardAuthorProps{Name: name, Src: refAvatar, Class: "min-w-0"}),
			CardDateComponent(CardDateProps{Datetime: datetime, Label: label, Prefix: "Edited", Class: "block"}),
		),
	)
}

func authorMembersFooter(name, datetime, label, count string) templ.Component {
	return summaryFooter("cluster cluster-between cluster-nowrap items-center text-sm",
		withChildren(cardWrapperComponent("cluster cluster-compact cluster-nowrap items-center min-w-0"),
			withChildren(CardAuthorComponent(CardAuthorProps{Name: name, Src: refAvatar, Class: "min-w-0 text-sm"}),
				CardDateComponent(CardDateProps{Datetime: datetime, Label: label, HiddenLabel: "Last updated", Icon: staticHTML(clockIcon), Class: "flex items-center gap-025"}),
			),
		),
		CardCountComponent(CardCountProps{Value: count, Label: "members", Icon: staticHTML(peopleIcon), Class: "inline-flex items-center gap-025 flex-shrink-0"}),
	)
}

func cardWrapperComponent(class string) templ.Component { return cardWrapper(class, nil) }

func listOf(class string, items ...templ.Component) templ.Component {
	return withChildren(cardWrapperTag("ul", class), items...)
}

func cardWrapperTag(tag, class string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "<"+tag+` class="`+class+`">`); err != nil {
			return err
		}
		children := templ.GetChildren(ctx)
		if err := children.Render(templ.ClearChildren(ctx), w); err != nil {
			return err
		}
		_, err := io.WriteString(w, "</"+tag+">")
		return err
	})
}

func TestReferenceSummaryCards(t *testing.T) {
	long := "Festival submissions 2027: shortlist, deadlines, fees and the long tail of paperwork"
	longDescription := "Where we track which festivals get which films, what each submission costs, who owns the paperwork, and which deadlines are hard. Fee waivers, screener links and rights questions live here too. If a deadline is inside two weeks, mention it in the title."
	patronage := "Revise the patronage formula for project-based work"
	patronageDescription := "Weights hours on client projects and internal work the same way, and phases the change in over two allocation years."
	fourDay := "Adopt a four-day week for the production team"
	fourDayDescription := "Moves the production team to four eight-hour days from January, with a review at the spring general meeting."

	grid1 := listOf("grid list-none",
		summaryCard(CardSummaryProps{Title: "Board room", Href: "#", Description: "Board and officer preparation between meetings. Decisions still go through proposals."},
			membersFooter("+2", "5 members", 3, conversations("4", "conversations, 1 unread", "inline-flex items-center gap-025 font-bold"))),
		summaryCard(CardSummaryProps{Title: long, Href: "#", Description: longDescription},
			membersFooter("+1", "4 members", 3, conversations("12", "conversations", "inline-flex items-center gap-025"))),
		summaryCard(CardSummaryProps{Title: "Studio commons", Href: "#"},
			membersFooter("+ All", "All members", 3, conversations("9", "conversations", "inline-flex items-center gap-025"))),
		summaryCard(CardSummaryProps{Title: patronage, Href: "#", Description: patronageDescription},
			membersFooter("", "2 collaborators", 2, CardDateComponent(CardDateProps{Datetime: "2026-10-17", Label: "2 days ago", Prefix: "Edited"}))),
	)
	grid2 := listOf("grid list-none",
		summaryCard(CardSummaryProps{Title: patronage, Href: "#", Description: patronageDescription}, authorFooter("Ada Okafor", "2026-10-17", "2 days ago")),
		summaryCard(CardSummaryProps{Title: fourDay, Href: "#", Description: fourDayDescription}, authorFooter("Mateus Rodrigues", "2026-10-18", "yesterday")),
		summaryCard(CardSummaryProps{Title: "Set a member loan rate", Href: "#"}, authorFooter("Wiremu Te Rangi-Hēnare-Whitehead", "2026-09-28", "3 weeks ago")),
	)
	grid3 := listOf("grid list-none",
		summaryCard(CardSummaryProps{Title: patronage, Href: "#", Description: patronageDescription, Badge: "New", HeadingRowClass: badgeRow}, authorMembersFooter("Ada Okafor", "2026-10-17", "Oct 17", "5")),
		summaryCard(CardSummaryProps{Title: fourDay, Href: "#", Description: fourDayDescription, Badge: "Updated", HeadingRowClass: badgeRow}, authorMembersFooter("Wiremu Te Rangi-Hēnare-Whitehead", "2026-10-18", "Yesterday", "12")),
		summaryCard(CardSummaryProps{Title: "Set a member loan rate", Href: "#", HeadingRowClass: badgeRow}, authorMembersFooter("Mateus Rodrigues", "2026-09-28", "Sep 28", "3")),
	)
	assertMatchesReference(t, "card-summary-reference.html", templ.Join(grid1, grid2, grid3))
}

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	return testutil.RenderToString(t, c)
}

func TestCardFooterEmitsFooterRegion(t *testing.T) {
	out := render(t, withChildren(CardFooterComponent(CardFooterProps{Class: "cluster cluster-between", Attrs: templ.Attributes{"data-x": "1"}}), staticHTML("<span>kid</span>")))
	tag := testutil.RawTag(t, out, "footer")
	for _, want := range []string{`class="footer box cluster cluster-between"`, `data-x="1"`} {
		if !strings.Contains(tag, want) {
			t.Errorf("footer tag %q missing %q", tag, want)
		}
	}
	if !strings.Contains(out, "<span>kid</span></footer>") {
		t.Errorf("children not rendered inside footer: %s", out)
	}
}

func TestCardFooterBakesNoLayoutPrimitives(t *testing.T) {
	out := render(t, CardFooterComponent(CardFooterProps{}))
	if !strings.Contains(out, `class="footer box"`) {
		t.Errorf("expected only footer box, got %s", out)
	}
	for _, banned := range []string{"cluster", "stack"} {
		if strings.Contains(out, banned) {
			t.Errorf("footer must not bake %q: %s", banned, out)
		}
	}
}

func TestCardFooterWrapper(t *testing.T) {
	out := render(t, withChildren(CardFooter("text-sm"), staticHTML("x")))
	if !strings.Contains(out, `<footer class="footer box text-sm">x</footer>`) {
		t.Errorf("unexpected wrapper output: %s", out)
	}
}

func TestCardFooterAfterHeaderStaysSiblings(t *testing.T) {
	out := render(t, withChildren(CardFrameComponent(CardFrameProps{Scheme: "subtle", Class: "rounded"}),
		CardHeaderComponent(CardHeaderProps{Title: "Panel"}),
		CardFooterComponent(CardFooterProps{Class: "text-sm"}),
	))
	if !strings.Contains(out, `<div class="card scheme-subtle rounded">`) {
		t.Errorf("frame must be a bare card: %s", out)
	}
	if !strings.Contains(out, `</div><footer class="footer box text-sm">`) {
		t.Errorf("footer must follow the header directly: %s", out)
	}
	if strings.Contains(out, `card box`) {
		t.Errorf("frame must not add box padding: %s", out)
	}
}

func TestCardHeaderComponentHeadingLevelAndEscaping(t *testing.T) {
	out := render(t, CardHeaderComponent(CardHeaderProps{Title: `A <b>&"x"`, HeadingLevel: "h2", Class: "extra", Attrs: templ.Attributes{"data-h": "1"}}))
	if !strings.Contains(out, `class="header box extra"`) || !strings.Contains(out, `data-h="1"`) {
		t.Errorf("class/attrs not passed: %s", out)
	}
	if !strings.Contains(out, "<h2>A &lt;b&gt;&amp;&#34;x&#34;</h2>") {
		t.Errorf("title not escaped in h2: %s", out)
	}
}

func TestCardFrame(t *testing.T) {
	out := render(t, withChildren(CardFrame(""), staticHTML("c")))
	if out != `<div class="card">c</div>` {
		t.Errorf("unexpected frame: %s", out)
	}
}

func TestCardSummaryStructure(t *testing.T) {
	out := render(t, withChildren(CardSummaryComponent(CardSummaryProps{
		Title: "Board room", Href: "/spaces/board", Description: "Prep.", Scheme: "subtle", Class: "rounded",
		HeaderClass: "stack", Attrs: templ.Attributes{"data-id": "7", "style": "--card-description-lines: 2"},
	}), staticHTML("<footer>f</footer>")))
	nodes := testutil.ParseFragment(t, out)
	li := testutil.FindElement(nodes, "li")
	if li == nil {
		t.Fatalf("expected li: %s", out)
	}
	if got, _ := testutil.AttrVal(li, "class"); got != "card card--summary scheme-subtle rounded" {
		t.Errorf("card class = %q", got)
	}
	if got, _ := testutil.AttrVal(li, "data-id"); got != "7" {
		t.Errorf("Attrs not passed through: %s", out)
	}
	if got, _ := testutil.AttrVal(li, "style"); got != "--card-description-lines: 2" {
		t.Errorf("style custom property not passed through: %q", got)
	}
	h3 := testutil.FindElement(nodes, "h3")
	if got, _ := testutil.AttrVal(h3, "class"); got != "card-heading" {
		t.Errorf("heading class = %q", got)
	}
	a := testutil.FindElement(nodes, "a")
	if got, _ := testutil.AttrVal(a, "href"); got != "/spaces/board" {
		t.Errorf("heading link href = %q", got)
	}
	if strings.Count(out, "<a ") != 1 || strings.Contains(out, "<button") {
		t.Errorf("heading link must be the only interactive element: %s", out)
	}
	header := testutil.FindElementByClass(nodes, "header")
	if got, _ := testutil.AttrVal(header, "class"); got != "header box stack" {
		t.Errorf("header class = %q", got)
	}
	if !strings.Contains(out, `<p class="description">Prep.</p></div><footer>f</footer></li>`) {
		t.Errorf("footer must follow the header inside the card: %s", out)
	}
}

func TestCardSummaryHeadingLevels(t *testing.T) {
	for level, tag := range map[string]string{"": "h3", "h2": "h2", "h4": "h4", "h6": "h6", "bogus": "h3"} {
		out := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Href: "/x", HeadingLevel: level}))
		if !strings.Contains(out, "<"+tag+` class="card-heading"><a href="/x">T</a></`+tag+">") {
			t.Errorf("level %q: expected %s heading, got %s", level, tag, out)
		}
	}
}

func TestCardSummaryEscaping(t *testing.T) {
	out := render(t, CardSummaryComponent(CardSummaryProps{
		Title: `<script>alert(1)</script>`, Href: `javascript:alert(1)`, Description: `a & <i>b</i>`, Badge: `<b>x</b>`,
	}))
	for _, bad := range []string{"<script>", "<i>", "<b>x</b>", `href="javascript:`} {
		if strings.Contains(out, bad) {
			t.Errorf("unescaped %q in %s", bad, out)
		}
	}
	if !strings.Contains(out, "&lt;script&gt;") || !strings.Contains(out, "a &amp; &lt;i&gt;b&lt;/i&gt;") {
		t.Errorf("expected escaped text: %s", out)
	}
}

func TestCardSummaryWithoutHrefHasNoLink(t *testing.T) {
	out := render(t, CardSummaryComponent(CardSummaryProps{Title: "Plain"}))
	if strings.Contains(out, "<a") {
		t.Errorf("no link expected without Href: %s", out)
	}
	if !strings.Contains(out, `<h3 class="card-heading">Plain</h3>`) {
		t.Errorf("heading text missing: %s", out)
	}
}

// sameHTML compares markup structurally. templ keeps one space between inline
// siblings written on separate lines, which does not change how Live Wires
// lays out the footer.
func sameHTML(t *testing.T, got, want string) {
	t.Helper()
	if g, w := testutil.CanonicalHTML(t, got), testutil.CanonicalHTML(t, want); g != w {
		t.Errorf("markup mismatch\n--- got ---\n%s\n--- want ---\n%s", g, w)
	}
}

func TestCardSummaryOptionalBadgeAndDescription(t *testing.T) {
	plain := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Href: "/x"}))
	if strings.Contains(plain, "badge") || strings.Contains(plain, "description") || strings.Contains(plain, "<div class=\"cluster") {
		t.Errorf("badge, description and heading row must be omitted: %s", plain)
	}
	if strings.Contains(plain, "min-w-0") {
		t.Errorf("min-w-0 belongs only to the heading inside a heading row: %s", plain)
	}
	badged := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Href: "/x", Badge: "Updated", HeadingRowClass: badgeRow}))
	sameHTML(t, badged, `<li class="card card--summary"><div class="header box"><div class="`+badgeRow+`"><h3 class="card-heading min-w-0"><a href="/x">T</a></h3><span class="badge badge--blue">Updated</span></div></div></li>`)
	rowOnly := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Href: "/x", HeadingRowClass: badgeRow}))
	sameHTML(t, rowOnly, `<li class="card card--summary"><div class="header box"><div class="`+badgeRow+`"><h3 class="card-heading min-w-0"><a href="/x">T</a></h3></div></div></li>`)
	other := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Badge: "New", BadgeVariant: "green"}))
	if !strings.Contains(other, "badge badge--green") {
		t.Errorf("badge variant ignored: %s", other)
	}
}

func TestCardSummaryOptionalFooter(t *testing.T) {
	out := render(t, CardSummaryComponent(CardSummaryProps{Title: "T"}))
	if strings.Contains(out, "footer") {
		t.Errorf("no footer expected without a footer child: %s", out)
	}
}

func TestCardSummaryBakesNoLayoutPrimitives(t *testing.T) {
	out := render(t, CardSummaryComponent(CardSummaryProps{Title: "T", Href: "/x", Description: "d", Badge: "New"}))
	for _, banned := range []string{"stack", "cluster", `box stack`, "rounded"} {
		if strings.Contains(out, banned) {
			t.Errorf("summary card must not bake %q: %s", banned, out)
		}
	}
}

func TestCardSummaryWrapper(t *testing.T) {
	out := render(t, withChildren(CardSummary("T", "/x", "d", "subtle", "rounded", "stack"), staticHTML("<footer>f</footer>")))
	if !strings.Contains(out, `<li class="card card--summary scheme-subtle rounded">`) || !strings.Contains(out, `<div class="header box stack">`) || !strings.Contains(out, `<footer>f</footer></li>`) {
		t.Errorf("unexpected wrapper output: %s", out)
	}
}

func TestCardPeopleDecorativeAvatars(t *testing.T) {
	out := render(t, CardPeopleComponent(CardPeopleProps{
		People: []CardPerson{{Name: "Ada Okafor", Src: "/a.jpg"}, {Name: "Bo Li"}, {Name: "Cy Ng"}, {Name: "Di Ox"}},
		Extra:  "+2", Label: "5 members", Class: "cluster cluster-compact items-center",
	}))
	nodes := testutil.ParseFragment(t, out)
	group := testutil.FindElementByClass(nodes, "avatar-group")
	if v, _ := testutil.AttrVal(group, "aria-hidden"); v != "true" {
		t.Errorf("avatar group must be aria-hidden: %s", out)
	}
	if n := strings.Count(out, `class="avatar avatar--sm"`); n != 3 {
		t.Errorf("expected at most 3 avatars, got %d: %s", n, out)
	}
	if strings.Contains(out, "role=") || strings.Contains(out, "aria-label") || strings.Contains(out, "Di Ox") {
		t.Errorf("avatars must have no accessible name or fourth avatar: %s", out)
	}
	if !strings.Contains(out, `<img src="/a.jpg" alt="">`) {
		t.Errorf("image avatar must use empty alt: %s", out)
	}
	if !strings.Contains(out, `<span class="initials" aria-hidden="true">B</span>`) {
		t.Errorf("initials avatar should show a single hidden initial: %s", out)
	}
	if !strings.Contains(out, `<span aria-hidden="true">+2</span>`) || !strings.Contains(out, `<span class="visually-hidden">5 members</span>`) {
		t.Errorf("overflow must be hidden and count visually hidden: %s", out)
	}
	if !strings.HasPrefix(out, `<div class="cluster cluster-compact items-center">`) {
		t.Errorf("wrapper class missing: %s", out)
	}
}

func TestCardPeopleOmitsEmptyParts(t *testing.T) {
	out := render(t, CardPeopleComponent(CardPeopleProps{Label: "No one yet"}))
	sameHTML(t, out, `<div><span class="visually-hidden">No one yet</span></div>`)
}

func TestCardPeopleWrapperAndAttrs(t *testing.T) {
	out := render(t, CardPeopleComponent(CardPeopleProps{People: people(1), Attrs: templ.Attributes{"data-p": "1"}}))
	if !strings.HasPrefix(out, `<div data-p="1">`) {
		t.Errorf("Attrs not passed: %s", out)
	}
	if got := render(t, CardPeople(people(1), "", "1 member", "c")); !strings.HasPrefix(got, `<div class="c">`) {
		t.Errorf("wrapper class: %s", got)
	}
}

func TestCardCountVisuallyHiddenText(t *testing.T) {
	out := render(t, CardCountComponent(CardCountProps{
		Value: "4", Label: "conversations, 1 unread", Icon: staticHTML(chatsIcon), Class: "inline-flex font-bold", Attrs: templ.Attributes{"data-c": "1"},
	}))
	if !strings.HasPrefix(out, `<span class="inline-flex font-bold" data-c="1"><svg class="icon flex-none" aria-hidden="true"`) {
		t.Errorf("icon/class/attrs: %s", out)
	}
	if !strings.Contains(out, `<span>4</span>`) || !strings.Contains(out, `<span class="visually-hidden"> conversations, 1 unread</span></span>`) {
		t.Errorf("value and hidden label: %s", out)
	}
	sameHTML(t, render(t, CardCount("3", "", "x")), `<span class="x"><span>3</span></span>`)
	if got := render(t, CardCount("3", "", "")); strings.Contains(got, "class=") || !strings.HasPrefix(got, "<span") {
		t.Errorf("empty Class must omit the class attribute: %s", got)
	}
}

func TestCardDateTimeElement(t *testing.T) {
	sentence := render(t, CardDateComponent(CardDateProps{Datetime: "2026-10-17", Label: "2 days ago", Prefix: "Edited"}))
	sameHTML(t, sentence, `<span class="text-xs text-muted">Edited <time datetime="2026-10-17">2 days ago</time></span>`)
	clock := render(t, CardDateComponent(CardDateProps{Datetime: "2026-10-17", Label: "Oct 17", HiddenLabel: "Last updated", Icon: staticHTML(clockIcon), Class: "flex items-center gap-025", Attrs: templ.Attributes{"data-d": "1"}}))
	sameHTML(t, clock, `<span class="flex items-center gap-025 text-xs text-muted" data-d="1">`+clockIcon+`<span class="visually-hidden">Last updated </span><time datetime="2026-10-17">Oct 17</time></span>`)
	if got := render(t, CardDate("", "Soon", "")); strings.Contains(got, "datetime") {
		t.Errorf("empty Datetime must omit the attribute: %s", got)
	}
	if got := render(t, CardDate("2026-10-17", `<b>x</b>`, "")); strings.Contains(got, "<b>") {
		t.Errorf("label must be escaped: %s", got)
	}
}

func TestCardAuthorKeepsFullNameAndHidesAvatar(t *testing.T) {
	name := "Wiremu Te Rangi-Hēnare-Whitehead"
	out := render(t, withChildren(CardAuthorComponent(CardAuthorProps{Name: name, Src: "/w.jpg", Class: "min-w-0", Attrs: templ.Attributes{"data-a": "1"}}), staticHTML("<i>date</i>")))
	want := `<span class="avatar"><img src="/w.jpg" alt=""></span><div class="min-w-0" data-a="1"><span class="block truncate"><span class="visually-hidden">By </span>` + name + `</span><i>date</i></div>`
	sameHTML(t, out, want)
	initials := render(t, CardAuthor("Ada Okafor", "", ""))
	if !strings.Contains(initials, `<span class="avatar"><span class="initials" aria-hidden="true">AO</span></span>`) || strings.Contains(initials, "aria-label") || strings.Contains(initials, "role=") {
		t.Errorf("initials avatar must be hidden from assistive technology: %s", initials)
	}
	custom := render(t, CardAuthorComponent(CardAuthorProps{Name: "Ada", HiddenPrefix: "Par"}))
	if !strings.Contains(custom, `<span class="visually-hidden">Par </span>Ada`) {
		t.Errorf("custom hidden prefix: %s", custom)
	}
}

func TestCardWithHeaderOutputUnchanged(t *testing.T) {
	out := render(t, withChildren(CardWithHeader(CardProps{Title: "Panel", Scheme: "subtle", Class: "rounded", Attrs: templ.Attributes{"data-k": "1"}}), staticHTML("<p>body</p>")))
	want := `<div class="card scheme-subtle rounded" data-k="1"><div class="header box"><h3>Panel</h3></div><div class="box"><p>body</p></div></div>`
	if out != want {
		t.Errorf("CardWithHeader output changed:\n got %s\nwant %s", out, want)
	}
	if strings.Contains(out, "footer") {
		t.Errorf("no footer expected: %s", out)
	}
}
