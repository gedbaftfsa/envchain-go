package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCmdStatsEmpty(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	if err := CmdStats(st, "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No projects found") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestCmdStatsMultiple(t *testing.T) {
	st := newTempStore(t)
	passphrase := "secret"

	// seed two projects
	if err := CmdSet(st, "alpha", passphrase, "FOO=bar"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, "alpha", passphrase, "BAZ=qux"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, "beta", passphrase, "X=1"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdStats(st, passphrase, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output")
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected beta in output")
	}
	if !strings.Contains(out, "2 project(s), 3 variable(s)") {
		t.Errorf("expected totals, got: %s", out)
	}
}

func TestCmdStatsWrongPassphrase(t *testing.T) {
	st := newTempStore(t)
	if err := CmdSet(st, "proj", "correct", "K=V"); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	err := CmdStats(st, "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
