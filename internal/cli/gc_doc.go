package cli

// GC command documentation.
//
// Usage:
//
//	envchain gc
//
// gc scans the store for snapshot entries that reference projects
// which no longer exist and removes them, reclaiming disk space.
//
// No passphrase is required because only metadata keys are inspected;
// encrypted blobs are deleted without being decrypted.
const gcDoc = `Remove orphaned snapshots for deleted projects.

Usage:
  envchain gc

Scans all stored snapshots and deletes any whose parent project no
longer exists in the store. Safe to run at any time.
`
