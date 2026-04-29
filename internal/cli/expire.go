package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdExpire marks a project's variables as expiring after a given number of days.
// Usage: envchain expire <project> <days>
func CmdExpire(st *store.Store, passphrase, project string, args []string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}
	if len(args) < 1 {
		return fmt.Errorf("usage: expire <project> <days>")
	}

	days, err := strconv.Atoi(args[0])
	if err != nil || days <= 0 {
		return fmt.Errorf("days must be a positive integer")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	expiry := time.Now().UTC().AddDate(0, 0, days)
	set.SetMeta("expires_at", expiry.Format(time.RFC3339))

	if err := st.Save(project, passphrase, set); err != nil {
		return err
	}

	fmt.Fprintf(w, "project %q will expire on %s\n", project, expiry.Format("2006-01-02"))
	return nil
}

// CmdCheckExpiry checks whether a project's variables have passed their expiry date.
// Usage: envchain check-expiry <project>
func CmdCheckExpiry(st *store.Store, passphrase, project string, w io.Writer) error {
	if project == "" {
		return fmt.Errorf("project name required")
	}

	set, err := st.Load(project, passphrase)
	if err != nil {
		return err
	}

	raw := set.Meta("expires_at")
	if raw == "" {
		fmt.Fprintf(w, "project %q has no expiry set\n", project)
		return nil
	}

	expiry, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fmt.Errorf("invalid expiry value stored: %w", err)
	}

	if time.Now().UTC().After(expiry) {
		fmt.Fprintf(w, "EXPIRED: project %q expired on %s\n", project, expiry.Format("2006-01-02"))
		os.Exit(1)
	}

	remaining := time.Until(expiry).Hours() / 24
	fmt.Fprintf(w, "project %q expires on %s (%.0f day(s) remaining)\n",
		project, expiry.Format("2006-01-02"), remaining)
	return nil
}
