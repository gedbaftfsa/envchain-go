package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCmdCompletionBash(t *testing.T) {
	var buf bytes.Buffer
	if err := CmdCompletion("bash", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "_envchain_go") {
		t.Error("bash completion missing function name")
	}
	if !strings.Contains(out, "complete -F") {
		t.Error("bash completion missing complete directive")
	}
}

func TestCmdCompletionZsh(t *testing.T) {
	var buf bytes.Buffer
	if err := CmdCompletion("zsh", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "#compdef") {
		t.Error("zsh completion missing #compdef header")
	}
}

func TestCmdCompletionFish(t *testing.T) {
	var buf bytes.Buffer
	if err := CmdCompletion("fish", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "envchain-go") {
		t.Error("fish completion missing binary name")
	}
}

func TestCmdCompletionCaseInsensitive(t *testing.T) {
	var buf bytes.Buffer
	if err := CmdCompletion("BASH", &buf); err != nil {
		t.Fatalf("expected case-insensitive match, got error: %v", err)
	}
}

func TestCmdCompletionUnknownShell(t *testing.T) {
	var buf bytes.Buffer
	err := CmdCompletion("powershell", &buf)
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("unexpected error message: %v", err)
	}
}
