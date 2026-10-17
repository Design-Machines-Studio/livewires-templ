package component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

// Static test fixtures model already-rendered caller components, not an API for
// accepting HTML strings from application input.
func commentFixture(markup string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, markup)
		return err
	})
}

func commentByClass(t *testing.T, output, class string) *html.Node {
	t.Helper()
	var found *html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if classes, _ := testutil.AttrVal(n, "class"); strings.Contains(" "+classes+" ", " "+class+" ") {
			found = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, n := range testutil.ParseFragment(t, output) {
		walk(n)
	}
	if found == nil {
		t.Fatalf("missing .%s in %s", class, output)
	}
	return found
}

func TestCommentArticleEscapesConvenienceText(t *testing.T) {
	output := testutil.RenderToString(t, CommentArticle(`<b>Lee & Morgan</b>`, `<script>alert(1)</script>`, "2026-03-15", "15 March <today>"))
	nodes := testutil.ParseFragment(t, output)
	for _, tag := range []string{"script", "b", "today"} {
		if testutil.FindElement(nodes, tag) != nil {
			t.Fatalf("unescaped %s: %s", tag, output)
		}
	}
	for _, want := range []string{"&lt;script&gt;", "&lt;b&gt;Lee &amp; Morgan", "15 March &lt;today&gt;"} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing escaped value %q: %s", want, output)
		}
	}
	if strings.Contains(output, "Mar 15, 2026") {
		t.Fatal("article reformatted the caller's timestamp")
	}
}

