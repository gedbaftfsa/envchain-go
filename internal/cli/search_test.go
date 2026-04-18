package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/envchain-go/internal/store"
)

func seedSearch(t *testing.T, st *store.Store) {
	t.Helper()
	seedProjectInStore(t, st, "alpha", "pass", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	seedProjectInStore(t, st, "beta", "pass", map[string]string{"DB_NAME": "mydb", "API_KEY": "secret"})
	seedProjectInStore(t, st, "gamma", "pass", map[string]string{"REDIS_URL": "redis://localhost"})
}

func seedProjectInStore(t *testing.T, st *store.Store, name, pass string, vars map[string]string) {
	t.Helper()
	set, _ := st.Load(name, pass)
	if set == nil {
		set, _ = newEnvSet()
	}
	for k, v := range vars {
		_ = set.Put(k, v)
	}
	if err := st.Save(name, set, pass); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdSearchFound(t *testing.T) {
	st := newTempStore(t)
	seedSearch(t, st)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := CmdSearch(st, "pass", "DB")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected alpha/DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "beta") || !strings.Contains(out, "DB_NAME") {
		t.Errorf("expected beta/DB_NAME in output, got: %s", out)
	}
}

func TestCmdSearchNotFound(t *testing.T) {
	st := newTempStore(t)
	seedSearch(t, st)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := CmdSearch(st, "pass", "NONEXISTENT")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "no keys") {
		t.Errorf("expected 'no keys' message")
	}
}

func TestCmdSearchEmptyQuery(t *testing.T) {
	st := newTempStore(t)
	if err := CmdSearch(st, "pass", ""); err == nil {
		t.Fatal("expected error for empty query")
	}
}
