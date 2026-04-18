package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/envchain-go/internal/store"
)

// Main is the top-level entry point for the envchain CLI.
// It parses os.Args and dispatches to the appropriate command.
func Main() int {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return 1
	}

	st, err := defaultStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "set":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "usage: envchain set <project> KEY=VALUE ...")
			return 1
		}
		pass, err := ReadPassphrase("Passphrase: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if err := CmdSet(st, pass, rest[0], rest[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	case "unset":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "usage: envchain unset <project> KEY ...")
			return 1
		}
		pass, err := ReadPassphrase("Passphrase: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if err := CmdUnset(st, pass, rest[0], rest[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	case "list":
		if len(rest) < 1 {
			fmt.Fprintln(os.Stderr, "usage: envchain list <project>")
			return 1
		}
		pass, err := ReadPassphrase("Passphrase: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if err := CmdList(st, pass, rest[0], true, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	case "delete":
		if len(rest) < 1 {
			fmt.Fprintln(os.Stderr, "usage: envchain delete <project>")
			return 1
		}
		if err := CmdDelete(st, rest[0]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	case "run":
		return Run(st, rest)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		return 1
	}
	return 0
}

func defaultStore() (*store.Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return store.New(filepath.Join(home, ".envchain"))
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `envchain — per-project encrypted environment variables

Usage:
  envchain set    <project> KEY=VALUE ...
  envchain unset  <project> KEY ...
  envchain list   <project>
  envchain delete <project>
  envchain run    <project> [--] <command> [args...]`)
}
