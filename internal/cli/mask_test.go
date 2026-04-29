package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
	"github.com/envchain/envchain-go/internal/store"
)

func newMaskStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedMask(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	es := env.NewSet()
	es.Put("API_KEY", "supersecretvalue")
	es.Put("DB_PASS", "short")
	es.Put("EMPTY_VAR", "")
	if err := st.Save(project, pass, es); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestCmdMaskOutput(t *testing.T) {
	st := newMaskStore(t)
	seedMask(t, st, "proj", "pass")

	var buf bytes.Buffer
	if err := CmdMask(st, "proj", "pass", &buf); err != nil {
		t.Fatalf("CmdMask: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "supersecretvalue") {
		t.Error("output must not contain plaintext secret")
	}
	if !strings.Contains(out, "API_KEY=") {
		t.Error("expected API_KEY in output")
	}
	// value capped at 8 asterisks
	if !strings.Contains(out, "API_KEY=********") {
		t.Errorf("expected 8 asterisks for long value, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS=") {
		t.Error("expected DB_PASS in output")
	}
	// short value: 5 chars → 5 asterisks
	if !strings.Contains(out, "DB_PASS=*****") {
		t.Errorf("expected 5 asterisks for 'short', got: %s", out)
	}
	if !strings.Contains(out, "EMPTY_VAR=(empty)") {
		t.Errorf("expected (empty) marker for empty value, got: %s", out)
	}
}

func TestCmdMaskEmptyProject(t *testing.T) {
	st := newMaskStore(t)
	var buf bytes.Buffer
	if err := CmdMask(st, "", "pass", &buf); err == nil {
		t.Error("expected error for empty project name")
	}
}

func TestCmdMaskNotFound(t *testing.T) {
	st := newMaskStore(t)
	var buf bytes.Buffer
	if err := CmdMask(st, "ghost", "pass", &buf); err == nil {
		t.Error("expected error for missing project")
	}
}

func TestCmdMaskWrongPassphrase(t *testing.T) {
	st := newMaskStore(t)
	seedMask(t, st, "proj", "correct")
	var buf bytes.Buffer
	if err := CmdMask(st, "proj", "wrong", &buf); err == nil {
		t.Error("expected error for wrong passphrase")
	}
}

func TestMaskValue(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", "(empty)"},
		{"abc", "***"},
		{"12345678", "********"},
		{"verylongpassword", "********"},
	}
	for _, c := range cases {
		got := maskValue(c.input)
		if got != c.want {
			t.Errorf("maskValue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
