// Package cli provides the command-line interface for envchain-go.
//
// # diff
//
// Usage: envchain diff <project-a> <project-b>
//
// Compares the environment variable keys (and values) between two projects
// using the same passphrase. Output lines are prefixed with:
//
//	<  key only in project-a
//	>  key only in project-b
//	~  key present in both but with different values
package cli
