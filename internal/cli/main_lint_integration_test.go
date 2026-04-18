package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
)

func TestMainLintDispatch(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)

	set := env.NewSet()
	set.Put("FOO", "bar")
	if err := st.Save("myproject", "s3cr3t", set); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := CmdLint(st, "myproject", "s3cr3t", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK output, got: %s", buf.String())
	}
}

func TestMainLintWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)

	set := env.NewSet()
	set.Put("FOO", "bar")
	if err := st.Save("myproject", "correct", set); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err := CmdLint(st, "myproject", "wrong", &buf)
	if err == nil {
		t.Fatal("expected decryption error")
	}
}
