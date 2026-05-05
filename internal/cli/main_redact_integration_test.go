package cli

import (
	"bytes"
	"testing"

	"github.com/envchain-go/internal/env"
)

func TestMainRedactDispatch(t *testing.T) {
	st := newRedactStore(t)

	es := env.NewSet()
	_ = es.Put("TOKEN", "topsecret")
	if err := st.Save("proj", "pw", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var buf bytes.Buffer
	err := CmdRedact(st, "proj", "pw", "my token is topsecret ok", &buf)
	if err != nil {
		t.Fatalf("CmdRedact: %v", err)
	}

	got := buf.String()
	if contains(got, "topsecret") {
		t.Errorf("secret leaked in output: %q", got)
	}
}

func TestMainRedactWrongPassphrase(t *testing.T) {
	st := newRedactStore(t)
	seedRedact(t, st)

	var buf bytes.Buffer
	err := CmdRedact(st, "myapp", "badpass", "irrelevant", &buf)
	if err == nil {
		t.Fatal("expected error")
	}
}
