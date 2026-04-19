package cli

import (
	"bytes"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
)

func seedPrune(t *testing.T, st interface{ Save(string, interface{}, string) error }) {
	t.Helper()
}

func TestCmdPruneRemovesEmpty(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	pass := "prunepass"

	// Save a project with keys.
	es := newEnvSet(t, pass)
	es.Put("KEY", "val")
	if err := st.Save("has-keys", es, pass); err != nil {
		t.Fatal(err)
	}

	// Save an empty project.
	empty := mustEnvSet(t, pass)
	if err := st.Save("empty-proj", empty, pass); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdPrune(st, pass, &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("empty-proj")) {
		t.Errorf("expected 'empty-proj' in output, got: %s", out)
	}

	// Confirm empty-proj is gone.
	projects, err := st.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range projects {
		if p == "empty-proj" {
			t.Error("empty-proj should have been pruned")
		}
	}
}

func TestCmdPruneNothingToDo(t *testing.T) {
	st, _ := newTempStore(t)
	pass := "prunepass2"

	es := newEnvSet(t, pass)
	es.Put("A", "1")
	if err := st.Save("full", es, pass); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := CmdPrune(st, pass, &buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("nothing to prune")) {
		t.Errorf("expected 'nothing to prune', got: %s", buf.String())
	}
}

func TestCmdPruneWrongPassphrase(t *testing.T) {
	st, _ := newTempStore(t)
	pass := "correct"

	es := newEnvSet(t, pass)
	es.Put("X", "y")
	if err := st.Save("proj", es, pass); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := CmdPrune(st, "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func newEnvSet(t *testing.T, pass string) *env.Set {
	t.Helper()
	es, err := env.NewSet(pass)
	if err != nil {
		t.Fatal(err)
	}
	return es
}
