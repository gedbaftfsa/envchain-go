package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestMainGCDispatch(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)

	// seed a project and snapshot it
	es := mustEnvSet(t, "TOKEN=abc")
	if err := st.Save("proj", "secret", es); err != nil {
		t.Fatalf("Save: %v", err)
	}
	var snap bytes.Buffer
	if err := CmdSnapshot(st, "proj", "secret", "v1", &snap); err != nil {
		t.Fatalf("CmdSnapshot: %v", err)
	}

	// delete project to orphan the snapshot
	if err := st.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// run gc via Main dispatch
	var out bytes.Buffer
	code := Main([]string{"envchain", "gc"}, st, noPassphrase, &out)
	if code != 0 {
		t.Fatalf("Main returned %d; output: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "removed 1 orphaned") {
		t.Errorf("expected removal message, got: %q", out.String())
	}
}

func TestMainGCHelp(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)
	var out bytes.Buffer
	code := Main([]string{"envchain", "help", "gc"}, st, noPassphrase, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "gc") {
		t.Errorf("help output missing 'gc': %q", out.String())
	}
}
