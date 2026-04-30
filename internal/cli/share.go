package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdShare prints a portable, passphrase-protected export of a project's
// environment variables to stdout in a format that CmdImport can consume.
// The output is identical to CmdExport but prefixed with a one-line header
// that carries the project name, making it easy to pipe between machines.
//
// Usage: envchain share <project> [passphrase]
func CmdShare(st *store.Store, project, passphrase string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	if len(keys) == 0 {
		return fmt.Errorf("project %q has no variables to share", project)
	}

	fmt.Fprintf(w, "# envchain-share project=%s\n", project)
	for _, k := range keys {
		v, _ := set.Get(k)
		fmt.Fprintf(w, "%s=%s\n", k, shellQuote(v))
	}
	return nil
}

// CmdReceive reads a share blob from r (produced by CmdShare) and saves it
// into st under the given project name (or the embedded name if destProject
// is empty) protected by newPassphrase.
func CmdReceive(st *store.Store, r io.Reader, destProject, newPassphrase string, w io.Writer) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading share input: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty share input")
	}

	embeddedProject := ""
	start := 0
	if strings.HasPrefix(lines[0], "# envchain-share project=") {
		embeddedProject = strings.TrimPrefix(lines[0], "# envchain-share project=")
		start = 1
	}

	target := destProject
	if target == "" {
		target = embeddedProject
	}
	if target == "" {
		return fmt.Errorf("cannot determine target project name; pass one explicitly")
	}

	set, _ := loadOrNew(st, target, newPassphrase)
	for _, line := range lines[start:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, err := parseEntryLine(line)
		if err != nil {
			return fmt.Errorf("parsing line %q: %w", line, err)
		}
		set.Put(k, v)
	}

	if err := st.Save(target, newPassphrase, set); err != nil {
		return err
	}
	fmt.Fprintf(w, "received %d variable(s) into project %q\n", len(set.Keys()), target)
	return nil
}

func parseEntryLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("missing '='")
	}
	return line[:idx], strings.Trim(line[idx+1:], "'"), nil
}

func init() {
	_ = os.Stdout // ensure os import is used
}
