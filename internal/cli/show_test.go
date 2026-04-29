package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envchain-go/internal/env"
	"github.com/user/envchain-go/internal/store"
)

func newShowStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func seedShow(t *testing.T, st *store.Store, project string) {
	t.Helper()
	set := env.NewSet()
	set.Put("ALPHA", "one")
	set.Put("BETA", "two")
	set.Put("GAMMA", "three")
	if err := st.Save(project, set, "pass"); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdShowSuccess(t *testing.T) {
	st := newShowStore(t)
	seedShow(t, st, "myproject")

	var buf bytes.Buffer
	if err := CmdShow(st, "myproject", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALPHA=one") {
		t.Errorf("expected ALPHA=one in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BETA=two") {
		t.Errorf("expected BETA=two in output, got:\n%s", out)
	}
	if !strings.Contains(out, "GAMMA=three") {
		t.Errorf("expected GAMMA=three in output, got:\n%s", out)
	}
	if !strings.Contains(out, "# project: myproject (3 variable(s))") {
		t.Errorf("expected header line in output, got:\n%s", out)
	}
}

func TestCmdShowSorted(t *testing.T) {
	st := newShowStore(t)
	seedShow(t, st, "sorted")

	var buf bytes.Buffer
	if err := CmdShow(st, "sorted", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var varLines []string
	for _, l := range lines {
		if !strings.HasPrefix(l, "#") {
			varLines = append(varLines, l)
		}
	}
	for i := 1; i < len(varLines); i++ {
		if varLines[i] < varLines[i-1] {
			t.Errorf("output not sorted: %q before %q", varLines[i-1], varLines[i])
		}
	}
}

func TestCmdShowNotFound(t *testing.T) {
	st := newShowStore(t)
	var buf bytes.Buffer
	err := CmdShow(st, "ghost", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCmdShowWrongPassphrase(t *testing.T) {
	st := newShowStore(t)
	seedShow(t, st, "locked")
	var buf bytes.Buffer
	err := CmdShow(st, "locked", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdShowEmptyProject(t *testing.T) {
	st := newShowStore(t)
	var buf bytes.Buffer
	err := CmdShow(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdShowEmptyVars(t *testing.T) {
	st := newShowStore(t)
	set := env.NewSet()
	if err := st.Save("empty", set, "pass"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var buf bytes.Buffer
	if err := CmdShow(st, "empty", "pass", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "has no variables") {
		t.Errorf("expected empty-project message, got: %s", buf.String())
	}
}
