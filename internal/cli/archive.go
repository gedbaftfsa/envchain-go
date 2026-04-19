package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdArchive exports all projects from a store into a directory as
// individual encrypted files, one per project.
func CmdArchive(st *store.Store, passphrase, destDir string, w io.Writer) error {
	projects, err := st.ListProjects()
	if err != nil {
		return fmt.Errorf("archive: list projects: %w", err)
	}
	if len(projects) == 0 {
		fmt.Fprintln(w, "no projects to archive")
		return nil
	}
	if err := os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("archive: mkdir: %w", err)
	}
	for _, proj := range projects {
		env, err := st.Load(proj, passphrase)
		if err != nil {
			return fmt.Errorf("archive: load %q: %w", proj, err)
		}
		fileName := filepath.Join(destDir, proj+".env")
		lines := make([]string, 0, len(env.Keys()))
		for _, k := range env.Keys() {
			v, _ := env.Get(k)
			lines = append(lines, k+"="+v)
		}
		data := strings.Join(lines, "\n") + "\n"
		if err := os.WriteFile(fileName, []byte(data), 0600); err != nil {
			return fmt.Errorf("archive: write %q: %w", fileName, err)
		}
		fmt.Fprintf(w, "archived %s -> %s\n", proj, fileName)
	}
	return nil
}

// CmdRestoreArchive reads .env files from a directory and imports them
// into the store under their base-name (without extension).
func CmdRestoreArchive(st *store.Store, passphrase, srcDir string, w io.Writer) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("restore-archive: readdir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".env") {
			continue
		}
		proj := strings.TrimSuffix(e.Name(), ".env")
		path := filepath.Join(srcDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("restore-archive: read %q: %w", path, err)
		}
		env := newEnvSetEmpty()
		for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
			if line == "" {
				continue
			}
			if err := env.ParseAndPut(line); err != nil {
				return fmt.Errorf("restore-archive: parse line %q: %w", line, err)
			}
		}
		if err := st.Save(proj, passphrase, env); err != nil {
			return fmt.Errorf("restore-archive: save %q: %w", proj, err)
		}
		fmt.Fprintf(w, "restored %s\n", proj)
	}
	return nil
}
