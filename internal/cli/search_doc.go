package cli

// search subcommand
//
// Usage:
//   envchain search <query>
//
// Searches for environment variable keys matching <query> (case-insensitive
// substring match) across all projects stored in the current store.
//
// The passphrase is prompted once and reused for each project decryption.
// Results are printed as "<project>\t<key>" pairs, sorted by project then key.
