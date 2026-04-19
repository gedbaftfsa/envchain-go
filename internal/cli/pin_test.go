package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/store"
)

func newPinStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestCmdPinAndList(t *testing.T) {
	st := newPinStore(t)
	var buf bytes.Buffer

	if err := CmdSet(st, "proj", "pass", "FOO=bar", &buf); err != nil {
		t.Fatalf("CmdSet: %v", err)
	}
	buf.Reset()

	if err := CmdPin(st, "proj", "pass", []string{"FOO", "BAR"}, &buf); err != nil {
		t.Fatalf("CmdPin: %v", err)
	}
	if !strings.Contains(buf.String(), "pinned 2") {
		t.Errorf("unexpected output: %s", buf.String())
	}

	buf.Reset()
	if err := CmdListPinned(st, "proj", "pass", &buf); err != nil {
		t.Fatalf("CmdListPinned: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected FOO and BAR in output, got: %s", out)
	}
}

func TestCmdUnpin(t *testing.T) {
	st := newPinStore(t)
	var buf bytes.Buffer

	_ = CmdSet(st, "proj", "pass", "FOO=bar", &buf)
	_ = CmdPin(st, "proj", "pass", []string{"FOO", "BAR"}, &buf)
	buf.Reset()

	if err := CmdUnpin(st, "proj", "pass", []string{"FOO"}, &buf); err != nil {
		t.Fatalf("CmdUnpin: %v", err)
	}

	buf.Reset()
	_ = CmdListPinned(st, "proj", "pass", &buf)
	out := buf.String()
	if strings.Contains(out, "FOO") {
		t.Errorf("FOO should have been unpinned, got: %s", out)
	}
	if !strings.Contains(out, "BAR") {
		t.Errorf("BAR should still be pinned, got: %s", out)
	}
}

func TestCmdPinEmptyProject(t *testing.T) {
	st := newPinStore(t)
	var buf bytes.Buffer
	err := CmdPin(st, "", "pass", []string{"FOO"}, &buf)
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdPinNoKeys(t *testing.T) {
	st := newPinStore(t)
	var buf bytes.Buffer
	err := CmdPin(st, "proj", "pass", []string{}, &buf)
	if err == nil {
		t.Fatal("expected error for no keys")
	}
}

func TestCmdListPinnedNone(t *testing.T) {
	st := newPinStore(t)
	var buf bytes.Buffer
	_ = CmdSet(st, "proj", "pass", "FOO=bar", &buf)
	buf.Reset()

	if err := CmdListPinned(st, "proj", "pass", &buf); err != nil {
		t.Fatalf("CmdListPinned: %v", err)
	}
	if !strings.Contains(buf.String(), "no pinned") {
		t.Errorf("expected 'no pinned' message, got: %s", buf.String())
	}
}
