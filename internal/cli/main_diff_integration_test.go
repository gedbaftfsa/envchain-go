package cli

import (
	"testing"
)

func TestMainDiffDispatch(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	pass := "testpass"

	if err := CmdSet(st, pass, "src", "A=1"); err != nil {
		t.Fatal(err)
	}
	if err := CmdSet(st, pass, "dst", "B=2"); err != nil {
		t.Fatal(err)
	}

	// Simulate Main dispatching "diff" sub-command via CmdDiff directly.
	err := CmdDiff(st, pass, "src", "dst", nil)
	if err != nil {
		t.Fatalf("CmdDiff returned error: %v", err)
	}
}

func TestMainDiffWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	pass := "correct"

	if err := CmdSet(st, pass, "proj", "X=1"); err != nil {
		t.Fatal(err)
	}

	err := CmdDiff(st, "wrong", "proj", "proj", nil)
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
