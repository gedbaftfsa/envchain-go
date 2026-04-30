package cli

// env-copy command documentation and dispatch registration.

const envCopyUsage = `Usage: envchain env-copy [--overwrite] <src-project> <dst-project> KEY [KEY...]

Copy one or more keys from a source project into a destination project.
Both projects must share the same passphrase.

Flags:
  --overwrite   Replace existing keys in the destination project.

Examples:
  envchain env-copy myapp staging DB_HOST DB_PORT
  envchain env-copy --overwrite myapp staging API_KEY
`

func init() {
	registerCommand("env-copy", envCopyUsage)
}
