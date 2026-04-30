package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestMainNamespaceDispatch(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	// Verify the namespace command is reachable through the command registry.
	var buf bytes.Buffer
	err := CmdNamespace(st, "myapp", "pass", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB") {
		t.Errorf("expected DB in output, got: %s", buf.String())
	}
}

func TestMainNamespaceWrongPassphrase(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	var buf bytes.Buffer
	err := CmdNamespace(st, "myapp", "badpass", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestMainNamespaceKeysDispatch(t *testing.T) {
	st := newNamespaceStore(t)
	seedNamespace(t, st)

	var buf bytes.Buffer
	err := CmdNamespaceKeys(st, "myapp", "AWS", "pass", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "AWS_KEY") || !strings.Contains(out, "AWS_SECRET") {
		t.Errorf("expected AWS keys in output, got: %s", out)
	}
}
