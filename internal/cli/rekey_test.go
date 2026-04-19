package cli

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
)

func seedRekey(t *testing.T, st interface{ Save(string, *env.Set, string) error }, pass string, names ...string) {
	t.Helper()
	for _, name := range names {
		set := newEnvSet(t, "KEY=val")
		if err := st.(interface {
			Save(string, *env.Set, string) error
		}).Save(name, set, pass); err != nil {
			t.Fatalf("seed %q: %v", name, err)
		}
	}
}

func TestCmdRekeySuccess(t *testing.T) {
	st := newTempStore(t)
	seedRekey(t, st, "old", "alpha", "beta")

	var buf bytes.Buffer
	if err := CmdRekey(st, "old", "new", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("2 project(s)")) {
		t.Errorf("expected count in output, got: %s", out)
	}

	// verify new passphrase works
	for _, name := range []string{"alpha", "beta"} {
		if _, err := st.Load(name, "new"); err != nil {
			t.Errorf("load %q with new pass: %v", name, err)
		}
	}
}

func TestCmdRekeyWrongOldPassphrase(t *testing.T) {
	st := newTempStore(t)
	seedRekey(t, st, "correct", "proj")

	var buf bytes.Buffer
	err := CmdRekey(st, "wrong", "new", &buf)
	if err == nil {
		t.Fatal("expected error for wrong old passphrase")
	}
}

func TestCmdRekeySamePassphrase(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	err := CmdRekey(st, "pass", "pass", &buf)
	if err == nil {
		t.Fatal("expected error when passphrases are identical")
	}
}

func TestCmdRekeyEmpty(t *testing.T) {
	st := newTempStore(t)
	var buf bytes.Buffer
	if err := CmdRekey(st, "old", "new", &buf); err != nil {
		t.Fatalf("unexpected error on empty store: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("no projects")) {
		t.Errorf("expected 'no projects' message, got: %s", buf.String())
	}
}
