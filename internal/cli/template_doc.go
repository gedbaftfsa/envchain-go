package cli

// Template command documentation.
//
// Usage:
//
//	envchain template <project> <template-string>
//	envchain template-file <project> <path>
//
// Renders a template by substituting {KEY} or {{{KEY}}} placeholders
// with values from the named project's environment set.
//
// Example:
//
//	envchain template myapp "host={DB_HOST} port={DB_PORT}"
//	envchain template-file myapp ./config.tmpl
const templateDoc = `render a template with project env vars`
