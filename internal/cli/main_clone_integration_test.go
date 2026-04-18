package cli

import (
	"os"
	"strings"
	"testing"
)

func TestMainCloneDispatch(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCHAIN_DIR", dir)

	s, err := defaultStore()
	if err != nil {
		t.Fatal(err)
	}

	const pass = "clonepass"
	seedClone(t, s, "original", pass)

	// Simulate: envchain clone original copy
	readPass = func(_ string) (string, error) { return pass, nil }
	defer func() { readPass = ReadPassphrase }()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = Main([]string{"envchain", "clone", "original", "copy"})
	w.Close()
	os.Stdout = old

	var buf strings.Builder
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("Main clone: %v", err)
	}
	if !strings.Contains(buf.String(), "original") {
		t.Errorf("expected output to mention source project, got: %s", buf.String())
	}
}
