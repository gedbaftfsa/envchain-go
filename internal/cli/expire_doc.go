package cli

// Expire and check-expiry commands.
//
// expire <project> <days>
//
//	Marks a project's environment variable set with an expiry date
//	set <days> from today. The expiry is stored as metadata alongside
//	the encrypted variables.
//
// check-expiry <project>
//
//	Reads the stored expiry date for a project and reports whether it
//	has passed. Exits with a non-zero status code if the project is
//	expired, making it safe to use in CI pre-flight checks.
