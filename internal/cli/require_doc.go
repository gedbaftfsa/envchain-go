package cli

// RequireDoc is the help text for the require command.
const RequireDoc = `require — assert that required keys exist in a project

Usage:
  envchain require <project> <KEY1> [KEY2 ...]

Description:
  Loads the named project and verifies that every listed key exists
  and has a non-empty value. Exits with a non-zero status if any key
  is absent or blank, printing each offending key to stdout.

  Useful in shell scripts or CI pipelines to guard against accidentally
  running with incomplete configuration.

Examples:
  envchain require myapp DATABASE_URL SECRET_KEY
  envchain require staging AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY
`
