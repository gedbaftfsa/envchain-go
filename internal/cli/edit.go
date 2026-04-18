package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/envchain-go/internal/env"
	"github.com/user/envchain-go/internal/store"
)

// CmdEdit opens the project's env set in $EDITOR for interactive editing.
func CmdEdit(st *store.Store, project, passphrase string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	set, err := loadOrNew(st, project, passphrase)
	if err != nil {
		return err
	}

	// Serialise current vars to a temp file
	tmp, err := os.CreateTemp("", "envchain-edit-*.env")
	if err != nil {
		return fmt.Errorf("edit: create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	for _, k := range set.Keys() {
		v, _ := set.Get(k)
		fmt.Fprintf(tmp, "%s=%s\n", k, v)
	}
	tmp.Close()

	// Open editor
	cmd := exec.Command(editor, tmpName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("edit: editor exited: %w", err)
	}

	// Re-read file
	data, err := os.ReadFile(tmpName)
	if err != nil {
		return fmt.Errorf("edit: read temp file: %w", err)
	}

	newSet := env.NewSet()
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, err := env.ParseEntry(line)
		if err != nil {
			return fmt.Errorf("edit: %w", err)
		}
		newSet.Put(k, v)
	}

	return st.Save(project, passphrase, newSet)
}
