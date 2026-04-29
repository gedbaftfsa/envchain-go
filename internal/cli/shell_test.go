package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/store"
)

func newShellStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedShell(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	es := newEnvSet(t, map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "s3cr3t value",
	})
	if err := st.Save(project, pass, es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdShellExportsVars(t *testing.T) {
	st := newShellStore(t)
	seedShell(t, st, "myapp", "pw")

	var buf bytes.Buffer
	if err := CmdShell(st, "pw", "myapp", &buf); err != nil {
		t.Fatalf("CmdShell: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "export DB_HOST=") {
		t.Errorf("expected DB_HOST export, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_PASS=") {
		t.Errorf("expected DB_PASS export, got:\n%s", out)
	}
	// value with space must be quoted
	if !strings.Contains(out, "'s3cr3t value'") {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestCmdShellSorted(t *testing.T) {
	st := newShellStore(t)
	seedShell(t, st, "myapp", "pw")

	var buf bytes.Buffer
	if err := CmdShell(st, "pw", "myapp", &buf); err != nil {
		t.Fatalf("CmdShell: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output not sorted: %q before %q", lines[i-1], lines[i])
		}
	}
}

func TestCmdUnshellUnsetsVars(t *testing.T) {
	st := newShellStore(t)
	seedShell(t, st, "myapp", "pw")

	var buf bytes.Buffer
	if err := CmdUnshell(st, "pw", "myapp", &buf); err != nil {
		t.Fatalf("CmdUnshell: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "unset DB_HOST") {
		t.Errorf("expected unset DB_HOST, got:\n%s", out)
	}
	if !strings.Contains(out, "unset DB_PASS") {
		t.Errorf("expected unset DB_PASS, got:\n%s", out)
	}
	if strings.Contains(out, "export") {
		t.Errorf("unexpected export in unshell output:\n%s", out)
	}
}

func TestCmdShellEmptyProject(t *testing.T) {
	st := newShellStore(t)
	var buf bytes.Buffer
	err := CmdShell(st, "pw", "", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdShellWrongPassphrase(t *testing.T) {
	st := newShellStore(t)
	seedShell(t, st, "myapp", "correct")

	var buf bytes.Buffer
	err := CmdShell(st, "wrong", "myapp", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
