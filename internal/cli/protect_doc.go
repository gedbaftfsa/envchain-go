package cli

// protect.go implements the protect/unprotect/list-protected sub-commands.
//
// Usage:
//
//	envchain protect <project> <KEY> [KEY...]
//	  Mark one or more keys as protected in the given project.
//	  Protected keys cannot be overwritten or deleted without --force.
//
//	envchain unprotect <project> <KEY> [KEY...]
//	  Remove protection from one or more keys in the given project.
//
//	envchain list-protected <project>
//	  List all protected keys for the given project.
//
// Protected key names are stored in the encrypted envelope under the
// reserved key "__protected__" as a comma-separated sorted list, so
// no additional metadata file is needed.
