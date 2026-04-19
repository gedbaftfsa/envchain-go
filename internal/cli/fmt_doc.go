package cli

// CmdFmt sorts the keys of a stored project env-set alphabetically and
// re-saves it with the same passphrase.  The encrypted file is rewritten
// in-place; no data is lost.
//
// Usage:
//
//	envchain fmt <project>
//
// CmdFmtDiff performs a dry-run, printing the new key order without saving.
//
// Usage:
//
//	envchain fmt --diff <project>
