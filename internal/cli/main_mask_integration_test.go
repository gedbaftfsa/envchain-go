package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envchain/envchain-go/internal/env"
)

func TestMainMaskDispatch(t *testing.T) {
	dir := t.TempDir()
	st := mustOpenStore(t, dir)

	es := env.NewSet()
	es.Put("SECRET", "topsecret")
	if err := st.Save("myapp", "hunter2", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var out bytes.Buffer
	passIdx := 0
	phrases := []string{"hunter2"}
	readPass := func(_ string) (string, error) {
		p := phrases[passIdx]
		passIdx++
		return p, nil
	}

	err := Main([]string{"envchain", "mask", "myapp"}, st, readPass, &out)
	if err != nil {
		t.Fatalf("Main mask: %v", err)
	}
	if strings.Contains(out.String(), "topsecret") {
		t.Error("plaintext secret must not appear in mask output")
	}
	if !strings.Contains(out.String(), "SECRET=") {
		t.Errorf("expected SECRET key in output, got: %s", out.String())
	}
}

func TestMainMaskWrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	st := mustOpenStore(t, dir)

	es := env.NewSet()
	es.Put("K", "v")
	if err := st.Save("proj", "correct", es); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var out bytes.Buffer
	readPass := func(_ string) (string, error) { return "wrong", nil }

	err := Main([]string{"envchain", "mask", "proj"}, st, readPass, &out)
	if err == nil {
		t.Error("expected error for wrong passphrase")
	}
}
