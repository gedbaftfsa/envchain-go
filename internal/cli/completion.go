package cli

import (
	"fmt"
	"io"
	"strings"
)

// supportedShells lists the shells for which completion scripts are available.
var supportedShells = []string{"bash", "zsh", "fish"}

// CmdCompletion prints shell completion scripts for the given shell.
// Supported shells are: bash, zsh, and fish.
func CmdCompletion(shell string, w io.Writer) error {
	switch strings.ToLower(shell) {
	case "bash":
		fmt.Fprint(w, bashCompletion)
	case "zsh":
		fmt.Fprint(w, zshCompletion)
	case "fish":
		fmt.Fprint(w, fishCompletion)
	default:
		return fmt.Errorf("unsupported shell %q: choose %s", shell, strings.Join(supportedShells, ", "))
	}
	return nil
}

const bashCompletion = `# envchain-go bash completion
_envchain_go() {
    local cur prev words
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    local commands="set unset list delete export import completion run"
    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
    fi
}
complete -F _envchain_go envchain-go
`

const zshCompletion = `# envchain-go zsh completion
#compdef envchain-go
_envchain_go() {
    local -a commands
    commands=(
        'set:Set an environment variable in a project'
        'unset:Remove an environment variable from a project'
        'list:List variables in a project'
        'delete:Delete an entire project'
        'export:Export project variables as shell exports'
        'import:Import variables from shell export format'
        'completion:Print shell completion script'
        'run:Run a command with project environment injected'
    )
    _describe 'command' commands
}
_envchain_go
`

const fishCompletion = `# envchain-go fish completion
set -l commands set unset list delete export import completion run
complete -c envchain-go -f -n "not __fish_seen_subcommand_from $commands" -a "$commands"
`
