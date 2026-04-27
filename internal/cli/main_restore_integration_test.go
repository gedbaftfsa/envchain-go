package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainRestoreDispatch(t *testing.T) {
	dir := t.TempDir()

	// Write a backup file
	backup := filepath.Join(dir, "proj.env")
	if err := os.WriteFile(backup, []byte("KEY=value\n"), 0600); err != nil {
		t.Fatalf("write backup: %v", err)
	}

	passReader := strings.NewReader("testpass\n")
	out, err := runMain([]string{"restore", "proj", backup}, dir, passReader)
	if err != nil {
		t.Fatalf("Main restore: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "restored project") {
		t.Errorf("expected restored message, got: %s", out)
	}
}

func TestMainRestoreWrongArgs(t *testing.T) {
	dir := t.TempDir()
	passReader := strings.NewReader("testpass\n")
	out, err := runMain([]string{"restore", "proj"}, dir, passReader)
	if err == nil {
		t.Fatalf("expected error for missing file arg, got output: %s", out)
	}
}
