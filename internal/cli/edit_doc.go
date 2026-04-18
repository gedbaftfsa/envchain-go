/*
Package cli — edit command

The edit command opens the environment variable set for a project in the
user's preferred text editor (determined by the $EDITOR environment variable,
defaulting to vi).

Usage:

	envchain edit <project>

The current variables are written to a temporary file in KEY=VALUE format.
After the editor exits the file is parsed and the project store is updated
with the new contents. Lines beginning with '#' and blank lines are ignored.
*/
package cli
