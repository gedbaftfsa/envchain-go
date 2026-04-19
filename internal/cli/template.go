package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdTemplate renders a project's env vars as a template string,
// substituting {{KEY}} placeholders with stored values.
func CmdTemplate(st *store.Store, project, passphrase, tmpl string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}
	if tmpl == "" {
		return fmt.Errorf("template string required")
	}

	es, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	result := tmpl
	for _, k := range es.Keys() {
		v, _ := es.Get(k)
		result = strings.ReplaceAll(result, "{{{"+k+"}}}", v)
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}

	fmt.Fprintln(w, result)
	return nil
}

// CmdTemplateFile renders a template file, substituting placeholders.
func CmdTemplateFile(st *store.Store, project, passphrase, path string, w io.Writer) error {
	if path == "" {
		return fmt.Errorf("template file path required")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	return CmdTemplate(st, project, passphrase, string(b), w)
}
