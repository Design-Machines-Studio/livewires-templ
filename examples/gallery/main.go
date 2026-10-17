// The gallery is local development tooling. Demo outcomes never persist.
package main

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/a-h/templ"
)

//go:embed assets.lock.json demo.js
var authored embed.FS

// Compiled by run.py from the immutable build snapshot, never runtime labels.
var buildReceipt string

type sourceIdentity struct {
	Checkout      string `json:"checkout"`
	Branch        string `json:"branch"`
	Commit        string `json:"commit"`
	Dirty         bool   `json:"dirty"`
	Fingerprint   string `json:"fingerprint"`
	AssetRevision string `json:"assetRevision"`
	LockDigest    string `json:"lockDigest"`
}
type asset struct {
	File   string `json:"file"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
	Size   int    `json:"size"`
	Serve  string `json:"serve"`
}
type assetLock struct {
	ProducerCommit string  `json:"producerCommit"`
	Assets         []asset `json:"assets"`
}

var examples = []string{"relative-time", "copy", "loading", "toast", "sortable", "inline-edit", "combobox", "save-status", "page-header", "comments", "loaders"}
var itemNames = map[string]string{"notes": "Field notes", "paper": "On paper", "index": "An index"}

func digest(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }

func loadAssets(directory string) (map[string][]byte, assetLock, error) {
	b, _ := authored.ReadFile("assets.lock.json")
	var lock assetLock
	if err := json.Unmarshal(b, &lock); err != nil {
		return nil, lock, err
	}
	files := map[string][]byte{}
	for _, a := range lock.Assets {
		p := filepath.Join(directory, a.File)
		real, err := filepath.EvalSymlinks(p)
		if err != nil || real != p {
			return nil, lock, fmt.Errorf("missing or unsafe asset: %s", a.File)
		}
		data, err := os.ReadFile(p)
		if err != nil || len(data) != a.Size || digest(data) != a.SHA256 {
			return nil, lock, fmt.Errorf("asset checksum mismatch: %s", a.File)
		}
		if a.Serve != "" {
			files[a.Serve] = data
		}
	}
	return files, lock, nil
}

func validOrder(value string) []string {
	order := strings.Split(value, ",")
	seen := map[string]bool{}
	for _, id := range order {
		if itemNames[id] == "" || seen[id] {
			return []string{"notes", "paper", "index"}
		}
		seen[id] = true
	}
	if len(order) != 3 {
		return []string{"notes", "paper", "index"}
	}
	return order
}

func moveOrder(order []string, item, before string) []string {
	if itemNames[item] == "" || (before != "" && itemNames[before] == "") || item == before {
		return order
	}
	result := []string{}
	for _, id := range order {
		if id != item {
			result = append(result, id)
		}
	}
	for i, id := range result {
		if id == before {
			return append(append(result[:i:i], item), result[i:]...)
		}
	}
	return append(result, item)
}

func outcome(scenario string) string {
	switch scenario {
	case "failure":
		return "Simulated failure. Nothing was saved."
	case "unknown":
		return "Simulated unknown result. Reload to recover; nothing is persisted."
	default:
		return "Simulated success. Nothing was saved."
	}
}

func routes(source sourceIdentity, files map[string][]byte, lock assetLock) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"source": source, "assets": lock.Assets})
	})
	for path, data := range files {
		m.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(path, ".js") {
				w.Header().Set("Content-Type", "text/javascript")
			} else if strings.HasSuffix(path, ".css") {
				w.Header().Set("Content-Type", "text/css")
			}
			w.Header().Set("Cache-Control", "no-cache")
			w.Write(data)
		})
	}
	m.HandleFunc("GET /demo.js", func(w http.ResponseWriter, r *http.Request) {
		b, _ := authored.ReadFile("demo.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(b)
	})
	m.HandleFunc("GET /demo/result", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("scenario") == "delayed" {
			select {
			case <-time.After(800 * time.Millisecond):
			case <-r.Context().Done():
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"scenario": r.URL.Query().Get("scenario"), "message": outcome(r.URL.Query().Get("scenario"))})
	})
	m.HandleFunc("GET /examples/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		found := false
		for _, n := range examples {
			if n == name {
				found = true
			}
		}
		if !found {
			http.NotFound(w, r)
			return
		}
		order := validOrder(r.URL.Query().Get("order"))
		message := ""
		variant := "info"
		if (name == "loading" || name == "toast") && r.URL.Query().Has("run") {
			message = outcome(r.URL.Query().Get("scenario"))
			if r.URL.Query().Get("scenario") == "failure" {
				variant = "error"
			}
		}
		if name == "inline-edit" && (r.URL.Query().Has("title-a") || r.URL.Query().Has("title-b")) {
			message = "Native example submission. Fictional text is carried in the URL; nothing was saved."
		}
		if name == "combobox" && (r.URL.Query().Has("interface_timezone") || r.URL.Query().Has("required_timezone")) {
			message = fmt.Sprintf("Native example submission. interface_timezone=%q; required_timezone=%q. These values came from this URL; nothing was saved.", r.URL.Query().Get("interface_timezone"), r.URL.Query().Get("required_timezone"))
		}
		if name == "sortable" && r.URL.Query().Has("itemId") {
			order = moveOrder(order, r.URL.Query().Get("itemId"), r.URL.Query().Get("before"))
			message = "Native demo move. Order is carried in this URL; reload the base URL to reset."
		}
		templ.Handler(page(name, source, order, message, variant)).ServeHTTP(w, r)
	})
	m.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(page("", source, nil, "", "info")).ServeHTTP(w, r)
	})
	return m
}

func main() {
	assets := flag.String("assets", "/assets", "verified development cache")
	listen := flag.String("listen", "0.0.0.0:8080", "container listen address")
	flag.Parse()
	files, lock, err := loadAssets(*assets)
	if err != nil {
		log.Fatal(err)
	}
	raw, err := base64.RawURLEncoding.DecodeString(buildReceipt)
	var source sourceIdentity
	b, _ := authored.ReadFile("assets.lock.json")
	if err != nil || json.Unmarshal(raw, &source) != nil || len(source.Commit) != 40 || len(source.Fingerprint) != 64 || source.LockDigest != digest(b) || source.AssetRevision != lock.ProducerCommit {
		log.Fatal("compiled source receipt required; use gallery run.sh")
	}
	server := &http.Server{Addr: *listen, Handler: routes(source, files, lock), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second}
	log.Fatal(server.ListenAndServe())
}
