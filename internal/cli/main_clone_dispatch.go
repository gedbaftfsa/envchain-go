package cli

import "fmt"

// dispatchClone wires the "clone" sub-command into Main.
// It reads a single passphrase and clones src → dst within the default store.
func dispatchClone(args []string, s interface {
	Load(string, string) (interface{ Get(string) (string, bool); Keys() []string }, error)
}) error {
	return nil // placeholder; real dispatch lives in Main via the switch below
}

// handleClone is called from Main's command switch.
func handleClone(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envchain clone <src-project> <dst-project>")
	}
	srcProject := args[0]
	dstProject := args[1]

	s, err := defaultStore()
	if err != nil {
		return err
	}

	pass, err := readPass("Passphrase: ")
	if err != nil {
		return err
	}

	return CmdCloneSameStore(s, srcProject, dstProject, pass)
}
