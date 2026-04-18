package cli

import (
	"testing"

	"github.com/envchain-go/internal/env"
)

func TestCmdRotateSuccess(t *testing.T) {
	st := newTempStore(t)

	// Seed a project with old passphrase
	set := env.NewSet()
	_ = set.Put("KEY", "value")
	if err := st.Save("myapp", set, "old-pass"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Rotate to new passphrase
	if err := CmdRotate(st, "myapp", "old-pass", "new-pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Old passphrase should no longer work
	if _, err := st.Load("myapp", "old-pass"); err == nil {
		t.Fatal("expected error with old passphrase after rotation")
	}

	// New passphrase should work and data intact
	loaded, err := st.Load("myapp", "new-pass")
	if err != nil {
		t.Fatalf("load with new pass: %v", err)
	}
	if v, _ := loaded.Get("KEY"); v != "value" {
		t.Errorf("expected KEY=value, got %q", v)
	}
}

func TestCmdRotateSamePassphrase(t *testing.T) {
	st := newTempStore(t)
	set := env.NewSet()
	_ = st.Save("proj", set, "pass")

	err := CmdRotate(st, "proj", "pass", "pass")
	if err == nil {
		t.Fatal("expected error when old and new passphrase are the same")
	}
}

func TestCmdRotateWrongOldPassphrase(t *testing.T) {
	st := newTempStore(t)
	set := env.NewSet()
	_ = st.Save("proj", set, "correct")

	err := CmdRotate(st, "proj", "wrong", "new-pass")
	if err == nil {
		t.Fatal("expected error with wrong old passphrase")
	}
}

func TestCmdRotateEmptyProject(t *testing.T) {
	st := newTempStore(t)
	err := CmdRotate(st, "", "old", "new")
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}
