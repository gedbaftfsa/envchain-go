package store

// Dir returns the filesystem path of the store directory.
func (s *Store) Dir() string {
	return s.dir
}
