// Package cli provides the command-line interface for envchain-go.
//
// # rename
//
// The rename command moves a project namespace to a new name while
// preserving all stored environment variables. The passphrase is used
// to decrypt the source project and re-encrypt it under the new name.
//
// Usage:
//
//	envchain rename <old-project> <new-project>
//
// The command fails if the destination project already exists or if the
// source project cannot be found.
package cli
