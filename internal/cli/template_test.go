package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain-go/internal/store"
)

func newTemplateStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func seedTemplate(t *testing.T, st *store.Store) {
	t.Helper()
	es := mustEnvSet(t, "HOST=db.local", "PORT=5432", "USER=admin")
	if err := st.Save("myapp", "pass", es); err != nil {
		t.Fatal(err)
	}
}

func TestCmdTemplateBasic(t *testing.T) {
	st := newTemplateStore(t)
	seedTemplate(t, st)

	var buf bytes.Buffer
	err := CmdTemplate(st, "myapp", "pass", "host={HOST} port={PORT}", &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "host=db.local port=5432" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestCmdTemplateTripleBrace(t *testing.T) {
	st := newTemplateStore(t)
	seedTemplate(t, st)

	var buf bytes.Buffer
	err := CmdTemplate(st, "myapp", "pass", "user={{{USER}}}", &buf)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "admin") {
		t.Errorf("expected admin in output, got %q", buf.String())
	}
}

func TestCmdTemplateWrongPassphrase(t *testing.T) {
	st := newTemplateStore(t)
	seedTemplate(t, st)

	var buf bytes.Buffer
	err := CmdTemplate(st, "myapp", "wrong", "host={HOST}", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdTemplateEmptyProject(t *testing.T) {
	st := newTemplateStore(t)
	err := CmdTemplate(st, "", "pass", "x={X}", os.Stdout)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdTemplateFile(t *testing.T) {
	st := newTemplateStore(t)
	seedTemplate(t, st)

	path := filepath.Join(t.TempDir(), "tmpl.txt")
	if err := os.WriteFile(path, []byte("connect {USER}@{HOST}:{PORT}"), 0600); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := CmdTemplateFile(st, "myapp", "pass", path, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "connect admin@db.local:5432" {
		t.Errorf("unexpected output: %q", out)
	}
}
