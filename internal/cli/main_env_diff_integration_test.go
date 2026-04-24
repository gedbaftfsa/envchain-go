package cli

import (
	"strings"
	"testing"
)

func TestMainEnvDiffDispatch(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)

	es := mustEnvSet(t, "DISPATCH_KEY=dispatch_val")
	if err := st.Save("dispatchproj", "secret", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Setenv("DISPATCH_KEY", "dispatch_val")

	passReader := passphraseReader("secret")
	out := captureStdout(t, func() {
		err := Main(
			[]string{"envchain", "env-diff", "dispatchproj"},
			st,
			passReader,
		)
		if err != nil {
			t.Fatalf("Main: %v", err)
		}
	})

	if !strings.Contains(out, "no differences") {
		t.Errorf("expected 'no differences', got: %s", out)
	}
}

func TestMainEnvDiffWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := mustStoreDir(t, dir)

	es := mustEnvSet(t, "K=v")
	if err := st.Save("proj", "correct", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	passReader := passphraseReader("wrong")
	err := Main(
		[]string{"envchain", "env-diff", "proj"},
		st,
		passReader,
	)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
