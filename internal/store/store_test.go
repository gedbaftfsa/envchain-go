package store_test

import (
	"os"
	"testing"

	"github.com/yourusername/envchain-go/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func TestSaveAndLoad(t *testing.T) {
	s := tempStore(t)
	set := &store.EnvSet{
		Name: "myproject",
		Vars: map[string]string{"API_KEY": "secret123", "DEBUG": "true"},
	}
	pass := "hunter2"
	if err := s.Save(set, pass); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := s.Load("myproject", pass)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != set.Name {
		t.Errorf("name mismatch: got %q want %q", loaded.Name, set.Name)
	}
	for k, v := range set.Vars {
		if loaded.Vars[k] != v {
			t.Errorf("var %q: got %q want %q", k, loaded.Vars[k], v)
		}
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	s := tempStore(t)
	set := &store.EnvSet{Name: "proj", Vars: map[string]string{"X": "1"}}
	if err := s.Save(set, "correct"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	_, err := s.Load("proj", "wrong")
	if err != store.ErrBadPassphrase {
		t.Errorf("expected ErrBadPassphrase, got %v", err)
	}
}

func TestLoadNotFound(t *testing.T) {
	s := tempStore(t)
	_, err := s.Load("nonexistent", "pass")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	s := tempStore(t)
	set := &store.EnvSet{Name: "todel", Vars: map[string]string{}}
	if err := s.Save(set, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := s.Delete("todel"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Load("todel", "pass"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("ghost"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
