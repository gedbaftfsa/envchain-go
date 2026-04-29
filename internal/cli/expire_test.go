package cli

import {
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
}

func newExpireStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func seedExpire(t *testing.T, st *store.Store) {
	t.Helper()
	set := mustEnvSet(t, "KEY=value")
	if err := st.Save("myproject", "pass", set); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdExpireSetsDate(t *testing.T) {
	st := newExpireStore(t)
	seedExpire(t, st)

	var buf bytes.Buffer
	if err := CmdExpire(st, "pass", "myproject", []string{"30"}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	expected := time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02")
	if !strings.Contains(out, expected) {
		t.Errorf("expected output to contain %q, got: %s", expected, out)
	}
}

func TestCmdExpireInvalidDays(t *testing.T) {
	st := newExpireStore(t)
	seedExpire(t, st)

	var buf bytes.Buffer
	if err := CmdExpire(st, "pass", "myproject", []string{"abc"}, &buf); err == nil {
		t.Fatal("expected error for non-integer days")
	}
}

func TestCmdExpireEmptyProject(t *testing.T) {
	st := newExpireStore(t)
	var buf bytes.Buffer
	if err := CmdExpire(st, "pass", "", []string{"7"}, &buf); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdCheckExpiryNotExpired(t *testing.T) {
	st := newExpireStore(t)
	seedExpire(t, st)

	var buf bytes.Buffer
	if err := CmdExpire(st, "pass", "myproject", []string{"10"}, &buf); err != nil {
		t.Fatalf("set expiry: %v", err)
	}

	buf.Reset()
	if err := CmdCheckExpiry(st, "pass", "myproject", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "day(s) remaining") {
		t.Errorf("expected remaining days in output, got: %s", buf.String())
	}
}

func TestCmdCheckExpiryNoExpiry(t *testing.T) {
	st := newExpireStore(t)
	seedExpire(t, st)

	var buf bytes.Buffer
	if err := CmdCheckExpiry(st, "pass", "myproject", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no expiry set") {
		t.Errorf("expected 'no expiry set', got: %s", buf.String())
	}
}

func TestCmdCheckExpiryWrongPassphrase(t *testing.T) {
	st := newExpireStore(t)
	seedExpire(t, st)

	var buf bytes.Buffer
	_ = fmt.Sprintf("%v", buf) // suppress unused warning
	if err := CmdCheckExpiry(st, "wrong", "myproject", &buf); err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
