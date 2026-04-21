package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
)

func seedEnvCheck(t *testing.T, st interface{ Save(string, interface{}, string) error }, project, pass string) {
	t.Helper()
	es := env.NewSet()
	es.Put("APP_HOST", "localhost")
	es.Put("APP_PORT", "8080")
	if err := mustEnvSet(st, project, pass, es); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func TestCmdEnvCheckMatch(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	pass := "pass"
	seedEnvCheck(t, st, "myapp", pass)

	t.Setenv("APP_HOST", "localhost")
	t.Setenv("APP_PORT", "8080")

	var buf bytes.Buffer
	err := CmdEnvCheck(st, "myapp", pass, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "matches") {
		t.Errorf("expected match message, got: %q", buf.String())
	}
}

func TestCmdEnvCheckMissing(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	pass := "pass"
	seedEnvCheck(t, st, "myapp", pass)

	t.Setenv("APP_HOST", "localhost")
	// APP_PORT intentionally not set

	var buf bytes.Buffer
	err := CmdEnvCheck(st, "myapp", pass, &buf)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(buf.String(), "missing") {
		t.Errorf("expected missing message, got: %q", buf.String())
	}
}

func TestCmdEnvCheckMismatch(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	pass := "pass"
	seedEnvCheck(t, st, "myapp", pass)

	t.Setenv("APP_HOST", "remotehost")
	t.Setenv("APP_PORT", "8080")

	var buf bytes.Buffer
	err := CmdEnvCheck(st, "myapp", pass, &buf)
	if err == nil {
		t.Fatal("expected error for mismatched value")
	}
	if !strings.Contains(buf.String(), "mismatch") {
		t.Errorf("expected mismatch message, got: %q", buf.String())
	}
}

func TestCmdEnvCheckEmptyProject(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	var buf bytes.Buffer
	err := CmdEnvCheck(st, "", "pass", &buf)
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdEnvCheckWrongPassphrase(t *testing.T) {
	st, dir := newTempStore(t)
	_ = dir
	pass := "correct"
	seedEnvCheck(t, st, "myapp", pass)

	var buf bytes.Buffer
	err := CmdEnvCheck(st, "myapp", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
