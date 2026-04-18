// Package cli provides the command-line interface for envchain-go.
//
// # Version
//
// The version command prints build metadata embedded at compile time.
// To embed version info, build with:
//
//	go build -ldflags "-X 'github.com/youruser/envchain-go/internal/cli.Version=1.0.0' \
//	  -X 'github.com/youruser/envchain-go/internal/cli.Commit=$(git rev-parse --short HEAD)' \
//	  -X 'github.com/youruser/envchain-go/internal/cli.BuildDate=$(date -u +%Y-%m-%d)'"
//
// Usage:
//
//	envchain version
package cli
