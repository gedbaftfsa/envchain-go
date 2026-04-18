package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain-go/internal/cli"
	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(filepath.Join(dir, "envchain"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func seedProject(t *testing.T, st *store.Store, project, pass string, vars map[string]string) {
	t.Helper()
	set := env.NewSet()
	for k, v := range vars {
		if err := set.Put(k, v); err != nil {
			t.Fatalf("set.Put: %v", err)
		}
	}
	if err := st.Save(project, pass, set); err != nil {
		t.Fatalf("st.Save: %v", err)
	}
}

func TestRunNoArgs(t *testing.T) {
	st := tempStore(t)
	err := cli.RunFallback(st, cli.RunOptions{
		Project: "myapp", Passphrase: "secret", Args: nil,
	}, os.Stdout, os.Stderr)
	if err == nil || !strings.Contains(err.Error(), "no command") {
		t.Fatalf("expected 'no command' error, got %v", err)
	}
}

func TestRunProjectNotFound(t *testing.T) {
	st := tempStore(t)
	err := cli.RunFallback(st, cli.RunOptions{
		Project: "ghost", Passphrase: "x", Args: []string{"env"},
	}, os.Stdout, os.Stderr)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestRunInjectsEnv(t *testing.T) {
	st := tempStore(t)
	seedProject(t, st, "myapp", "pass", map[string]string{"ENVCHAIN_TEST_VAR": "hello"})

	pr, pw, _ := os.Pipe()
	defer pr.Close()

	err := cli.RunFallback(st, cli.RunOptions{
		Project:   "myapp",
		Passphrase: "pass",
		Args:      []string{"sh", "-c", "echo $ENVCHAIN_TEST_VAR"},
		Overwrite: true,
	}, pw, os.Stderr)
	pw.Close()
	if err != nil {
		t.Fatalf("RunFallback: %v", err)
	}

	buf := make([]byte, 64)
	n, _ := pr.Read(buf)
	got := strings.TrimSpace(string(buf[:n]))
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}
