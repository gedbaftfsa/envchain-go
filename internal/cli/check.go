package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envchain-go/internal/store"
)

// CmdCheck verifies that all keys in a project have non-empty values and
// optionally checks that a required set of keys are present.
//
// Usage: envchain check <project> [required_key...]
func CmdCheck(st *store.Store, passphrase, project string, required []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name is required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	keys := set.Keys()
	present := make(map[string]bool, len(keys))
	for _, k := range keys {
		present[k] = true
	}

	var issues []string

	// Check for empty values.
	for _, k := range keys {
		v, _ := set.Get(k)
		if strings.TrimSpace(v) == "" {
			issues = append(issues, fmt.Sprintf("key %q has an empty value", k))
		}
	}

	// Check required keys are present.
	for _, req := range required {
		if !present[req] {
			issues = append(issues, fmt.Sprintf("required key %q is missing", req))
		}
	}

	if len(issues) == 0 {
		fmt.Fprintf(w, "project %q passed all checks (%d key(s))\n", project, len(keys))
		return nil
	}

	for _, iss := range issues {
		fmt.Fprintf(w, "WARN: %s\n", iss)
	}
	fmt.Fprintf(os.Stderr, "project %q has %d issue(s)\n", project, len(issues))
	return fmt.Errorf("check failed with %d issue(s)", len(issues))
}
