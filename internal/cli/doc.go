// Package cli provides high-level command handlers for the envchain-go CLI.
//
// Each exported function corresponds to a top-level subcommand (e.g. run, set,
// get) and operates on a [store.Store] to load or persist encrypted environment
// variable sets. The functions are intentionally decoupled from any specific
// CLI framework so they can be wired to cobra, flag, or tested directly.
//
// # Subcommands
//
//   - Run: executes a child process with secrets injected into its environment.
//   - Set: prompts for a secret value and persists it to the store.
//   - Get: retrieves and prints a secret value from the store.
//   - List: enumerates all namespaces or keys within a namespace.
//
// # Error Handling
//
// Functions return descriptive errors that callers should present to the user.
// Secrets are never included in error messages.
package cli
