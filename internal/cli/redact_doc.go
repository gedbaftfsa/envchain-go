package cli

// Redact command documentation.
//
// Usage:
//
//	envchain redact <project> <passphrase> [text]
//	envchain redact-file <project> <passphrase> <file|->
//
// Description:
//	Scans the supplied text (or file) and replaces every known secret value
//	stored under <project> with the placeholder ***REDACTED***.
//
//	This is useful when you need to share log files or command output that
//	might inadvertently contain credentials.
//
// Examples:
//
//	# Redact a literal string
//	envchain redact myapp s3cr3t "token is abc123"
//
//	# Pipe a log file through redact
//	cat app.log | envchain redact-file myapp s3cr3t -
