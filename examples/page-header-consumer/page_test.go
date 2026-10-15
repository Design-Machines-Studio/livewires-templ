package consumer

import (
	"context"
	"strings"
	"testing"
)

func TestExternalPackageImportsAndRendersPublicPageHeader(t *testing.T) {
	var output strings.Builder
	if err := PageHeaderPage().Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<header class="section section-snug">`,
		`<h1>External fixture page <span id="fixture-member-count">(12)</span></h1>`,
		`<p class="lead">Rendered through the published component package API.</p>`,
		`<a href="/members" class="button">View members</a>`,
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("external render missing %q in %s", want, output.String())
		}
	}
	if !strings.Contains(output.String(), `View members</a></div><div class="cluster cluster-between items-baseline">`) {
		t.Fatalf("external footer must follow the action layout: %s", output.String())
	}
}
