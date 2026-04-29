package cli

import (
	"bytes"
	"testing"

	"github.com/your-org/envchain-go/internal/store"
)

func newProtectStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func seedProtect(t *testing.T, st *store.Store) {
	t.Helper()
	set := mustEnvSet(t, "API_KEY=secret", "DB_PASS=hunter2", "PORT=8080")
	if err := st.Save("myapp", "pass", set); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdProtectAndList(t *testing.T) {
	st := newProtectStore(t)
	seedProtect(t, st)

	var buf bytes.Buffer
	if err := CmdProtect(st, "myapp", "pass", []string{"API_KEY", "DB_PASS"}, &buf); err != nil {
		t.Fatalf("protect: %v", err)
	}
	if got := buf.String(); got == "" {
		t.Error("expected non-empty output from protect")
	}

	buf.Reset()
	if err := CmdListProtected(st, "myapp", "pass", &buf); err != nil {
		t.Fatalf("list-protected: %v", err)
	}
	out := buf.String()
	if !contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in protected list, got: %s", out)
	}
	if !contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in protected list, got: %s", out)
	}
}

func TestCmdUnprotect(t *testing.T) {
	st := newProtectStore(t)
	seedProtect(t, st)

	var buf bytes.Buffer
	_ = CmdProtect(st, "myapp", "pass", []string{"API_KEY", "DB_PASS"}, &buf)

	buf.Reset()
	if err := CmdUnprotect(st, "myapp", "pass", []string{"API_KEY"}, &buf); err != nil {
		t.Fatalf("unprotect: %v", err)
	}

	buf.Reset()
	_ = CmdListProtected(st, "myapp", "pass", &buf)
	out := buf.String()
	if contains(out, "API_KEY") {
		t.Errorf("API_KEY should no longer be protected, got: %s", out)
	}
	if !contains(out, "DB_PASS") {
		t.Errorf("DB_PASS should still be protected, got: %s", out)
	}
}

func TestCmdListProtectedEmpty(t *testing.T) {
	st := newProtectStore(t)
	seedProtect(t, st)

	var buf bytes.Buffer
	if err := CmdListProtected(st, "myapp", "pass", &buf); err != nil {
		t.Fatalf("list-protected: %v", err)
	}
	if !contains(buf.String(), "no protected keys") {
		t.Errorf("expected 'no protected keys' message, got: %s", buf.String())
	}
}

func TestCmdProtectEmptyProject(t *testing.T) {
	st := newProtectStore(t)
	var buf bytes.Buffer
	if err := CmdProtect(st, "", "pass", []string{"KEY"}, &buf); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestCmdProtectNoKeys(t *testing.T) {
	st := newProtectStore(t)
	seedProtect(t, st)
	var buf bytes.Buffer
	if err := CmdProtect(st, "myapp", "pass", []string{}, &buf); err == nil {
		t.Error("expected error for empty key list")
	}
}

func TestCmdProtectWrongPassphrase(t *testing.T) {
	st := newProtectStore(t)
	seedProtect(t, st)
	var buf bytes.Buffer
	if err := CmdProtect(st, "myapp", "wrong", []string{"API_KEY"}, &buf); err == nil {
		t.Error("expected error for wrong passphrase")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
