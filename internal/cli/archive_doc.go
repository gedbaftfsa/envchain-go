package cli

// Archive / restore-archive commands.
//
// Usage:
//
//	envchain archive <passphrase> <dest-dir>
//		Exports every project in the current store as a plain <project>.env
//		file inside <dest-dir>.  Each file contains KEY=VALUE lines.
//
//	envchain restore-archive <passphrase> <src-dir>
//		Reads every *.env file from <src-dir> and imports it into the store
//		using the file base-name (without extension) as the project name.
//
// Files are written with mode 0600; the destination directory is created
// with mode 0700 if it does not already exist.
