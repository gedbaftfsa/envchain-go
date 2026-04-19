package cli

// verify command documentation.
//
// Usage:
//
//	envchain verify <project>
//
// Verifies that the encrypted store for <project> can be decrypted with the
// supplied passphrase. Exits with a non-zero status if decryption fails.
// No environment variable values are printed.
const verifyDoc = `verify <project>  — check passphrase against stored data without revealing values`
