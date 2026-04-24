package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envchain-go/internal/store"
)

func newEnvDiffStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedEnvDiff(t *testing.T, st *store.Store) {
	t.Helper()
	es := mustEnvSet(t, "STORED_KEY=stored_val", "SHARED_KEY=project_val")
	if err := st.Save("myproject", "pass", es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdEnvDiffMissing(t *testing.T) {
	st := newEnvDiffStore(t)
	seedEnvDiff(t, st)

	t.Setenv("SHARED_KEY", "project_val")
	// STORED_KEY intentionally absent

	var buf bytes.Buffer
	if err := CmdEnvDiff(st, "myproject", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- STORED_KEY") {
		t.Errorf("expected missing marker for STORED_KEY, got:\n%s", out)
	}
}

func TestCmdEnvDiffChanged(t *testing.T) {
	st := newEnvDiffStore(t)
	seedEnvDiff(t, st)

	t.Setenv("STORED_KEY", "stored_val")
	t.Setenv("SHARED_KEY", "live_val") // differs from "project_val"

	var buf bytes.Buffer
	if err := CmdEnvDiff(st, "myproject", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ SHARED_KEY") {
		t.Errorf("expected changed marker for SHARED_KEY, got:\n%s", out)
	}
}

func TestCmdEnvDiffNoDifferences(t *testing.T) {
	st := newEnvDiffStore(t)

	es := mustEnvSet(t, "ONLY_KEY=val")
	if err := st.Save("clean", "pass", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("ONLY_KEY", "val")

	var buf bytes.Buffer
	if err := CmdEnvDiff(st, "clean", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no differences") {
		t.Errorf("expected 'no differences', got: %s", buf.String())
	}
}

func TestCmdEnvDiffEmptyProject(t *testing.T) {
	st := newEnvDiffStore(t)
	var buf bytes.Buffer
	err := CmdEnvDiff(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdEnvDiffWrongPassphrase(t *testing.T) {
	st := newEnvDiffStore(t)
	seedEnvDiff(t, st)
	var buf bytes.Buffer
	err := CmdEnvDiff(st, "myproject", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
