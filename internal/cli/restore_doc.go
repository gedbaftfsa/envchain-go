package cli

// restore_doc.go documents the restore command.
//
// Usage:
//
//	envchain restore <project> <file>
//
// Description:
//
//	Restores a project's environment variable set from a backup file
//	previously created with the archive command. The file must be in
//	the envchain export format (KEY=value lines). You will be prompted
//	for a passphrase to encrypt the restored data at rest.
//
// Example:
//
//	envchain restore myproject myproject.env
