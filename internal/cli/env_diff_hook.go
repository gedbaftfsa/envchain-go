package cli

import "fmt"

func init() {
	registerCommand("env-diff", envDiffHandler)
}

// envDiffHandler is the Main dispatch handler for the "env-diff" sub-command.
// It reads the passphrase, then delegates to CmdEnvDiff.
func envDiffHandler(args []string, st interface{ Loader }, passReader func() (string, error), w interface{ Writer }) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain env-diff <project>")
	}
	project := args[0]

	pass, err := passReader()
	if err != nil {
		return fmt.Errorf("reading passphrase: %w", err)
	}

	s, ok := st.(storeInterface)
	if !ok {
		return fmt.Errorf("internal: store does not implement required interface")
	}

	return CmdEnvDiff(s, project, pass, w.(writerInterface))
}
