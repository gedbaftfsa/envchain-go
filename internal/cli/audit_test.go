package cli

import (
	"bytes"
	"strings"
	"testing"
)

func seedAudit(t *testing.T) *storeT {
	t.Helper()
	st := newTempStore(t)
	if err := CmdSet(st.s, "secret", "myproject", []string{"FOO=bar", "BAZ=qux"}, &bytes.Buffer{}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return st
}

func TestCmdAuditSuccess(t *testing.T) {
	st := seedAudit(t)
	var buf bytes.Buffer
	if err := CmdAudit(st.s, "secret", "myproject", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "myproject") {
		t.Errorf("expected project name in output, got: %s", out)
	}
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAZ") {
		t.Errorf("expected key names in output, got: %s", out)
	}
	if !strings.Contains(out, "keys    : 2") {
		t.Errorf("expected key count 2, got: %s", out)
	}
}

func TestCmdAuditWrongPassphrase(t *testing.T) {
	st := seedAudit(t)
	var buf bytes.Buffer
	err := CmdAudit(st.s, "wrong", "myproject", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdAuditNotFound(t *testing.T) {
	st := seedAudit(t)
	var buf bytes.Buffer
	err := CmdAudit(st.s, "secret", "noproject", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdAuditEmptyProject(t *testing.T) {
	st := seedAudit(t)
	var buf bytes.Buffer
	err := CmdAudit(st.s, "secret", "", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
