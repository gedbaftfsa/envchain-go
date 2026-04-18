package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newTempStoreLP(t *testing.T) *store.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-lp-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return store.New(filepath.Join(dir, "store"))
}

func seedLP(t *testing.T, st *store.Store, project, passphrase string) {
	t.Helper()
	set := env.NewSet()
	_ = set.Put("KEY=val")
	if err := st.Save(project, passphrase, set); err != nil {
		t.Fatal(err)
	}
}

func TestCmdListProjectsEmpty(t *testing.T) {
	st := newTempStoreLP(t)
	var buf bytes.Buffer
	if err := CmdListProjects(st, nil, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no projects") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestCmdListProjectsMultiple(t *testing.T) {
	st := newTempStoreLP(t)
	seedLP(t, st, "alpha", "pass1")
	seedLP(t, st, "beta", "pass2")
	seedLP(t, st, "gamma", "pass3")

	var buf bytes.Buffer
	if err := CmdListProjects(st, nil, &buf); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 projects, got %d: %v", len(lines), lines)
	}
	if lines[0] != "alpha" || lines[1] != "beta" || lines[2] != "gamma" {
		t.Errorf("unexpected order or names: %v", lines)
	}
}
