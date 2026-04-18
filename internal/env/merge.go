package env

import "os"

// MergeOption controls how Merge handles conflicts.
type MergeOption int

const (
	// MergeSkip keeps existing values on conflict.
	MergeSkip MergeOption = iota
	// MergeOverwrite replaces existing values on conflict.
	MergeOverwrite
)

// Merge copies variables from src into dst.
// Behaviour on key conflicts is controlled by opt.
func Merge(dst, src *Set, opt MergeOption) {
	for k, v := range src.Vars {
		if _, exists := dst.Vars[k]; exists && opt == MergeSkip {
			continue
		}
		dst.Vars[k] = v
	}
}

// ApplyToProcess sets the variables in s as environment variables
// on the current process. Existing OS variables are not cleared.
func ApplyToProcess(s *Set) error {
	for k, v := range s.Vars {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// FromProcess creates a Set populated from the current process environment
// for the given keys. Missing keys are silently skipped.
func FromProcess(name string, keys []string) *Set {
	s := NewSet(name)
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			s.Vars[k] = v
		}
	}
	return s
}
