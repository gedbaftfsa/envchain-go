package cli

// recap command documentation
const recapDoc = `
Usage: envchain recap <project>

Print a summary of all environment variable keys stored in <project>
without revealing their values. Each line shows:

  KEY_NAME   set|empty   (~N chars)

This is useful for a quick audit of what is configured without
exposing sensitive data to the terminal or logs.

Examples:
  envchain recap myapp
  envchain recap production
`
