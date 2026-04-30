package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newScatterStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedScatter(t *testing.T, st *store.Store) {
	t.Helper()
	es := mustEnvSet(t, "DB_HOST=localhost", "DB_PORT=5432", "API_KEY=secret")
	if err := st.Save("myapp", es, "pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdScatterWritesFiles(t *testing.T) {
	st := newScatterStore(t)
	seedScatter(t, st)
	out := t.TempDir()

	var buf bytes.Buffer
	if err := CmdScatter(st, "myapp", out, "pass", &buf); err != nil {
		t.Fatalf("CmdScatter: %v", err)
	}

	for _, name := range []string{"API_KEY", "DB_HOST", "DB_PORT"} {
		path := filepath.Join(out, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected file %s: %v", path, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("file %s is empty", path)
		}
	}

	if !strings.Contains(buf.String(), "wrote ") {
		t.Errorf("expected progress output, got: %q", buf.String())
	}
}

func TestCmdScatterEmptyProject(t *testing.T) {
	st := newScatterStore(t)
	var buf bytes.Buffer
	err := CmdScatter(st, "", t.TempDir(), "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdScatterNotFound(t *testing.T) {
	st := newScatterStore(t)
	var buf bytes.Buffer
	err := CmdScatter(st, "ghost", t.TempDir(), "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdScatterWrongPassphrase(t *testing.T) {
	st := newScatterStore(t)
	seedScatter(t, st)
	var buf bytes.Buffer
	err := CmdScatter(st, "myapp", t.TempDir(), "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdScatterCreatesDir(t *testing.T) {
	st := newScatterStore(t)
	seedScatter(t, st)
	newDir := filepath.Join(t.TempDir(), "nested", "secrets")

	var buf bytes.Buffer
	if err := CmdScatter(st, "myapp", newDir, "pass", &buf); err != nil {
		t.Fatalf("CmdScatter: %v", err)
	}

	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
