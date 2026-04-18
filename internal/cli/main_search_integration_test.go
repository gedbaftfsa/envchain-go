package cli

import (
	"testing"
)

func TestMainSearchDispatch(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	seedProjectInStore(t, st, "myproject", "hunter2", map[string]string{
		"AWS_ACCESS_KEY": "AKIA...",
		"AWS_SECRET_KEY": "secret",
	})

	err := CmdSearch(st, "hunter2", "AWS")
	if err != nil {
		t.Fatalf("dispatch search: %v", err)
	}
}

func TestMainSearchWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := storeFromDir(t, dir)
	seedProjectInStore(t, st, "proj", "correct", map[string]string{"KEY": "val"})

	err := CmdSearch(st, "wrong", "KEY")
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
