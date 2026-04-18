package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func seedLint(t *testing.T, st *store.Store) {
	t.Helper()
	set := env.NewSet()
	set.Put("DB_HOST", "localhost")
	set.Put("DB_PASS", "secret")
	if err := st.Save("clean", "pass", set); err != nil {
		t.Fatal(err)
	}

	dirty := env.NewSet()
	dirty.Put("DB_HOST", "")
	dirty.Put("PATH", "/custom/bin")
	dirty.Put("API_KEY", "abc123")
	if err := st.Save("dirty", "pass", dirty); err != nil {
		t.Fatal(err)
	}
}

func TestCmdLintClean(t *testing.T) {
	st := newTempStore(t)
	seedLint(t, st)

	var buf bytes.Buffer
	if err := CmdLint(st, "clean", "pass", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK, got: %s", out)
	}
}

func TestCmdLintDirty(t *testing.T) {
	st := newTempStore(t)
	seedLint(t, st)

	var buf bytes.Buffer
	if err := CmdLint(st, "dirty", "pass", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN lines, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST warning, got: %s", out)
	}
	if !strings.Contains(out, "PATH") {
		t.Errorf("expected PATH shadow warning, got: %s", out)
	}
}

func TestCmdLintNotFound(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdLint(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdLintWrongPassphrase(t *testing.T) {
	st := newTempStore(t)
	seedLint(t, st)
	var buf bytes.Buffer
	err := CmdLint(st, "clean", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
