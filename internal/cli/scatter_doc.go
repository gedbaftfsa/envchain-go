package cli

// scatter writes each environment variable in a project to its own file
// inside a specified directory. The file name is the key name and the file
// content is the raw value.
//
// Usage:
//
//	envchain scatter <project> <dir>
//
// Each file is created with mode 0600. The target directory is created with
// mode 0700 if it does not already exist. This is useful for integrating with
// systems that expect secrets to be mounted as files (e.g. Docker secrets,
// Kubernetes projected volumes, or any twelve-factor app that reads config
// from the filesystem).
