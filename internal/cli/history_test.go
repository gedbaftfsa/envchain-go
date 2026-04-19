package cli

import (
	"bytes"
	"strings"
	"testing"
)

func seedHistory(t *testing.T) (*testStore, string) {
	t.Helper()
	st, pass := newTempStoreSnap(t)
	project := "histproj"
	// create two snapshots
	if err := CmdSnapshot(st, project, "v1", pass); err != nil {
		t.Fatalf("snapshot v1: %v", err)
	}
	if err := CmdSnapshot(st, project, "v2", pass); err != nil {
		t.Fatalf("snapshot v2: %v", err)
	}
	return st, pass
}

func TestCmdHistoryLists(t *testing.T) {
	st, pass := seedHistory(t)
	var buf bytes.Buffer
	if err := CmdHistory(st, "histproj", pass, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "v1") {
		t.Errorf("expected v1 in output, got: %s", out)
	}
	if !strings.Contains(out, "v2") {
		t.Errorf("expected v2 in output, got: %s", out)
	}
}

func TestCmdHistoryEmpty(t *testing.T) {
	st, pass := newTempStoreSnap(t)
	var buf bytes.Buffer
	if err := CmdHistory(st, "nosnaps", pass, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no history") {
		t.Errorf("expected 'no history' message, got: %s", buf.String())
	}
}

func TestCmdHistoryEmptyProject(t *testing.T) {
	st, pass := newTempStoreSnap(t)
	_ = pass
	var buf bytes.Buffer
	err := CmdHistory(st, "", "", &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}
