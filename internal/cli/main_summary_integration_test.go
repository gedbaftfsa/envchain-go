package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
	"github.com/envchain/envchain-go/internal/store"
)

func TestMainSummaryDispatch(t *testing.T) {
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}

	set := env.NewSet()
	_ = set.Put("HELLO", "world")
	if err := st.Save("myapp", set, "pw"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var buf bytes.Buffer
	err = CmdSummary(st, "pw", &buf)
	if err != nil {
		t.Fatalf("CmdSummary: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected myapp in summary output, got: %q", out)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected key count 1 in summary output, got: %q", out)
	}
}

func TestMainSummaryWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st, _ := store.New(dir)
	set := env.NewSet()
	_ = set.Put("K", "V")
	_ = st.Save("proj", set, "correct")

	err := CmdSummary(st, "bad", nil)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}