func TestLegacyCommentCompatibility(t *testing.T) {
	for input, visible := range map[string]string{"2026-03-15": "Mar 15, 2026", "2026-03-15T09:30:00-07:00": "Mar 15, 2026", "unknown <date>": "unknown &lt;date&gt;", "": ""} {
		output := testutil.RenderToString(t, Comment("Lee", "<script>plain</script>", input, "Member"))
		if other := testutil.RenderToString(t, CommentComponent(CommentProps{Author: "Lee", Body: "<script>plain</script>", Date: input, Badge: "Member"})); output != other {
			t.Fatal("legacy wrapper and Props diverged")
		}
		if !strings.Contains(output, "<time>"+visible+"</time>") || !strings.Contains(output, "&lt;script&gt;plain&lt;/script&gt;") {
			t.Fatalf("legacy date/body changed: %s", output)
		}
		if !strings.Contains(output, "comment cluster items-start text-sm") || !strings.Contains(output, "comment-body flex-1 stack stack-compact") || !strings.Contains(output, "badge--small") {
			t.Fatalf("legacy presentation changed: %s", output)
		}
	}
	output := testutil.RenderToString(t, CommentComponent(CommentProps{Size: "lg", Align: "right", Class: "caller", Attrs: templ.Attributes{"data-comment": "legacy"}}))
	for _, want := range []string{"text-lg", "text-right", "caller", `data-comment="legacy"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("legacy Props lost %s", want)
		}
	}
}

func TestCommentArticleStructureAndSlots(t *testing.T) {
	props := CommentArticleProps{
		ID: "post-815", Class: "caller", Attrs: templ.Attributes{"data-message": "815"},
		Author: commentFixture(`<a class="author" href="/members/sam">Sam</a>`),
		Avatar: commentFixture(`<img class="avatar" src="/sam.jpg" alt="">`), AvatarHref: "/members/sam", AvatarLabel: "Sam's profile",
		DateTime: "2026-03-15T10:05:00-07:00", Timestamp: "15 March at 10:05", Permalink: "#post-815", PermalinkLabel: "15 March at 10:05, link to message 815",
		Content:   commentFixture(`<p>Rich <a href="/notes">notes</a></p><ul><li>A step</li></ul><pre><code>x &lt; y</code></pre><img src="/repair.jpg" alt="Repaired stair">`),
		Edited:    commentFixture(`Edited <time datetime="2026-03-15T10:10:00-07:00">10:10</time>`),
		Footer:    commentFixture(`<ul class="comment-actions" aria-label="Comment actions"><li><button type="button">Reply</button></li><li><a href="/history">History</a></li></ul>`),
		BodyClass: "text-base mt-075", TimestampClass: "text-muted text-xs", EditedClass: "text-xs",
	}
	output := testutil.RenderToString(t, CommentArticleComponent(props))
	article := testutil.FindElement(testutil.ParseFragment(t, output), "article")
	for key, value := range map[string]string{"id": "post-815", "class": "comment comment--layout caller", "data-message": "815"} {
		if got, _ := testutil.AttrVal(article, key); got != value {
			t.Fatalf("%s: %q != %q", key, got, value)
		}
	}
	content := commentByClass(t, output, "comment-content")
	var order []string
	for n := content.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode {
			v, _ := testutil.AttrVal(n, "class")
			order = append(order, n.Data+":"+v)
		}
	}
	if got := strings.Join(order, "|"); got != "header:comment-meta|div:comment-body prose text-base mt-075|p:comment-edited text-xs|footer:comment-footer" {
		t.Fatal(got)
	}
	stamp := commentByClass(t, output, "comment-timestamp")
	if got, _ := testutil.AttrVal(stamp, "datetime"); got != props.DateTime || stamp.FirstChild.Data != props.Timestamp {
		t.Fatal("caller time lost")
	}
	if stamp.Parent.Data != "a" || stamp.Parent.Parent.Data != "header" {
		t.Fatal("permalink is not a sibling of the author")
	}
	avatar := commentByClass(t, output, "comment-avatar")
	if label, _ := testutil.AttrVal(avatar, "aria-label"); avatar.Data != "a" || label != "Sam's profile" {
		t.Fatal("avatar link has no accessible name")
	}
	// Inspect token nesting before HTML parser repair could hide invalid links.
	z := html.NewTokenizer(strings.NewReader(output))
	interactive := ""
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		tok := z.Token()
		if tok.Data != "a" && tok.Data != "button" {
			continue
		}
		if tt == html.StartTagToken {
			if interactive != "" {
				t.Fatal("nested interactive elements")
			}
			interactive = tok.Data
		} else if tt == html.EndTagToken {
			interactive = ""
		}
	}
	for _, unwanted := range []string{"aria-live", `role="log"`, "cluster", "flex-1", "stack"} {
		if strings.Contains(output, unwanted) {
			t.Fatalf("unexpected automatic behavior/layout: %s", unwanted)
		}
	}
}

func TestCommentArticleOptionalContentAndChildren(t *testing.T) {
	child := commentFixture("<p>Children body</p>")
	render := func(props CommentArticleProps) string {
		return testutil.RenderToString(t, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return CommentArticleComponent(props).Render(templ.WithChildren(ctx, child), w)
		}))
	}
	output := render(CommentArticleProps{})
	if !strings.Contains(output, "Children body") {
		t.Fatal("children lost")
	}
	for _, unwanted := range []string{"comment-avatar", "comment-author", "comment-timestamp", "comment-permalink", "comment-edited", "comment-footer", `id=""`} {
		if strings.Contains(output, unwanted) {
			t.Fatalf("empty optional output: %s", unwanted)
		}
	}
	output = render(CommentArticleProps{Content: commentFixture("<p>Explicit body</p>")})
	if !strings.Contains(output, "Explicit body") || strings.Contains(output, "Children body") {
		t.Fatal("Content must replace children")
	}
	for _, id := range []string{"comment-one", "comment-two"} {
		output := render(CommentArticleProps{ID: id, Avatar: commentFixture(`<img alt="" src="/avatar.jpg">`), AvatarHref: "/unnamed-profile", Timestamp: "March 15", DateTime: "2026-03-15"})
		avatar := commentByClass(t, output, "comment-avatar")
		if hidden, _ := testutil.AttrVal(avatar, "aria-hidden"); avatar.Data != "span" || hidden != "true" {
			t.Fatal("unnamed avatar must remain decorative and unlinked")
		}
		stamp := commentByClass(t, output, "comment-timestamp")
		if stamp.Parent.Data != "header" || !strings.Contains(output, `id="`+id+`"`) {
			t.Fatal("unlinked semantic time or caller ID lost")
		}
	}
}
