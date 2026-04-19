package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain-go/internal/store"
)

func newDoctorStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestCmdDoctorClean(t *testing.T) {
	st := newDoctorStore(t)
	var buf bytes.Buffer
	if err := CmdDoctor(st, &buf); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[ OK ] store directory:") {
		t.Errorf("expected OK for store directory, got:\n%s", out)
	}
	if !strings.Contains(out, "[ OK ] store directory permissions:") {
		t.Errorf("expected OK for permissions, got:\n%s", out)
	}
}

func TestCmdDoctorBadPermissions(t *testing.T) {
	st := newDoctorStore(t)
	dir := st.Dir()
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o700)
	var buf bytes.Buffer
	err := CmdDoctor(st, &buf)
	if err == nil {
		t.Fatal("expected error for bad permissions")
	}
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("expected WARN in output, got:\n%s", buf.String())
	}
}

func TestCmdDoctorMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	st, _ := store.New(dir)
	// Remove dir after creation
	os.RemoveAll(dir)
	var buf bytes.Buffer
	err := CmdDoctor(st, &buf)
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
	if !strings.Contains(buf.String(), "FAIL") {
		t.Errorf("expected FAIL in output, got:\n%s", buf.String())
	}
}
