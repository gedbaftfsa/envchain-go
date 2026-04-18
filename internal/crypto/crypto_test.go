package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("mysecret")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := DeriveKey("mysecret")
	if !bytes.Equal(key, key2) {
		t.Fatal("expected deterministic key derivation")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "test-passphrase"
	key := DeriveKey(passphrase)
	plaintext := []byte("MY_VAR=supersecret")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := DeriveKey("correct-passphrase")
	wrongKey := DeriveKey("wrong-passphrase")

	ciphertext, err := Encrypt(key, []byte("secret=value"))
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	_, err = Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestDecryptTooShort(t *testing.T) {
	key := DeriveKey("passphrase")
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
