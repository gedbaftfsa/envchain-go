// Package crypto provides AES-GCM encryption and decryption utilities
// for envchain-go. It is used to encrypt environment variable sets at
// rest using a user-supplied passphrase.
//
// Key derivation is performed via SHA-256 hashing of the passphrase,
// producing a 32-byte key suitable for AES-256.
//
// Usage:
//
//	key := crypto.DeriveKey("my-passphrase")
//	ciphertext, err := crypto.Encrypt(key, []byte("API_KEY=abc123"))
//	plaintext, err := crypto.Decrypt(key, ciphertext)
package crypto
