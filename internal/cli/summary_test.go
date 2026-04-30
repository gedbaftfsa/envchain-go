package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
	"github.com/envchain/envchain-go/internal/store"
)

func newSummaryStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedSummary(t *testing.T, st *store.Store, name, pass string, pairs [][2]string) {
	t.Helper()
	set := env.NewSet()
	for _, p := range pairs {
		if err := set.Put(p[0], p[1]); err != nil {
			t.Fatalf("set.Put: %v", err)
		}
	}
	if err := st.Save(name, set, pass); err != nil {
		t.Fatalf("store.Save: %v", err)
	}
}

func TestCmdSummaryEmpty(t *testing.T) {
	st := newSummaryStore(t)
	var buf bytes.Buffer
	if err := CmdSummary(st, "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no projects found") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestCmdSummaryMultiple(t *testing.T) {
	st := newSummaryStore(t)
	const pass = "s3cr3t"
	seedSummary(t, st, "alpha", pass, [][2]string{{"A", "1"}, {"B", "2"}})
	seedSummary(t, st, "beta", pass, [][2]string{{"X", "9"}})

	var buf bytes.Buffer
	if err := CmdSummary(st, pass, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("missing alpha in output: %q", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("missing beta in output: %q", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected key count 2 for alpha: %q", out)
	}
}

func TestCmdSummaryWrongPassphrase(t *testing.T) {
	st := newSummaryStore(t)
	seedSummary(t, st, "proj", "correct", [][2]string{{"K", "V"}})

	var buf bytes.Buffer
	err := CmdSummary(st, "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestCmdSummaryHeader(t *testing.T) {
	st := newSummaryStore(t)
	seedSummary(t, st, "myproject", "pw", [][2]string{{"FOO", "bar"}})

	var buf bytes.Buffer
	if err := CmdSummary(st, "pw", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, col := range []string{"PROJECT", "KEYS", "PINNED", "PROTECTED"} {
		if !strings.Contains(out, col) {
			t.Errorf("missing column header %q in: %q", col, out)
		}
	}
}
