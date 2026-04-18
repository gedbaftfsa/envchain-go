// Package env provides types and helpers for managing environment variable sets.
package env

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidEntry is returned when an env entry is malformed.
var ErrInvalidEntry = errors.New("invalid env entry: expected KEY=VALUE")

// Set represents a named collection of environment variables.
type Set struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// NewSet creates an empty Set with the given name.
func NewSet(name string) *Set {
	return &Set{Name: name, Vars: make(map[string]string)}
}

// Put adds or updates a key-value pair.
func (s *Set) Put(key, value string) error {
	if key == "" {
		return fmt.Errorf("%w: empty key", ErrInvalidEntry)
	}
	s.Vars[key] = value
	return nil
}

// Delete removes a key from the set. Returns false if key did not exist.
func (s *Set) Delete(key string) bool {
	_, ok := s.Vars[key]
	delete(s.Vars, key)
	return ok
}

// Keys returns a sorted list of variable names.
func (s *Set) Keys() []string {
	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// ParseEntry parses a "KEY=VALUE" string into its components.
func ParseEntry(entry string) (key, value string, err error) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", "", fmt.Errorf("%w: %q", ErrInvalidEntry, entry)
	}
	return parts[0], parts[1], nil
}

// ToEnvSlice returns the set as a slice of "KEY=VALUE" strings.
func (s *Set) ToEnvSlice() []string {
	result := make([]string, 0, len(s.Vars))
	for k, v := range s.Vars {
		result = append(result, k+"="+v)
	}
	return result
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
