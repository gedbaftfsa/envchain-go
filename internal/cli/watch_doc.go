package cli

// watch_doc.go documents the watch sub-command.
//
// Usage:
//
//	envchain watch <project> [--] <cmd> [args...]
//
// watch loads the named project's environment variables and runs the given
// command with them injected into its environment. The process is restarted
// automatically when a SIGHUP signal is received, allowing live reloads after
// updating the project's stored variables. Send SIGINT or SIGTERM to stop.
