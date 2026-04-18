package env_test

import (
	"testing"

	"github.com/yourorg/envchain-go/internal/env"
)

func TestNewSet(t *testing.T) {
	s := env.NewSet("myproject")
	if s.Name != "myproject" {
		t.Fatalf("expected name myproject, got %s", s.Name)
	}
	if len(s.Vars) != 0 {
		t.Fatal("expected empty vars")
	}
}

func TestPutAndKeys(t *testing.T) {
	s := env.NewSet("test")
	_ = s.Put("ZEBRA", "z")
	_ = s.Put("ALPHA", "a")
	_ = s.Put("MIDDLE", "m")

	keys := s.Keys()
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("expected %s at index %d, got %s", k, i, keys[i])
		}
	}
}

func TestPutEmptyKey(t *testing.T) {
	s := env.NewSet("test")
	if err := s.Put("", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestDelete(t *testing.T) {
	s := env.NewSet("test")
	_ = s.Put("FOO", "bar")
	if !s.Delete("FOO") {
		t.Fatal("expected true when deleting existing key")
	}
	if s.Delete("FOO") {
		t.Fatal("expected false when deleting missing key")
	}
}

func TestParseEntry(t *testing.T) {
	tests := []struct {
		input   string
		key     string
		value   string
		wantErr bool
	}{
		{"KEY=VALUE", "KEY", "VALUE", false},
		{"KEY=val=with=equals", "KEY", "val=with=equals", false},
		{"KEY=", "KEY", "", false},
		{"NOEQUALS", "", "", true},
		{"=VALUE", "", "", true},
	}
	for _, tt := range tests {
		k, v, err := env.ParseEntry(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseEntry(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (k != tt.key || v != tt.value) {
			t.Errorf("ParseEntry(%q) = %q, %q; want %q, %q", tt.input, k, v, tt.key, tt.value)
		}
	}
}

func TestToEnvSlice(t *testing.T) {
	s := env.NewSet("test")
	_ = s.Put("A", "1")
	_ = s.Put("B", "2")
	slice := s.ToEnvSlice()
	if len(slice) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(slice))
	}
}
