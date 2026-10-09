package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

type listState struct {
	Order    []string `json:"order"`
	Revision uint64   `json:"revision"`
}

var initialOrders = map[string][]string{
	"reading":      {"essay-01", "essay-02", "essay-03"},
	"requirements": {"orientation", "agreements", "governance"},
}

var errStale = errors.New("The list changed. Reload before trying again.")
var errInvalidMove = errors.New("The item or insertion point is invalid for this list.")

// One process owns this demo file. The mutex includes validation and the entire
// durable commit; request delays are deliberately outside this operation.
type store struct {
	mu          sync.Mutex
	path        string
	lists       map[string]listState
	unavailable bool
	// Narrow I/O seams exercise real rename and directory-sync failures.
	rename    func(string, string) error
	syncDir   func(string) error
	writeTemp func(*os.File, []byte) error
}

func openStore(path string) (*store, error) {
	s := &store{path: path, rename: os.Rename, syncDir: syncDirectory, writeTemp: writeTemporary}
	if err := s.reload(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		s.lists = make(map[string]listState)
		for key, order := range initialOrders {
			s.lists[key] = listState{Order: slices.Clone(order), Revision: 1}
		}
		if err := s.commit(s.lists); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func syncDirectory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = f.Sync()
	return errors.Join(err, f.Close())
}

func writeTemporary(f *os.File, b []byte) error {
	if _, err := f.Write(b); err != nil {
		return err
	}
	return f.Sync()
}

func (s *store) reload() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > 16384 {
		return errors.New("invalid order file size or type")
	}
	var state map[string]listState
	d := json.NewDecoder(io.LimitReader(f, 16385))
	d.DisallowUnknownFields()
	if err := d.Decode(&state); err != nil {
		return fmt.Errorf("invalid order file: %w", err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("invalid trailing order data")
	}
	if len(state) != len(initialOrders) {
		return errors.New("invalid list keys")
	}
	for key, expected := range initialOrders {
		v, ok := state[key]
		if !ok || v.Revision == 0 || len(v.Order) != len(expected) {
			return errors.New("invalid list state")
		}
		seen := map[string]bool{}
		for _, id := range v.Order {
			if !slices.Contains(expected, id) || seen[id] {
				return errors.New("invalid order membership")
			}
			seen[id] = true
		}
	}
	s.lists = state
	return nil
}

func (s *store) snapshot(key string) listState {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.lists[key]
	v.Order = slices.Clone(v.Order)
	return v
}

func (s *store) move(key, item, before string, revision uint64) (listState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.lists[key]
	v.Order = slices.Clone(v.Order)
	if !ok {
		return listState{}, errInvalidMove
	}
	if s.unavailable {
		return listState{}, errors.New("Saving is unavailable. Reload and restart the demo.")
	}
	if revision != v.Revision {
		return v, errStale
	}
	if !slices.Contains(v.Order, item) || before == item || (before != "" && !slices.Contains(v.Order, before)) {
		return v, errInvalidMove
	}
	order := slices.Clone(v.Order)
	index := slices.Index(order, item)
	order = slices.Delete(order, index, index+1)
	if before == "" {
		order = append(order, item)
	} else {
		order = slices.Insert(order, slices.Index(order, before), item)
	}
	if slices.Equal(order, v.Order) {
		return v, nil
	}
	if v.Revision == ^uint64(0) {
		return v, errors.New("List revision exhausted.")
	}
	candidate := make(map[string]listState, len(s.lists))
	for key, state := range s.lists {
		candidate[key] = state
	}
	candidate[key] = listState{Order: order, Revision: v.Revision + 1}
	if err := s.commit(candidate); err != nil {
		if s.unavailable {
			return listState{}, err
		}
		current := s.lists[key]
		current.Order = slices.Clone(current.Order)
		return current, err
	}
	v = s.lists[key]
	v.Order = slices.Clone(v.Order)
	return v, nil
}

func (s *store) commit(candidate map[string]listState) error {
	b, err := json.Marshal(candidate)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".orders-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	writeErr := s.writeTemp(f, append(b, '\n'))
	err = errors.Join(writeErr, f.Close())
	if err != nil {
		return err
	}
	if err := s.rename(f.Name(), s.path); err != nil {
		return err
	}
	if err := s.syncDir(filepath.Dir(s.path)); err != nil {
		// Rename has already selected the authoritative file. Never publish the
		// previous in-memory state as though that rename had been undone.
		if reloadErr := s.reload(); reloadErr != nil {
			s.unavailable = true
			return errors.New("Commit outcome is uncertain and the authoritative file could not be read. Restart required.")
		}
		return errors.New("Commit outcome is uncertain after directory sync. Reload before trying again.")
	}
	s.lists = candidate
	return nil
}
