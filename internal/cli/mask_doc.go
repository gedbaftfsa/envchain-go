package cli

// mask command documentation.
//
// Usage:
//
//	envchain mask <project>
//
// Prints all environment variable keys stored under <project> with their
// values replaced by asterisks. The number of asterisks is capped at 8
// regardless of actual value length to avoid leaking secret length.
//
// Example:
//
//	$ envchain mask myapp
//	AWS_ACCESS_KEY_ID=********
//	AWS_SECRET_ACCESS_KEY=********
//	DATABASE_URL=********
