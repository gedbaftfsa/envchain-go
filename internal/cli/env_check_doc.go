package cli

// env-check command documentation.
const envCheckUsage = `Usage: envchain env-check <project>

Compare the keys stored in <project> against the variables currently
present in the running process environment.

Exits 0 when every stored key is present in the environment with an
identical value.  Exits non-zero and prints a summary when any key is
missing or has a different value.

This is useful in CI pipelines to assert that the execution environment
has been seeded correctly before running a job.
`
