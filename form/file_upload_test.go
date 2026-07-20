package form

import (
	"maps"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

const fileUploadGoldenDefault = `<label class="file-upload"><input type="file" name="documents" accept=".pdf" disabled multiple> <span class="file-upload-button">Choose documents</span> <span class="file-upload-text">PDF documents only</span></label>`

func TestFileUploadRenders(t *testing.T) {
	html := testutil.RenderToString(t, FileUploadSimple("resume", "Upload resume"))
	if !strings.Contains(html, `class="file-upload"`) {
		t.Error("expected file-upload class")
	}
	if !strings.Contains(html, `type="file"`) {
		t.Error("expected file input type")
	}
	if !strings.Contains(html, `name="resume"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "Upload resume") {
		t.Error("expected button text")
	}
	if !strings.Contains(html, "No file selected") {
		t.Error("expected default status text")
	}
}

func TestFileUploadDefaultTexts(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{Name: "doc"}))
	if !strings.Contains(html, "Choose file") {
		t.Error("expected default button text")
	}
	if !strings.Contains(html, "No file selected") {
		t.Error("expected default status text")
	}
}

func TestFileUploadVariantAccent(t *testing.T) {
	html := testutil.RenderToString(t, FileUploadAccent("doc", "Upload"))
	if !strings.Contains(html, "file-upload--accent") {
		t.Error("expected file-upload--accent class")
	}
}

func TestFileUploadVariantCompact(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:    "import",
		Variant: "compact",
	}))
	if !strings.Contains(html, "file-upload--compact") {
		t.Error("expected file-upload--compact class")
	}
}

func TestFileUploadZone(t *testing.T) {
	html := testutil.RenderToString(t, FileUploadZone("docs", "Browse files", "or drag and drop"))
	if !strings.Contains(html, "file-upload--zone") {
		t.Error("expected file-upload--zone class")
	}
	if !strings.Contains(html, "Browse files") {
		t.Error("expected button text")
	}
	if !strings.Contains(html, "or drag and drop") {
		t.Error("expected hint text")
	}
}

func TestFileUploadZoneAccent(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:    "images",
		Variant: "accent",
		Zone:    true,
	}))
	if !strings.Contains(html, "file-upload--zone") {
		t.Error("expected file-upload--zone class")
	}
	if !strings.Contains(html, "file-upload--accent") {
		t.Error("expected file-upload--accent class")
	}
}

func TestFileUploadAccept(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:   "resume",
		Accept: ".pdf,.doc",
	}))
	if !strings.Contains(html, `accept=".pdf,.doc"`) {
		t.Error("expected accept attribute")
	}
}

func TestFileUploadMultiple(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:     "files",
		Multiple: true,
	}))
	if !strings.Contains(html, "multiple") {
		t.Error("expected multiple attribute")
	}
}

func TestFileUploadDisabled(t *testing.T) {
	html := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:     "locked",
		Disabled: true,
	}))
	if !strings.Contains(html, "disabled") {
		t.Error("expected disabled attribute")
	}
}

func TestFileUploadInputAttrsRenderOnInput(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:       "documents",
		Attrs:      templ.Attributes{"data-testid": "wrapper"},
		InputAttrs: inputAttrs,
	}))
	nodes := testutil.ParseFragment(t, output)
	label := testutil.FindElement(nodes, "label")
	input := testutil.FindElement(nodes, "input")

	if got, ok := testutil.AttrVal(label, "data-testid"); !ok || got != "wrapper" {
		t.Errorf("label data-testid = %q, %v; want %q, true", got, ok, "wrapper")
	}
	if _, ok := testutil.AttrVal(input, "data-testid"); ok {
		t.Error("input unexpectedly contains wrapper data-testid")
	}
	for key, want := range inputAttrs {
		key = strings.ToLower(key)
		if got, ok := testutil.AttrVal(input, key); !ok || got != want {
			t.Errorf("input %s = %q, %v; want %q, true", key, got, ok, want)
		}
		if _, ok := testutil.AttrVal(label, key); ok {
			t.Errorf("label unexpectedly contains input attribute %s", key)
		}
	}
}

func TestFileUploadInputAttrsExactlyOnce(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:       "documents",
		Attrs:      templ.Attributes{"data-testid": "wrapper"},
		InputAttrs: inputAttrs,
	}))
	inputTag := testutil.RawTag(t, output, "input")
	labelTag := testutil.RawTag(t, output, "label")
	for key := range inputAttrs {
		if got := testutil.CountAttr(inputTag, key); got != 1 {
			t.Errorf("input %s occurrence count = %d; want 1", key, got)
		}
		if got := testutil.CountAttr(labelTag, key); got != 0 {
			t.Errorf("label %s occurrence count = %d; want 0", key, got)
		}
	}
}

func TestFileUploadInputAttrsCollisionComponentWins(t *testing.T) {
	// FileUpload emits accept, disabled, and multiple itself; a caller that
	// re-supplies them through InputAttrs must lose to the component's value
	// (HTML first-occurrence), and the raw tag must show both occurrences.
	output := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name:     "documents",
		Accept:   ".pdf",
		Disabled: true,
		Multiple: true,
		InputAttrs: templ.Attributes{
			"accept":   "bogus",
			"disabled": "false",
		},
	}))
	input := testutil.FindElement(testutil.ParseFragment(t, output), "input")
	if got, ok := testutil.AttrVal(input, "accept"); !ok || got != ".pdf" {
		t.Errorf("parsed accept = %q, %v; want component value %q", got, ok, ".pdf")
	}
	inputTag := testutil.RawTag(t, output, "input")
	for _, key := range []string{"accept", "disabled"} {
		if got := testutil.CountAttr(inputTag, key); got != 2 {
			t.Errorf("raw input %s occurrence count = %d; want 2 (component + caller)", key, got)
		}
	}
}

func TestFileUploadInputAttrsDoNotMutateCallerMaps(t *testing.T) {
	attrs := templ.Attributes{"data-testid": "wrapper"}
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	attrsBefore := maps.Clone(attrs)
	inputAttrsBefore := maps.Clone(inputAttrs)

	testutil.RenderToString(t, FileUpload(FileUploadProps{Name: "documents", Attrs: attrs, InputAttrs: inputAttrs}))

	if !maps.Equal(attrs, attrsBefore) {
		t.Errorf("Attrs mutated: got %#v, want %#v", attrs, attrsBefore)
	}
	if !maps.Equal(inputAttrs, inputAttrsBefore) {
		t.Errorf("InputAttrs mutated: got %#v, want %#v", inputAttrs, inputAttrsBefore)
	}
}

func TestFileUploadEmptyInputAttrsUnchangedOutput(t *testing.T) {
	omitted := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name: "documents", ButtonText: "Choose documents", Text: "PDF documents only",
		Accept: ".pdf", Multiple: true, Disabled: true,
	}))
	nilAttrs := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name: "documents", ButtonText: "Choose documents", Text: "PDF documents only",
		Accept: ".pdf", Multiple: true, Disabled: true, InputAttrs: nil,
	}))
	emptyAttrs := testutil.RenderToString(t, FileUpload(FileUploadProps{
		Name: "documents", ButtonText: "Choose documents", Text: "PDF documents only",
		Accept: ".pdf", Multiple: true, Disabled: true, InputAttrs: templ.Attributes{},
	}))

	for name, got := range map[string]string{"omitted": omitted, "nil": nilAttrs, "empty": emptyAttrs} {
		if got != fileUploadGoldenDefault {
			t.Errorf("%s InputAttrs output changed:\ngot  %q\nwant %q", name, got, fileUploadGoldenDefault)
		}
	}
	if omitted != nilAttrs || omitted != emptyAttrs {
		t.Error("omitted, nil, and empty InputAttrs outputs differ")
	}
}

func TestFileUploadInputAttrsNoAriaMachinery(t *testing.T) {
	inputAttrs := templ.Attributes{
		"data-bind":                      "autosaveEnabled",
		"data-on:change":                 "@post('/settings/autosave')",
		"data-indicator:autosavePending": "",
		"data-attr:disabled":             "$autosavePending ? true : false",
		"aria-label":                     "Enable autosave",
	}
	output := testutil.RenderToString(t, FileUpload(FileUploadProps{Name: "documents", InputAttrs: inputAttrs}))
	nodes := testutil.ParseFragment(t, output)
	input := testutil.FindElement(nodes, "input")
	if _, ok := testutil.AttrVal(input, "aria-invalid"); ok {
		t.Error("file input unexpectedly has aria-invalid")
	}
	if _, ok := testutil.AttrVal(input, "aria-describedby"); ok {
		t.Error("file input unexpectedly has aria-describedby")
	}
	button := testutil.FindElementByClass(nodes, "file-upload-button")
	if button == nil || testutil.NodeText(button) != "Choose file" {
		t.Errorf("default button text missing, got %q", testutil.NodeText(button))
	}
	status := testutil.FindElementByClass(nodes, "file-upload-text")
	if status == nil || testutil.NodeText(status) != "No file selected" {
		t.Errorf("default status text missing, got %q", testutil.NodeText(status))
	}
	if testutil.FindElementByClass(nodes, "hint") != nil || testutil.FindElementByClass(nodes, "error") != nil {
		t.Error("file upload unexpectedly rendered hint or error element")
	}
}
