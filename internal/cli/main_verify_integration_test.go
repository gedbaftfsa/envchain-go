package cli

import (
	"bytes"
	"testing"
)

func TestMainVerifyDispatch(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	seedVerify(t, st)

	var out bytes.Buffer
	passFn := func(_ string) (string, error) { return "s3cr3t", nil }

	err := CmdVerify(st, "myproject", "s3cr3t", &out)
	if err != nil {
		t.Fatalf("dispatch verify: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected output from verify")
	}
	_ = passFn
}

func TestMainVerifyWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	seedVerify(t, st)

	var out bytes.Buffer
	err := CmdVerify(st, "myproject", "badpass", &out)
	if err == nil {
		t.Fatal("expected failure with wrong passphrase")
	}
}
