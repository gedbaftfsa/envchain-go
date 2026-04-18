// Package cli — export/import subcommands.
//
// CmdExport writes all variables for a project to an io.Writer in
// POSIX shell "export KEY=VALUE" format, suitable for eval or sourcing:
//
//	$ envchain export myapp > myapp.env
//	$ eval "$(envchain export myapp)"
//
// CmdImport reads KEY=VALUE (or "export KEY=VALUE") lines from a byte
// slice and merges them into the named project, creating it if needed.
// Lines beginning with '#' and blank lines are ignored.
//
//	$ envchain import myapp < myapp.env
package cli
