// Package cli implements the command-line interface for envchain-go.
package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/envchain-go/internal/env"
	"github.com/envchain-go/internal/store"
)

// RunOptions holds parameters for the run subcommand.
type RunOptions struct {
	Project    string
	Passphrase string
	Args       []string
	Overwrite  bool
}

// Run loads the env set for a project and executes the given command with
// those variables injected into the process environment.
func Run(st *store.Store, opts RunOptions) error {
	if len(opts.Args) == 0 {
		return fmt.Errorf("no command specified")
	}

	set, err := st.Load(opts.Project, opts.Passphrase)
	if err != nil {
		return fmt.Errorf("load project %q: %w", opts.Project, err)
	}

	base := env.FromProcess()
	merged := env.Merge(base, set, opts.Overwrite)

	cmdPath, err := exec.LookPath(opts.Args[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", opts.Args[0])
	}

	envSlice := env.ApplyToProcess(merged)

	return syscall.Exec(cmdPath, opts.Args, envSlice)
}

// RunFallback is like Run but uses exec.Cmd instead of syscall.Exec,
// which is useful on platforms where syscall.Exec is unavailable or in tests.
func RunFallback(st *store.Store, opts RunOptions, stdout, stderr *os.File) error {
	if len(opts.Args) == 0 {
		return fmt.Errorf("no command specified")
	}

	set, err := st.Load(opts.Project, opts.Passphrase)
	if err != nil {
		return fmt.Errorf("load project %q: %w", opts.Project, err)
	}

	base := env.FromProcess()
	merged := env.Merge(base, set, opts.Overwrite)

	cmd := exec.Command(opts.Args[0], opts.Args[1:]...)
	cmd.Env = env.ApplyToProcess(merged)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
