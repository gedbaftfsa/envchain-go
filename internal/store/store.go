// Package store manages encrypted storage of environment variable sets.
package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/yourusername/envchain-go/internal/crypto"
)

// EnvSet represents a named set of environment variables.
type EnvSet struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Store handles reading and writing encrypted env sets to disk.
type Store struct {
	Dir string
}

// New creates a Store rooted at dir, creating it if necessary.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Store{Dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.Dir, name+".enc")
}

// Save encrypts and writes an EnvSet to disk using the given passphrase.
func (s *Store) Save(set *EnvSet, passphrase string) error {
	data, err := json.Marshal(set)
	if err != nil {
		return err
	}
	key, salt, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return err
	}
	ciphertext, err := crypto.Encrypt(key, data)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(encryptedFile{Salt: salt, Data: ciphertext})
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(set.Name), payload, 0600)
}

// Load decrypts and reads an EnvSet from disk using the given passphrase.
func (s *Store) Load(name, passphrase string) (*EnvSet, error) {
	raw, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var ef encryptedFile
	if err := json.Unmarshal(raw, &ef); err != nil {
		return nil, err
	}
	key, _, err := crypto.DeriveKey(passphrase, ef.Salt)
	if err != nil {
		return nil, err
	}
	plaintext, err := crypto.Decrypt(key, ef.Data)
	if err != nil {
		return nil, ErrBadPassphrase
	}
	var set EnvSet
	if err := json.Unmarshal(plaintext, &set); err != nil {
		return nil, err
	}
	return &set, nil
}

// Delete removes an EnvSet from disk.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

type encryptedFile struct {
	Salt []byte `json:"salt"`
	Data []byte `json:"data"`
}
