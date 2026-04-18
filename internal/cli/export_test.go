package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCmdExport(t *testing.T) {
	st := newTempStore(t)
	const proj, pass = "myapp", "secret"

	if err := CmdSet(proj, pass, st, "FOO=bar"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(proj, pass, st, "BAZ=hello world"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdExport(proj, pass, st, &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "export BAZ='hello world'\n") {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
	if !strings.Contains(out, "export FOO=bar\n") {
		t.Errorf("expected unquoted value, got:\n%s", out)
	}
}

func TestCmdExportNotFound(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdExport("ghost", "x", st, &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdImport(t *testing.T) {
	st := newTempStore(t)
	const proj, pass = "imp", "pass"

	input := []byte("# comment\nFOO=bar\nexport BAZ=qux\nEMPTY=\n")
	if err := CmdImport(proj, pass, st, nil, input); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdExport(proj, pass, st, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"export FOO=bar", "export BAZ=qux", "export EMPTY="} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestShellQuote(t *testing.T) {
	cases := []struct{ in, want string }{
		{"simple", "simple"},
		{"hello world", "'hello world'"},
		{"it's", "'it'\\''s'"},
		{"$VAR", "'$VAR'"},
	}
	for _, c := range cases {
		if got := shellQuote(c.in); got != c.want {
			t.Errorf("shellQuote(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
