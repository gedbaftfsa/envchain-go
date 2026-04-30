package cli

func init() {
	registerSummaryCommand()
}

// registerSummaryCommand wires the "summary" subcommand into the global
// dispatch table used by Main. It is called automatically via init.
func registerSummaryCommand() {
	registerCommand("summary", func(args []string, st interface{ ListProjects() ([]string, error) }, passphrase string) commandResult {
		// Actual dispatch is handled in Main via the switch statement;
		// this hook exists so the command appears in help output and
		// shell-completion lists without modifying main.go directly.
		return commandResult{handled: false}
	})
}

// commandResult is a lightweight sentinel used by hook registrations.
type commandResult struct {
	handled bool
}

// commandHook is the signature for lazily-registered command hooks.
type commandHook func(args []string, st interface{ ListProjects() ([]string, error) }, passphrase string) commandResult

var commandRegistry = map[string]commandHook{}

func registerCommand(name string, h commandHook) {
	commandRegistry[name] = h
}
