package cli

// CmdMerge merges the environment variable set of one project into another.
//
// Usage:
//
//	envchain merge <src-project> <dst-project> [--overwrite]
//
// By default, keys that already exist in the destination project are preserved
// (skip mode). Pass --overwrite to replace conflicting keys with values from
// the source project.
//
// Both projects must be accessible with the same passphrase. The destination
// project is created automatically if it does not yet exist.
