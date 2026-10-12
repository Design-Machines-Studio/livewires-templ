package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoutesRenderCurrentComponentsAndNativeMoves(t *testing.T) {
	app := routes(sourceIdentity{Commit: strings.Repeat("a", 40), AssetRevision: "producer"}, map[string][]byte{"/assets/main.css": []byte("verified css")}, assetLock{})
	cases := map[string]string{"/": "Component gallery", "/examples/relative-time": "lw-relative-time", "/examples/copy": "data-copy-source", "/examples/loading": "data-loading-region", "/examples/toast": "data-toast-region", "/examples/inline-edit": "data-inline-edit", "/examples/sortable": "lw-sortable-list", "/examples/sortable?order=notes,paper,index&itemId=notes&before=index": "value=\"paper,notes,index\"", "/examples/loading?run=1&scenario=failure": "Simulated failure", "/assets/main.css": "verified css"}
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
