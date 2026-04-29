package cli

// Drain command documentation.
//
// Usage:
//
//	envchain drain <project>
//
// Removes all environment variable keys from the given project without
// deleting the project itself. The project remains in the store with an
// empty key set, ready to be repopulated.
//
// This is useful when you want to reset a project's variables while
// retaining its entry in the store (e.g., for auditing or re-initialisation).
const drainDoc = "drain <project>  — remove all keys from a project, keeping the project intact"
