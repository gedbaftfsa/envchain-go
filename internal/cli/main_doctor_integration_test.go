package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestMainDoctorDispatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ENVCHAIN_DIR", dir)

	var buf bytes.Buffer
	code := Main([]string{"envchain", "doctor"}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\noutput: %s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "[ OK ]") {
		t.Errorf("expected OK lines in output:\n%s", buf.String())
	}
}

func TestMainDoctorHelp(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCHAIN_DIR", dir)
	var buf bytes.Buffer
	Main([]string{"envchain", "doctor", "--help"}, &buf)
	if !strings.Contains(buf.String(), "doctor") {
		t.Errorf("expected doctor in help output:\n%s", buf.String())
	}
}
