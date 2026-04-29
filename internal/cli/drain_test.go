package cli

import (
	"bytes"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
	"github.com/envchain/envchain-go/internal/store"
)

func newDrainStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func seedDrain(t *testing.T, st *store.Store, project, pass string) {
	t.Helper()
	es := env.NewSet()
	es.Put("ALPHA", "1")
	es.Put("BETA", "2")
	es.Put("GAMMA", "3")
	if err := st.Save(project, pass, es); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdDrainSuccess(t *testing.T) {
	st := newDrainStore(t)
	seedDrain(t, st, "myproject", "secret")

	var buf bytes.Buffer
	if err := CmdDrain(st, "myproject", "secret", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Fatal("expected output, got none")
	}

	// Project should still exist but be empty.
	es, err := st.Load("myproject", "secret")
	if err != nil {
		t.Fatalf("load after drain: %v", err)
	}
	if len(es.Keys()) != 0 {
		t.Fatalf("expected empty set, got keys: %v", es.Keys())
	}
}

func TestCmdDrainAlreadyEmpty(t *testing.T) {
	st := newDrainStore(t)
	es := env.NewSet()
	if err := st.Save("empty", "pass", es); err != nil {
		t.Fatalf("save: %v", err)
	}

	var buf bytes.Buffer
	if err := CmdDrain(st, "empty", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() == "" {
		t.Fatal("expected informational output")
	}
}

func TestCmdDrainWrongPassphrase(t *testing.T) {
	st := newDrainStore(t)
	seedDrain(t, st, "proj", "correct")

	var buf bytes.Buffer
	if err := CmdDrain(st, "proj", "wrong", &buf); err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdDrainNotFound(t *testing.T) {
	st := newDrainStore(t)

	var buf bytes.Buffer
	if err := CmdDrain(st, "ghost", "pass", &buf); err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdDrainEmptyProject(t *testing.T) {
	st := newDrainStore(t)

	var buf bytes.Buffer
	if err := CmdDrain(st, "", "pass", &buf); err == nil {
		t.Fatal("expected error for empty project name")
	}
}
