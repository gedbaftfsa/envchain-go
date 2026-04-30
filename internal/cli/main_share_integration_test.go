package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

func shareIntegrationStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestMainShareDispatch(t *testing.T) {
	st := shareIntegrationStore(t)
	set := env.NewSet()
	set.Put("INTEGRATION", "yes")
	_ = st.Save("intproj", "secret", set)

	var buf bytes.Buffer
	err := CmdShare(st, "intproj", "secret", &buf)
	if err != nil {
		t.Fatalf("CmdShare dispatch: %v", err)
	}
	if !strings.Contains(buf.String(), "INTEGRATION=") {
		t.Errorf("output missing INTEGRATION key: %s", buf.String())
	}
}

func TestMainShareWrongPassphrase(t *testing.T) {
	st := shareIntegrationStore(t)
	set := env.NewSet()
	set.Put("K", "v")
	_ = st.Save("proj", "correct", set)

	var buf bytes.Buffer
	err := CmdShare(st, "proj", "wrong", &buf)
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}
