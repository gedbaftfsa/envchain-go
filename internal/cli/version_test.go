package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCmdVersionDefault(t *testing.T) {
	out := captureStdout(Cchain-go") {
		t.Errorf("expected binary name in output, got: %s", out)
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("expected default version 'dev' in output, got: %s", out)
	}
}

func TestCmdVersionCustom(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origDate := BuildDate
	t.Cleanup(func() {
		Version = origVersion
		Commit = origCommit
		BuildDate = origDate
	})

	Version = "1.2.3"
	Commit = "abc1234"
	BuildDate = "2024-01-01"

	out := captureStdout(CmdVersion)
	if !strings.Contains(out, "1.2.3") {
		t.Errorf("expected version 1.2.3 in output, got: %s", out)
	}
	if !strings.Contains(out, "abc1234") {
		t.Errorf("expected commit abc1234 in output, got: %sn	if !strings.Contains(out, "2024-01-01") {
		t.Errorf("expected build date in output, got: %s", out)
	}
}
