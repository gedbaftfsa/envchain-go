// Package env defines the core data types for envchain-go environment
// variable sets.
//
// A [Set] is a named collection of key-value pairs that can be persisted
// via the store package, merged together, or applied to the current
// process environment.
//
// Typical usage:
//
//	s := env.NewSet("myproject")
//	_ = s.Put("DATABASE_URL", "postgres://localhost/dev")
//	_ = env.ApplyToProcess(s)
package env
