package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainArchiveDispatch(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)
	env := newEnvSet(t, "KEY=val")
	if err := st.Save("myproj", "pw", env); err != nil {
		t.Fatalf("save: %v", err)
	}

	dest := filepath.Join(dir, "archive")
	out := captureStdout(t, func() {
		os.Args = []string{"envchain", "archive", "myproj", dest}
		Main(st, passphraseReader("pw"))
	})
	if !strings.Contains(out, "archived myproj") {
		t.Errorf("expected archive output, got: %s", out)
	}
	if _, err := os.Stat(filepath.Join(dest, "myproj.env")); err != nil {
		t.Errorf("archive file not created: %v", err)
	}
}

func TestMainArchiveWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)
	env := newEnvSet(t, "K=v")
	if err := st.Save("p", "right", env); err != nil {
		t.Fatalf("save: %v", err)
	}

	defer func() { recover() }()
	os.Args = []string{"envchain", "archive", "p", t.TempDir()}
	Main(st, passphraseReader("wrong"))
}
