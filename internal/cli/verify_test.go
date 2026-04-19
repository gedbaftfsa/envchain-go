package cli

import (
	"bytes"
	"testing"

	"github.com/envchain-go/internal/store"
)

func seedVerify(t *testing.T, st *store.Store) {
	t.Helper()
	es := mustEnvSet(t, map[string]string{"TOKEN": "abc123"})
	if err := st.Save("myproject", es, "s3cr3t"); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdVerifySuccess(t *testing.T) {
	st := newTempStore(t)
	seedVerify(t, st)
	var buf bytes.Buffer
	if err := CmdVerify(st, "myproject", "s3cr3t", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got == "" {
		t.Error("expected non-empty output")
	}
}

func TestCmdVerifyWrongPassphrase(t *testing.T) {
	st := newTempStore(t)
	seedVerify(t, st)
	var buf bytes.Buffer
	err := CmdVerify(st, "myproject", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdVerifyNotFound(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdVerify(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdVerifyEmptyProject(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdVerify(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
