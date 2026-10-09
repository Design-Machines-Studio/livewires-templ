// The sortable example is a local, isolated consumer, not an authentication SDK.
package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//go:embed assets.lock.json bridge.js
var authored embed.FS

// run.sh embeds the receipt into the compiled binary, rather than trusting an
// environment label at runtime. Assets and bridge bytes are embedded/verified.
var buildReceipt string

type assetIdentity struct {
	File   string `json:"file"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
	Serve  string `json:"serve"`
}
type assetLock struct {
	ProducerCommit string          `json:"producerCommit"`
	Assets         []assetIdentity `json:"assets"`
}
type sourceReceipt struct {
	Application string `json:"application"`
	Checkout    string `json:"physicalCheckout"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	Fingerprint string `json:"sourceFingerprint"`
	Dirty       bool   `json:"dirty"`
	LockDigest  string `json:"lockDigest"`
}

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }

// Freeze only verified bytes into the allowlist. No FileServer, run-directory
// traversal, store file, session data, or mutable on-demand file reads.
func loadAssets(dir string) (map[string][]byte, assetLock, error) {
	b, _ := authored.ReadFile("assets.lock.json")
	var lock assetLock
	if err := json.Unmarshal(b, &lock); err != nil {
		return nil, lock, err
	}
	files := make(map[string][]byte)
	for _, asset := range lock.Assets {
		if filepath.IsAbs(asset.File) || filepath.Clean(asset.File) != asset.File || strings.HasPrefix(asset.File, "..") {
			return nil, lock, errors.New("unsafe asset path")
		}
		path := filepath.Join(dir, asset.File)
		real, err := filepath.EvalSymlinks(path)
		if err != nil || real != path {
			return nil, lock, fmt.Errorf("missing or unsafe locked asset: %s", asset.File)
		}
		bytes, err := os.ReadFile(path)
		if err != nil || digest(bytes) != asset.SHA256 {
			return nil, lock, fmt.Errorf("unverified locked asset: %s", asset.File)
		}
		if asset.Serve != "" {
			files[asset.Serve] = bytes
		}
	}
	return files, lock, nil
}

func main() {
	assets := flag.String("assets", "", "absolute verified asset directory")
	orders := flag.String("store", "", "absolute task-owned order file")
	listen := flag.String("listen", "0.0.0.0:8080", "isolated container HTTP address")
	flag.Parse()
	if !filepath.IsAbs(*assets) || !filepath.IsAbs(*orders) {
		log.Fatal("explicit absolute assets/store paths required")
	}
	files, lock, err := loadAssets(*assets)
	if err != nil {
		log.Fatal(err)
	}
	var receipt sourceReceipt
	raw, err := base64.RawURLEncoding.DecodeString(buildReceipt)
	if err != nil || json.Unmarshal(raw, &receipt) != nil || receipt.Application != "livewires-templ-sortable" || len(receipt.Fingerprint) != 64 || len(receipt.Commit) != 40 {
		log.Fatal("compiled build receipt required; use run.sh start")
	}
	locked, _ := authored.ReadFile("assets.lock.json")
	if receipt.LockDigest != digest(locked) {
		log.Fatal("compiled receipt/asset lock mismatch")
	}
	s, err := openStore(*orders)
	if err != nil {
		log.Fatal(err)
	}
	a := newApp(s, files, lock, receipt)
	server := &http.Server{Addr: *listen, Handler: a.routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 30 * time.Second}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		// Close listeners first, then drain synchronous handlers/durable writes.
		drain, stop := context.WithTimeout(context.Background(), 30*time.Second)
		defer stop()
		if err := server.Shutdown(drain); err != nil {
			log.Fatal(err)
		}
	}
}
