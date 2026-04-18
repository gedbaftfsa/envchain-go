package cli

// lint subcommand documentation.
//
// Usage:
//
//	envchain lint <project>
//
// Checks the named project's variable set for common issues:
//   - Empty or whitespace-only values
//   - Keys that shadow well-known system variables (PATH, HOME, USER, …)
//
// Exits with a non-zero status when any issues are found.
