package cli

// share.go — CmdShare / CmdReceive
//
// CmdShare serialises a project's environment variables to a human-readable,
// line-oriented text format prefixed with a project-name header:
//
//	# envchain-share project=myapp
//	DATABASE_URL='postgres://...'
//	SECRET_KEY='s3cr3t'
//
// The output can be redirected to a file or piped directly to another
// invocation of envchain via CmdReceive:
//
//	envchain share myapp | ssh remote envchain receive
//
// CmdReceive reads that format from any io.Reader, optionally overriding the
// embedded project name, and persists the variables into the local store
// protected by the supplied passphrase.
