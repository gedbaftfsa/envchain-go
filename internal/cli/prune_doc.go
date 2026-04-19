// Package cli — prune command
//
// Usage:
//
//	envchain prune <passphrase>
//
// CmdPrune scans all projects in the store and deletes any that contain
// no environment variable entries. This is useful for cleaning up
// projects that were initialised but never populated.
//
// Example:
//
//	$ envchain prune
//	pruned empty-project
//	$ envchain prune
//	nothing to prune
package cli
