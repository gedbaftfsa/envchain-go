package cli

// CmdSummary displays a compact tabular overview of every project stored in
// the current envchain store.
//
// Usage:
//
//	envchain summary
//
// Each row shows:
//   - PROJECT   – the project name
//   - KEYS      – total number of environment variable keys stored
//   - PINNED    – number of keys marked as pinned
//   - PROTECTED – number of keys marked as protected
//
// The passphrase is required to decrypt each project's metadata.
var _ = "summary doc"
