// Package cli provides high-level command handlers for the envchain-go CLI.
//
// Each exported function corresponds to a top-level subcommand (e.g. run, set,
// get) and operates on a [store.Store] to load or persist encrypted environment
// variable sets. The functions are intentionally decoupled from any specific
// CLI framework so they can be wired to cobra, flag, or tested directly.
package cli
