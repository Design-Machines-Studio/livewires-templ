package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

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
