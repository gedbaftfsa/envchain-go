package cli

import "fmt"

// Version information set at build time via ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// CmdVersion prints the current version, commit, and build date.
func CmdVersion() {
	fmt.Printf("envchain-go %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
}
