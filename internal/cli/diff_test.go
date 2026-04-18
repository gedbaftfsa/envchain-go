package cli

import (
	"bytes"
	"strings"
	"testing"
)

func seedDiff(t *testing.T, st interface{ Save(string, interface{}) error }) {}

func TestCmdDiffDistinctKeys(t *testing.T) {
	st := newTempStore(t)
	pass := "secret"

	// project alpha: FOO, BAR
	if err := CmdSet(st, pass, "alpha", "FOO=1"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, pass, "alpha", "BAR=2"); err != nil {
		t.Fatal(err)
	}

	// project beta: BAR, BAZ
	if err := CmdSet(st, pass, "beta", "BAR=2"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, pass, "beta", "BAZ=3"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdDiff(st, pass, "alpha", "beta", &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "< FOO") {
		t.Errorf("expected '< FOO' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "> BAZ") {
		t.Errorf("expected '> BAZ' in output, got:\n%s", out)
	}
	if strings.Contains(out, "BAR") {
		t.Errorf("BAR values are equal, should not appear; got:\n%s", out)
	}
}

func TestCmdDiffChangedValue(t *testing.T) {
	st := newTempStore(t)
	pass := "secret"

	if err := CmdSet(st, pass, "p1", "KEY=hello"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, pass, "p2", "KEY=world"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdDiff(st, pass, "p1", "p2", &buf); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "~ KEY") {
		t.Errorf("expected '~ KEY', got: %s", buf.String())
	}
}

func TestCmdDiffProjectNotFound(t *testing.T) {
	st := newTempStore(t)
	pass := "secret"

	if err := CmdSet(st, pass, "exists", "X=1"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := CmdDiff(st, pass, "exists", "missing", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdDiffIdentical(t *testing.T) {
	st := newTempStore(t)
	pass := "secret"

	if err := CmdSet(st, pass, "a", "K=v"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, pass, "b", "K=v"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdDiff(st, pass, "a", "b", &buf); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty diff, got: %s", buf.String())
	}
}
