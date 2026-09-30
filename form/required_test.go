package form

import (
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
	"strings"
	"testing"
)

func TestRequiredIndicators(t *testing.T) {
	for _, required := range []bool{false, true} {
		for _, err := range []string{"", "Invalid"} {
			for _, component := range []templ.Component{Field(FieldProps{Name: "name", Label: "Long field label", Required: required, Error: err}), Select(SelectProps{Name: "choice", Label: "Choice", Required: required, Error: err}), Textarea(TextareaProps{Name: "text", Label: "Text", Required: required, Error: err})} {
				html := testutil.RenderToString(t, component)
				if strings.Contains(html, `>Required</span>`) != required {
					t.Fatalf("marker: %s", html)
				}
				if required && (!strings.Contains(html, `aria-hidden="true"`) || !strings.Contains(html, `aria-required="true"`)) {
					t.Fatal(html)
				}
				if err != "" && !strings.Contains(html, `aria-invalid="true"`) {
					t.Fatal(html)
				}
			}
		}
	}
}
