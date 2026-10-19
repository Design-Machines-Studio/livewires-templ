# Card footer and summary cards

Adds a footer region for any card and a summary card for items in a grid. This
package emits HTML only. Use Live Wires **v2.26.0**, commit
`d4c197b29b1675bf6cada14d051b7a8e1425fe2f`, or a later release with the same
`.card .footer` and `.card--summary` CSS. Earlier versions render the markup
without the footer rule, equal heights or line limits.

The authority is the producer's
[reference markup](https://github.com/Design-Machines-Studio/livewires/blob/d4c197b29b1675bf6cada14d051b7a8e1425fe2f/public/reference/components/cards.html)
and [manual](https://github.com/Design-Machines-Studio/livewires/blob/d4c197b29b1675bf6cada14d051b7a8e1425fe2f/public/manual/components/cards.html)
(Live Wires PR #39). The summary card section of the reference is retained in
`component/testdata/livewires/card-summary-reference.html` and compared with the
composed output in `component/card_summary_test.go`.

Existing cards (`Card`, `CardLink`, `CardWithMedia`, `CardWithHeader`) keep
their output. `CardWithHeader` still renders the header rule because a body
follows the header.

## Components

| Component | Emits | Props |
| --- | --- | --- |
| `CardFooterComponent` / `CardFooter(class)` | `footer.footer.box` with children | `Class`, `Attrs` |
| `CardFrameComponent` / `CardFrame(scheme)` | bare `div.card` with children (no `box`) | `CardProps`: `Scheme`, `Class`, `Attrs` |
| `CardHeaderComponent` | `div.header.box` with a heading | `Title`, `HeadingLevel`, `Class`, `Attrs` |
| `CardSummaryComponent` / `CardSummary(title, href, description, scheme, headerClass)` | `li.card.card--summary`: header, then children | see below |
| `CardPeopleComponent` / `CardPeople(people, extra, label, class)` | up to three decorative avatars, `+N`, hidden label | `People`, `Extra`, `Label`, `Class`, `Attrs` |
| `CardCountComponent` / `CardCount(value, label, class)` | icon, value, hidden label | `Value`, `Label`, `Icon`, `Class`, `Attrs` |
| `CardDateComponent` / `CardDate(datetime, label, class)` | `<time datetime>` with optional prefix, hidden label and icon | `Datetime`, `Label`, `Prefix`, `HiddenLabel`, `Icon`, `Class`, `Attrs` |
| `CardAuthorComponent` / `CardAuthor(name, src, class)` | decorative avatar, then name block with children | `Name`, `Src`, `HiddenPrefix`, `Class`, `Attrs` |

`CardSummaryProps`: `Title`, `HeadingLevel` (default `h3`), `Href`,
`Description`, `Badge`, `BadgeVariant` (default `blue`), `Scheme`, `HeaderClass`,
`HeadingRowClass`, `Class`, `Attrs`.

## Footer on any card

`CardFooter` works inside any `.card`. A footer is ruled off from what comes
before it. When it follows a header directly, its rule replaces the header's.
Build that card with `CardFrame` and `CardHeaderComponent`:

```templ
@component.CardFrameComponent(component.CardProps{Scheme: "subtle", Class: "rounded"}) {
    @component.CardHeaderComponent(component.CardHeaderProps{Title: "Header directly followed by footer"})
    @component.CardFooter("cluster cluster-between items-center text-sm") {
        <span>Footer</span>
    }
}
```

Put a body between them (`<div class="box">...</div>`) and the header keeps its
rule. `CardWithHeader` always renders a body, so it never loses the rule.

## Summary card

```templ
<ul class="grid list-none">
    @component.CardSummaryComponent(component.CardSummaryProps{
        Title: "Board room", Href: "/spaces/board",
        Description: "Board and officer preparation between meetings.",
        Scheme: "subtle", Class: "rounded", HeaderClass: "stack",
    }) {
        @component.CardFooter("cluster cluster-between items-center text-sm") {
            @component.CardPeople(people, "+2", "5 members", "cluster cluster-compact items-center")
            @component.CardCountComponent(component.CardCountProps{
                Value: "4", Label: "conversations, 1 unread",
                Icon: chatsIcon(), Class: "inline-flex items-center gap-025 font-bold",
            })
        }
    }
</ul>
```

- The card is an `li`. Render it inside a `ul` or `ol`, such as `ul.grid.list-none`.
- Children render after the header. Pass a `CardFooter` as the child. A card with
  no footer child renders without a footer.
- The heading link is the only interactive element in the heading. With no
  `Href` the heading renders without a link. Links use `templ.URL`, so unsafe
  schemes such as `javascript:` are neutralized.
- `Badge` renders `badge badge--<BadgeVariant>` after the heading, for example
  `New` or `Updated`. The heading row (`HeadingRowClass`) renders when `Badge` or
  `HeadingRowClass` is set, so cards without a badge can keep the same structure
  as cards with one.

## Layout classes are caller-owned

The components emit no `cluster`, `stack` or other layout primitives. They add
only `box` on the footer (like `header box`), `block truncate` on the author
name, `text-xs text-muted` on dates and `min-w-0` on a heading beside a badge.
The caller passes the canonical classes from the producer markup:

| Where | Class to pass | Through |
| --- | --- | --- |
| Header | `stack` | `HeaderClass` |
| Heading and badge row | `cluster cluster-between cluster-nowrap items-start` | `HeadingRowClass` |
| Footer, people and count | `cluster cluster-between items-center text-sm` | `CardFooter(class)` |
| Footer, author | `cluster cluster-compact cluster-nowrap items-center text-sm` | `CardFooter(class)` |
| Footer, author and members | `cluster cluster-between cluster-nowrap items-center text-sm` | `CardFooter(class)` |
| People cluster | `cluster cluster-compact items-center` | `CardPeople(..., class)` |
| Author and date wrapper | `cluster cluster-compact cluster-nowrap items-center min-w-0` | a caller `<div>` around `CardAuthorComponent` |
| Count | `inline-flex items-center gap-025` (add `font-bold` for unread, `flex-shrink-0` beside an author) | `CardCount(..., class)` |
| Clock date | `flex items-center gap-025` | `CardDateComponent.Class` |
| Card | `rounded` | `Class` |

`CardDateComponent` places `Class` before `text-xs text-muted`, matching the
producer's order.

## Line limits are CSS-owned

Live Wires stops headings after 3 lines and descriptions after 4 with an
ellipsis, through `--card-heading-lines` and `--card-description-lines`. The full
text stays in the DOM. There is no line-limit prop. To change a limit, set the
property on the card through `Attrs`:

```templ
@component.CardSummaryComponent(component.CardSummaryProps{
    Title: title, Href: href, Description: description,
    Attrs: templ.Attributes{"style": "--card-heading-lines: 2; --card-description-lines: 2;"},
}) { ... }
```

`Attrs` and `Class` follow the package rule: never populate them from user input.

## Footer parts

- **People.** `CardPeople` shows the first three `CardPerson{Name, Src}` as
  `avatar--sm` in an `aria-hidden` avatar group. Images use an empty `alt`.
  Without `Src` the avatar shows an `aria-hidden` initial and carries no accessible name.
  `Extra` (`+2`, `+ All`) is `aria-hidden`. `Label` ("5 members") is visually
  hidden and carries the count. The caller supplies both strings.
- **Count.** `CardCount` renders the value and a visually hidden label with a
  leading space, such as ` conversations, 1 unread`.
- **Date.** `CardDate` renders `<time datetime>` with a caller-supplied short
  label (`Oct 17`, `Yesterday`, `2 days ago`). Nothing is formatted or computed.
  `Prefix: "Edited"` gives the sentence form. `HiddenLabel: "Last updated"` with
  an `Icon` gives the clock form.
- **Author.** `CardAuthorComponent` renders a decorative avatar and a text block
  with a visually hidden `By` (`HiddenPrefix` changes it). The name sits in
  `block truncate`, so long names get an ellipsis while the DOM keeps the full
  name. Pass a `CardDate` as the child. The avatar and text are siblings, so a
  footer cluster lays them out directly, or wrap them in a caller cluster.
- **Icons.** The package has no icon set. `Icon` takes any decorative
  `templ.Component` (no focusable controls). The producer uses Phosphor Bold
  `chats`, `users-three` and `clock`. The gallery defines them locally.

## Producer caveat

Initial avatars blend into `--color-subtle` backgrounds and can disappear on
`scheme-subtle` summary cards. Prefer image avatars there. The hidden label still
names the people.

## Gallery

`/examples/cards` in `examples/gallery` shows existing cards, footers on any card
(including header directly followed by footer), summary grids for every
variation, the line limits and the initials caveat. All names, counts, dates and
copy on that page are placeholder data. The avatar image is the producer's
guide image served from the locked asset set.
