package cli

import (
	"os"
	"testing"

	"github.com/user/envchain-go/internal/env"
)

func writeFakeEditor(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "fake-editor-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	// Overwrite the file passed as $1 with fixed content
	f.WriteString("#!/bin/sh\ncat > \"$1\" <<'EOF'\n" + content + "EOF\n")
	f.Close()
	os.Chmod(f.Name(), 0755)
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestCmdEditCreatesAndUpdates(t *testing.T) {
	st := newTempStore(t)
	const proj = "editproj"
	const pass = "secret"

	// Seed initial value
	set := env.NewSet()
	set.Put("FOO", "bar")
	if err := st.Save(proj, pass, set); err != nil {
		t.Fatal(err)
	}

	// Use fake editor that writes new content
	editor := writeFakeEditor(t, "FOO=newval\nBAZ=qux\n")
	t.Setenv("EDITOR", editor)

	if err := CmdEdit(st, proj, pass); err != nil {
		t.Fatalf("CmdEdit: %v", err)
	}

	loaded, err := st.Load(proj, pass)
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := loaded.Get("FOO"); v != "newval" {
		t.Errorf("FOO = %q, want newval", v)
	}
	if v, _ := loaded.Get("BAZ"); v != "qux" {
		t.Errorf("BAZ = %q, want qux", v)
	}
}

func TestCmdEditNewProject(t *testing.T) {
	st := newTempStore(t)
	editor := writeFakeEditor(t, "NEW=value\n")
	t.Setenv("EDITOR", editor)

	if err := CmdEdit(st, "brand-new", "pass"); err != nil {
		t.Fatalf("CmdEdit new project: %v", err)
	}

	loaded, err := st.Load("brand-new", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := loaded.Get("NEW"); v != "value" {
		t.Errorf("NEW = %q, want value", v)
	}
}
