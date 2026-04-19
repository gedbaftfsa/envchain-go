package cli_test

import (
	"bytes"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/cli"
	"github.com/nicholasgasior/envchain-go/internal/env"
)

func seedFmt(t *testing.T, st interface{ Save(string, string, *env.Set) error }, project, pass string) {
	t.Helper()
	es := env.NewSet()
	_ = es.Put("ZEBRA", "z")
	_ = es.Put("APPLE", "a")
	_ = es.Put("MANGO", "m")
	if err := st.Save(project, pass, es); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdFmtSorts(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	const pass = "pass"
	seedFmt(t, st, "myapp", pass)

	var buf bytes.Buffer
	if err := cli.CmdFmt(st, "myapp", pass, &buf); err != nil {
		t.Fatalf("CmdFmt: %v", err)
	}

	es, err := st.Load("myapp", pass)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	keys := es.Keys()
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestCmdFmtEmptyProject(t *testing.T) {
	st, _ := newTempStore(t)
	var buf bytes.Buffer
	if err := cli.CmdFmt(st, "", "pass", &buf); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestCmdFmtNotFound(t *testing.T) {
	st, _ := newTempStore(t)
	var buf bytes.Buffer
	if err := cli.CmdFmt(st, "ghost", "pass", &buf); err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestCmdFmtDiffAlreadySorted(t *testing.T) {
	st, _ := newTempStore(t)
	const pass = "pass"
	es := env.NewSet()
	_ = es.Put("AAA", "1")
	_ = es.Put("BBB", "2")
	_ = st.Save("sorted", pass, es)

	var buf bytes.Buffer
	if err := cli.CmdFmtDiff(st, "sorted", pass, &buf); err != nil {
		t.Fatalf("CmdFmtDiff: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("already sorted")) {
		t.Fatalf("expected 'already sorted', got: %s", buf.String())
	}
}

func TestCmdFmtDiffUnsorted(t *testing.T) {
	st, _ := newTempStore(t)
	const pass = "pass"
	seedFmt(t, st, "unsorted", pass)

	var buf bytes.Buffer
	if err := cli.CmdFmtDiff(st, "unsorted", pass, &buf); err != nil {
		t.Fatalf("CmdFmtDiff: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("would reorder")) {
		t.Fatalf("expected reorder message, got: %s", buf.String())
	}
}
