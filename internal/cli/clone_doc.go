// Package cli — clone command
//
// Usage:
//
//	envchain clone <src-project> <dst-project>
//
// Clones all environment variables from <src-project> into a new project
// called <dst-project>, encrypted with the same passphrase.
//
// To re-encrypt with a different passphrase use the --rekey flag (future).
package cli
