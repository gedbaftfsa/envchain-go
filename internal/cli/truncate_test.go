package cli

import (
	"bytes"
	"testing"
)

func seedTruncate(t *testing.T) (*testStore, string) {
	t.Helper()
	st, pass := newTempStoreSnap(t)
	project := "myapp"

	// Create 4 snapshots.
	for _, tag := range []string{"v1", "v2", "v3", "v4"} {
		if err := CmdSnapshot(st, project, pass, tag, &bytes.Buffer{}); err != nil {
			t.Fatalf("seed snapshot %s: %v", tag, err)
		}
	}
	return st, pass
}

func TestCmdTruncateKeepTwo(t *testing.T) {
	st, pass := seedTruncate(t)
	_ = pass
	var buf bytes.Buffer
	if err := CmdTruncate(st, "myapp", pass, 2, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snaps, err := st.ListSnapshots("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snaps))
	}
}

func TestCmdTruncateNothingToDo(t *testing.T) {
	st, pass := seedTruncate(t)
	var buf bytes.Buffer
	if err := CmdTruncate(st, "myapp", pass, 10, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("nothing to truncate")) {
		t.Fatalf("expected nothing-to-truncate message, got: %s", buf.String())
	}
}

func TestCmdTruncateEmptyProject(t *testing.T) {
	st, pass := newTempStoreSnap(t)
	var buf bytes.Buffer
	err := CmdTruncate(st, "", pass, 2, &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdTruncateKeepZero(t *testing.T) {
	st, pass := newTempStoreSnap(t)
	var buf bytes.Buffer
	err := CmdTruncate(st, "myapp", pass, 0, &buf)
	if err == nil {
		t.Fatal("expected error for keep=0")
	}
}
