package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
	"github.com/nicholasgasior/envchain-go/internal/store"
)

func newPeekStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("newPeekStore: %v", err)
	}
	return st
}

func seedPeek(t *testing.T, st *store.Store, passphrase string) {
	t.Helper()
	es := env.NewSet()
	es.Put("DB_HOST", "localhost")
	es.Put("DB_PORT", "5432")
	es.Put("API_KEY", "secret")
	if err := st.Save("myproject", passphrase, es); err != nil {
		t.Fatalf("seedPeek: %v", err)
	}
}

func TestCmdPeekFound(t *testing.T) {
	st := newPeekStore(t)
	seedPeek(t, st, "pass")

	var buf bytes.Buffer
	if err := CmdPeek(st, "pass", "myproject", "DB_HOST", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "localhost" {
		t.Errorf("expected %q, got %q", "localhost", got)
	}
}

func TestCmdPeekMissingKey(t *testing.T) {
	st := newPeekStore(t)
	seedPeek(t, st, "pass")

	var buf bytes.Buffer
	err := CmdPeek(st, "pass", "myproject", "MISSING", &buf)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention 'not found', got: %v", err)
	}
}

func TestCmdPeekWrongPassphrase(t *testing.T) {
	st := newPeekStore(t)
	seedPeek(t, st, "pass")

	var buf bytes.Buffer
	if err := CmdPeek(st, "wrong", "myproject", "DB_HOST", &buf); err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdPeekEmptyArgs(t *testing.T) {
	st := newPeekStore(t)
	var buf bytes.Buffer
	if err := CmdPeek(st, "pass", "", "KEY", &buf); err == nil {
		t.Fatal("expected error for empty project")
	}
	if err := CmdPeek(st, "pass", "proj", "", &buf); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestCmdPeekAll(t *testing.T) {
	st := newPeekStore(t)
	seedPeek(t, st, "pass")

	var buf bytes.Buffer
	if err := CmdPeekAll(st, "pass", "myproject", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, kv := range []string{"API_KEY=secret", "DB_HOST=localhost", "DB_PORT=5432"} {
		if !strings.Contains(out, kv) {
			t.Errorf("expected output to contain %q\ngot:\n%s", kv, out)
		}
	}
}

func TestCmdPeekAllEmpty(t *testing.T) {
	st := newPeekStore(t)
	es := env.NewSet()
	if err := st.Save("empty", "pass", es); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var buf bytes.Buffer
	if err := CmdPeekAll(st, "pass", "empty", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "(empty)") {
		t.Errorf("expected '(empty)' output, got: %s", buf.String())
	}
}
