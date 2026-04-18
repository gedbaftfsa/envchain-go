package env_test

import (
	"os"
	"testing"

	"github.com/yourorg/envchain-go/internal/env"
)

func TestMergeOverwrite(t *testing.T) {
	dst := env.NewSet("dst")
	_ = dst.Put("A", "old")

	src := env.NewSet("src")
	_ = src.Put("A", "new")
	_ = src.Put("B", "added")

	env.Merge(dst, src, env.MergeOverwrite)

	if dst.Vars["A"] != "new" {
		t.Errorf("expected A=new, got %s", dst.Vars["A"])
	}
	if dst.Vars["B"] != "added" {
		t.Errorf("expected B=added, got %s", dst.Vars["B"])
	}
}

func TestMergeSkip(t *testing.T) {
	dst := env.NewSet("dst")
	_ = dst.Put("A", "original")

	src := env.NewSet("src")
	_ = src.Put("A", "overridden")
	_ = src.Put("C", "new")

	env.Merge(dst, src, env.MergeSkip)

	if dst.Vars["A"] != "original" {
		t.Errorf("expected A=original, got %s", dst.Vars["A"])
	}
	if dst.Vars["C"] != "new" {
		t.Errorf("expected C=new, got %s", dst.Vars["C"])
	}
}

func TestApplyToProcess(t *testing.T) {
	s := env.NewSet("test")
	_ = s.Put("ENVCHAIN_TEST_VAR", "hello")

	if err := env.ApplyToProcess(s); err != nil {
		t.Fatalf("ApplyToProcess error: %v", err)
	}
	if got := os.Getenv("ENVCHAIN_TEST_VAR"); got != "hello" {
		t.Errorf("expected hello, got %s", got)
	}
	os.Unsetenv("ENVCHAIN_TEST_VAR")
}

func TestFromProcess(t *testing.T) {
	os.Setenv("ENVCHAIN_PICK_A", "aaa")
	os.Setenv("ENVCHAIN_PICK_B", "bbb")
	defer os.Unsetenv("ENVCHAIN_PICK_A")
	defer os.Unsetenv("ENVCHAIN_PICK_B")

	s := env.FromProcess("picked", []string{"ENVCHAIN_PICK_A", "ENVCHAIN_PICK_B", "ENVCHAIN_MISSING"})
	if s.Vars["ENVCHAIN_PICK_A"] != "aaa" {
		t.Errorf("expected aaa")
	}
	if s.Vars["ENVCHAIN_PICK_B"] != "bbb" {
		t.Errorf("expected bbb")
	}
	if _, ok := s.Vars["ENVCHAIN_MISSING"]; ok {
		t.Error("missing key should not be present")
	}
}
