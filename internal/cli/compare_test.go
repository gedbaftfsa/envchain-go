package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
)

func seedCompare(t *testing.T) (*testStore, string) {
	t.Helper()
	st := newTempStore(t)
	pass := "comparepass"

	aSet := newEnvSet(t, "KEY_ONLY_A=1", "SHARED=hello")
	if err := st.Save("proj-a", pass, aSet); err != nil {
		t.Fatalf("save proj-a: %v", err)
	}
	bSet := newEnvSet(t, "KEY_ONLY_B=2", "SHARED=world")
	if err := st.Save("proj-b", pass, bSet); err != nil {
		t.Fatalf("save proj-b: %v", err)
	}
	return st, pass
}

func TestCmdCompareOutput(t *testing.T) {
	st, pass := seedCompare(t)
	var buf bytes.Buffer
	if err := CmdCompare(st.Store, pass, "proj-a", "proj-b", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"KEY_ONLY_A", "KEY_ONLY_B", "SHARED"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected 'missing' in output")
	}
	if !strings.Contains(out, "present") {
		t.Error("expected 'present' in output")
	}
}

func TestCmdCompareProjectNotFound(t *testing.T) {
	st, pass := seedCompare(t)
	var buf bytes.Buffer
	err := CmdCompare(st.Store, pass, "proj-a", "no-such", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdCompareEmptyNames(t *testing.T) {
	st, pass := seedCompare(t)
	var buf bytes.Buffer
	err := CmdCompare(st.Store, pass, "", "proj-b", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdCompareWrongPassphrase(t *testing.T) {
	st, _ := seedCompare(t)
	var buf bytes.Buffer
	err := CmdCompare(st.Store, "wrong", "proj-a", "proj-b", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

// ensure env import is used
var _ = env.NewSet
