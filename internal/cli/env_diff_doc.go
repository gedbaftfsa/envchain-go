package cli

// env-diff command documentation.
//
// Usage:
//
//	envchain env-diff <project>
//
// Compares the variables stored in <project> against the current process
// environment and prints a three-section diff:
//
//	-  key  — present in the project but missing from the environment
//	~  key  — present in both but the live value differs from the stored one
//	+  key  — present in the environment but not tracked by the project
//
// Exits 0 when there are no differences.
var _ = "env-diff doc anchor"
