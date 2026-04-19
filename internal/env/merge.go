package env

// MergeMode controls how conflicts are resolved during a merge.
type MergeMode int

const (
	// MergeSkip leaves existing keys in dst untouched.
	MergeSkip MergeMode = iota
	// MergeOverwrite replaces existing keys in dst with values from src.
	MergeOverwrite
)

// Merge copies all key/value pairs from src into dst according to mode.
// It returns the number of keys added and the number skipped.
func Merge(dst, src *Set, mode MergeMode) (added, skipped int) {
	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		_, exists := dst.Get(k)
		if exists && mode == MergeSkip {
			skipped++
			continue
		}
		_ = dst.Put(k, v)
		added++
	}
	return
}

// ApplyToProcess returns a slice of "KEY=VALUE" strings suitable for
// passing to exec.Cmd.Env, merging s on top of base.
func ApplyToProcess(base []string, s *Set) []string {
	result := make([]string, len(base))
	copy(result, base)
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		result = append(result, k+"="+v)
	}
	return result
}

// FromProcess parses a slice of "KEY=VALUE" strings into a Set.
func FromProcess(environ []string) *Set {
	s := NewSet()
	for _, e := range environ {
		if k, v, ok := ParseEntry(e); ok {
			_ = s.Put(k, v)
		}
	}
	return s
}
