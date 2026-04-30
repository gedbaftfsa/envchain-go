package cli

// Namespace command documentation.
//
// Usage:
//
//	envchain namespace <project>
//		List all namespace prefixes (e.g. DB, AWS, APP) derived from
//		keys stored in the given project.
//
//	envchain namespace <project> <namespace>
//		List all keys that belong to the given namespace prefix.
//
// A key is considered part of a namespace when it contains an underscore.
// For example, DB_HOST and DB_PASS both belong to the "DB" namespace.
