package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRoutesRenderCurrentComponentsAndNativeMoves(t *testing.T) {
	app := routes(sourceIdentity{Commit: strings.Repeat("a", 40), AssetRevision: "producer"}, map[string][]byte{"/assets/main.css": []byte("verified css")}, assetLock{})
	cases := map[string]string{"/": "Component gallery", "/examples/relative-time": "lw-relative-time", "/examples/copy": "data-copy-source", "/examples/loading": "data-loading-region", "/examples/toast": "data-toast-region", "/examples/inline-edit": "data-inline-edit", "/examples/sortable": "lw-sortable-list", "/examples/combobox": "data-combobox", "/examples/save-status": "proposal-save-status", "/examples/combobox?interface_timezone=&required_timezone=Europe%2FParis": "Native example submission", "/examples/sortable?order=notes,paper,index&itemId=notes&before=index": "value=\"paper,notes,index\"", "/examples/loading?run=1&scenario=failure": "Simulated failure", "/assets/main.css": "verified css"}
	for path, want := range cases {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
			if w.Code != 200 || !strings.Contains(w.Body.String(), want) {
				t.Fatalf("%d %s", w.Code, w.Body.String())
			}
		})
	}
	for _, path := range []string{"/examples/missing", "/assets/unknown.js"} {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != http.StatusNotFound {
			t.Fatalf("%s: %d", path, w.Code)
		}
	}
	if got := strings.Join(validOrder("notes,notes,index"), ","); got != "notes,paper,index" {
		t.Fatal(got)
	}
}

func TestEveryAssetMustMatchLock(t *testing.T) {
	var lock assetLock
	b, _ := authored.ReadFile("assets.lock.json")
	if err := json.Unmarshal(b, &lock); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for _, a := range lock.Assets {
		path := filepath.Join(dir, a.File)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte("wrong bytes"), 0644)
	}
	if _, _, err := loadAssets(dir); err == nil {
		t.Fatal("corrupt assets accepted")
	}
	if _, _, err := loadAssets(t.TempDir()); err == nil {
		t.Fatal("missing assets accepted")
	}
}

func TestNativeScenarioControlsDoNotInventComponentOutcomes(t *testing.T) {
	app := routes(sourceIdentity{}, nil, assetLock{})
	for _, name := range []string{"sortable", "inline-edit"} {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest("GET", "/examples/"+name+"?scenario=failure&run=1&order=paper,notes,index", nil))
		if strings.Contains(w.Body.String(), "Simulated failure.") || strings.Contains(w.Body.String(), `name="run"`) {
			t.Fatalf("%s invented a result or exposed an inert run control", name)
		}
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/examples/save-status", nil))
	if strings.Count(w.Body.String(), `role="status"`) != 2 || !strings.Contains(w.Body.String(), "Rules and collaborators still need saving") {
		t.Fatal("save-status example must retain independent status regions and partial-scope text")
	}
	w = httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/examples/toast?scenario=failure&run=1", nil))
	if !strings.Contains(w.Body.String(), "toast--error") {
		t.Fatal("failure scenario did not render the error component variant")
	}
}

func TestComboboxTimezoneSamplesAreValidLocations(t *testing.T) {
	for _, option := range comboboxTimezoneOptions() {
		if _, err := time.LoadLocation(option.Value); err != nil {
			t.Errorf("timezone %q is not a valid location: %v", option.Value, err)
		}
	}
}
