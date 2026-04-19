package cli

import "github.com/nicholasgasior/envchain-go/internal/store"

// fmtOverride allows tests to intercept fmt dispatch without a real TTY.
var fmtOverride func(*store.Store, string, string) error

// RegisterFmtOverride sets a test hook for the fmt command dispatcher.
// Pass nil to clear the hook.
func RegisterFmtOverride(fn func(*store.Store, string, string) error) {
	fmtOverride = fn
}
