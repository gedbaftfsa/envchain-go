package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/store"
)

func newInitStore(t *testing.T) *store.Store {
	t.Helper()
	return newTempStore(t)
}

func passphraseSeq(vals ...string) func(string) (string, error) {
	i := 0
	return func(_ string) (string, error) {
		v := vals[i]
		i++
		return v, nil
	}
}

func TestCmdInitSuccess(t *testing.T) {
	st := newInitStore(t)
	var buf bytes.Buffer
	err := CmdInit(st, "myapp", &buf, passphraseSeq("secret", "secret"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "myapp") {
		t.Errorf("expected output to mention project name, got: %q", buf.String())
	}
}

func TestCmdInitEmptyProject(t *testing.T) {
	st := newInitStore(t)
	var buf bytes.Buffer
	err := CmdInit(st, "", &buf, passphraseSeq())
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestCmdInitMismatchedPassphrase(t *testing.T) {
	st := newInitStore(t)
	var buf bytes.Buffer
	err := CmdInit(st, "proj", &buf, passphraseSeq("abc", "xyz"))
	if err == nil {
		t.Fatal("expected error for mismatched passphrases")
	}
	if !strings.Contains(err.Error(), "do not match") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCmdInitEmptyPassphrase(t *testing.T) {
	st := newInitStore(t)
	var buf bytes.Buffer
	err := CmdInit(st, "proj", &buf, passphraseSeq("", ""))
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}
