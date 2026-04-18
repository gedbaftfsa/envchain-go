// Package cli wires together store, env, and crypto into user-facing commands.
package cli

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadPassphrase reads a passphrase from the terminal without echo.
// prompt is printed to stderr before reading.
func ReadPassphrase(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Fallback: read a plain line (e.g. piped input in tests).
		var line string
		_, err := fmt.Fscanln(os.Stdin, &line)
		if err != nil {
			return "", fmt.Errorf("reading passphrase: %w", err)
		}
		return strings.TrimRight(line, "\r\n"), nil
	}
	bytes, err := term.ReadPassword(fd)
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %w", err)
	}
	fmt.Fprintln(os.Stderr) // newline after hidden input
	return string(bytes), nil
}

// ReadPassphraseConfirm reads a passphrase twice and returns an error if they
// do not match. Useful when creating a new project namespace.
func ReadPassphraseConfirm(prompt, confirmPrompt string) (string, error) {
	first, err := ReadPassphrase(prompt)
	if err != nil {
		return "", err
	}
	second, err := ReadPassphrase(confirmPrompt)
	if err != nil {
		return "", err
	}
	if first != second {
		return "", fmt.Errorf("passphrases do not match")
	}
	return first, nil
}
