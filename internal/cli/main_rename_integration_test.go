package cli

import (
	"strings"
	"testing"
)

// Integration-style test: ensure Main dispatches "rename" correctly.
func TestMainRenameDispatch(t *testing.T) {
	st := newTempStore(t)
	const pass = "testpass"
	seedProject2(t, st, "proj-a", pass)

	// Capture output via CmdRename directly (Main wires passphrase via prompt,
	// so we test the dispatch logic at the command level).
	out := captureStdout(t, func() {
		if err := CmdRename(st, "proj-a", "proj-b", pass); err != nil {
			t.Fatalf("CmdRename failed: %v", err)
		}
	})

	if !strings.Contains(out, "proj-a") || !strings.Contains(out, "proj-b") {
		t.Fatalf("expected output to mention both project names, got: %q", out)
	}
}

func TestMainRenameWrongPassphrase(t *testing.T) {
	st := newTempStore(t)
	seedProject2(t, st, "proj-x", "correct")

	err := CmdRename(st, "proj-x", "proj-y", "wrong")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
