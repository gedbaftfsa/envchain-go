package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func newNamespaceStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func seedNamespace(t *testing.T, st *store.Store) {
	t.Helper()
	es := env.NewSet()
	for _, entry := range []string{
		"DB_HOST=localhost",
		"DB_PASS=secret",
		"DB_PORT=5432",
		"AWS_KEY=AKID",
		"AWS_SECRET=abc123",
		"PLAIN=noprefix",
	} {
		k, v, _ := env.ParseEntry(entry)
		es.Put(k, v)
	}
	if err := st.Save("myapp", es, "pass"); err != nil {
		t.Fatal(err)
	}
}

func TestCmdNamespaceList(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	var buf bytes.Buffer
	if err := CmdNamespace(st, "myapp", "pass", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "AWS") {
		t.Errorf("expected AWS namespace, got: %s", out)
	}
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB namespace, got: %s", out)
	}
	if !strings.Contains(out, "3 key(s)") {
		t.Errorf("expected DB to have 3 keys, got: %s", out)
	}
}

func TestCmdNamespaceKeys(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	var buf bytes.Buffer
	if err := CmdNamespaceKeys(st, "myapp", "DB", "pass", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, k := range []string{"DB_HOST", "DB_PASS", "DB_PORT"} {
		if !strings.Contains(out, k) {
			t.Errorf("expected key %s in output, got: %s", k, out)
		}
	}
}

func TestCmdNamespaceKeysNoMatch(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	var buf bytes.Buffer
	if err := CmdNamespaceKeys(st, "myapp", "GCP", "pass", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no keys found") {
		t.Errorf("expected no-match message")
	}
}

func TestCmdNamespaceEmptyProject(t *testing.T) {
	st := newNamespaceStore(t)
	var buf bytes.Buffer
	err := CmdNamespace(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdNamespaceWrongPassphrase(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)
	var buf bytes.Buffer
	err := CmdNamespace(st, "myapp", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
