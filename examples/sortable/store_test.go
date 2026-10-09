package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func testStore(t *testing.T) *store {
	t.Helper()
	s, err := openStore(filepath.Join(t.TempDir(), "orders.json"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}
func TestPersistenceDuplicatesAndIndependentLists(t *testing.T) {
	s := testStore(t)
	state, err := s.move("reading", "essay-02", "essay-01", 1)
	if err != nil || state.Revision != 2 || state.Order[0] != "essay-02" {
		t.Fatalf("move: %+v %v", state, err)
	}
	b, _ := os.ReadFile(s.path)
	if _, err := s.move("reading", "essay-02", "essay-01", 1); !errors.Is(err, errStale) {
		t.Fatal("duplicate not stale", err)
	}
	after, _ := os.ReadFile(s.path)
	if string(after) != string(b) {
		t.Fatal("replay wrote file")
	}
	if _, err := s.move("reading", "essay-02", "essay-01", 2); err != nil {
		t.Fatal("current no-op rejected", err)
	}
	after, _ = os.ReadFile(s.path)
	if string(after) != string(b) {
		t.Fatal("no-op wrote file")
	}
	if _, err := s.move("requirements", "governance", "orientation", 1); err != nil {
		t.Fatal(err)
	}
	reopened, err := openStore(s.path)
	if err != nil {
		t.Fatal(err)
	}
	for _, list := range []string{"reading", "requirements"} {
		if !reflect.DeepEqual(reopened.snapshot(list), s.snapshot(list)) {
			t.Fatal("lost persisted list", list)
		}
	}
	if reopened.snapshot("reading").Revision != 2 || reopened.snapshot("requirements").Revision != 2 {
		t.Fatal("independent revisions lost")
	}
}
func TestConcurrentSameRevision(t *testing.T) {
	s := testStore(t)
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, item := range []string{"essay-02", "essay-03"} {
		wg.Add(1)
		go func(item string) { defer wg.Done(); _, err := s.move("reading", item, "essay-01", 1); results <- err }(item)
	}
	wg.Wait()
	close(results)
	accepted, stale := 0, 0
	for err := range results {
		if err == nil {
			accepted++
		} else if errors.Is(err, errStale) {
			stale++
		} else {
			t.Fatal(err)
		}
	}
	if accepted != 1 || stale != 1 || s.snapshot("reading").Revision != 2 || s.snapshot("requirements").Revision != 1 {
		t.Fatal("serialization or list independence failed")
	}
}
func TestInvalidMovesDoNotWrite(t *testing.T) {
	s := testStore(t)
	before, _ := os.ReadFile(s.path)
	for _, move := range [][3]string{{"unknown", "essay-01", ""}, {"reading", "governance", ""}, {"reading", "essay-01", "orientation"}, {"reading", "essay-01", "essay-01"}, {"reading", "missing", ""}} {
		if _, err := s.move(move[0], move[1], move[2], 1); err == nil {
			t.Fatal("invalid move accepted", move)
		}
	}
	after, _ := os.ReadFile(s.path)
	if string(before) != string(after) {
		t.Fatal("invalid requests wrote state")
	}
}
func TestCorruptStoreIsNeverReseeded(t *testing.T) {
	for _, bad := range []string{"{", "{}", `{"reading":{"order":["essay-01","essay-01","essay-03"],"revision":1},"requirements":{"order":["orientation","agreements","governance"],"revision":1}}`, `{} {}`} {
		path := filepath.Join(t.TempDir(), "orders.json")
		if err := os.WriteFile(path, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := openStore(path); err == nil {
			t.Fatal("corruption accepted")
		}
		b, _ := os.ReadFile(path)
		if string(b) != bad {
			t.Fatal("corrupt file reseeded")
		}
	}
}

func TestOversizedStoreIsRejected(t *testing.T) {
	s := testStore(t)
	valid, _ := os.ReadFile(s.path)
	// A limited decoder must not hide corrupt data beyond its read boundary.
	bad := append(valid, []byte(strings.Repeat(" ", 16385)+"corrupt trailing bytes")...)
	if err := os.WriteFile(s.path, bad, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := openStore(s.path); err == nil {
		t.Fatal("oversized/corrupt tail accepted")
	}
}
func corruptAfterRename(s *store) {
	s.rename = func(from, to string) error {
		if err := os.Rename(from, to); err != nil {
			return err
		}
		return os.WriteFile(to, []byte("{}"), 0600)
	}
	s.syncDir = func(dir string) error { return syncDirectory(filepath.Join(dir, "missing")) }
}

func TestUncertainReloadFailureDoesNotPublishOldCache(t *testing.T) {
	s := testStore(t)
	corruptAfterRename(s)
	state, err := s.move("reading", "essay-02", "essay-01", 1)
	if err == nil || !strings.Contains(err.Error(), "uncertain") || state.Revision != 0 || !s.unavailable {
		t.Fatal("unreadable authoritative state misreported", state, err)
	}
	if state, err := s.move("reading", "essay-03", "essay-01", 1); err == nil || state.Revision != 0 {
		t.Fatal("unavailable store returned old state")
	}
	for list := range initialOrders {
		if state := s.snapshot(list); state.Revision != 0 || len(state.Order) != 0 {
			t.Fatal("unavailable snapshot returned old cache", list, state)
		}
	}
}

func TestDurableFailures(t *testing.T) {
	for _, stage := range []string{"create", "write", "rename", "after-rename"} {
		t.Run(stage, func(t *testing.T) {
			s := testStore(t)
			path := s.path
			before, _ := os.ReadFile(path)
			switch stage {
			case "create":
				s.path = filepath.Join(path, "missing", "orders.json")
			case "write":
				s.writeTemp = func(f *os.File, b []byte) error { f.Close(); return writeTemporary(f, b) }
			case "rename":
				s.rename = func(from, to string) error { return os.Rename(from, filepath.Join(to, "missing")) }
			case "after-rename":
				s.syncDir = func(dir string) error { return syncDirectory(filepath.Join(dir, "missing")) }
			}
			state, err := s.move("reading", "essay-02", "essay-01", 1)
			if err == nil {
				t.Fatal("failure accepted")
			}
			after, _ := os.ReadFile(path)
			if stage == "after-rename" {
				if !strings.Contains(err.Error(), "uncertain") || state.Revision != 2 || state.Order[0] != "essay-02" || string(before) == string(after) {
					t.Fatal("post-rename state misreported", state, err)
				}
				reopened, err := openStore(path)
				if err != nil || !reflect.DeepEqual(reopened.snapshot("reading"), state) {
					t.Fatal("authoritative reload mismatch")
				}
			} else if string(before) != string(after) || state.Revision != 1 {
				t.Fatal("pre-rename failure changed state", state)
			}
			files, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".orders-*"))
			if len(files) != 0 {
				t.Fatal("temporary files leaked", files)
			}
		})
	}
}
