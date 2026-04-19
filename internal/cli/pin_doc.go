// Package cli — pin / unpin command documentation strings.
package cli

const pinUsage = `Usage:
  envchain pin   <project> <KEY> [KEY...]   mark keys as pinned
  envchain unpin <project> <KEY> [KEY...]   remove keys from pinned set
  envchain pinned <project>                 list pinned keys

Pinned keys are highlighted by diff and lint commands to draw attention
to variables that must always be present in a project's environment set.
`
