package env

import (
	"testing"
)

func TestMergeOverwrite(t *testing.T) {
	dst := NewSet()
	_ = dst.Put("A", "old")

	src := NewSet()
	_ = src.Put("A", "new")
	_ = src.Put("B", "val")

	added, skipped := Merge(dst, src, MergeOverwrite)
	if added != 2 || skipped != 0 {
		t.Fatalf("expected added=2 skipped=0, got %d %d", added, skipped)
	}
	v, _ := dst.Get("A")
	if v != "new" {
		t.Fatalf("expected A=new, got %s", v)
	}
}

func TestMergeSkip(t *testing.T) {
	dst := NewSet()
	_ = dst.Put("A", "old")

	src := NewSet()
	_ = src.Put("A", "new")
	_ = src.Put("B", "val")

	added, skipped := Merge(dst, src, MergeSkip)
	if added != 1 || skipped != 1 {
		t.Fatalf("expected added=1 skipped=1, got %d %d", added, skipped)
	}
	v, _ := dst.Get("A")
	if v != "old" {
		t.Fatalf("expected A=old, got %s", v)
	}
}

func TestApplyToProcess(t *testing.T) {
	s := NewSet()
	_ = s.Put("FOO", "bar")
	result := ApplyToProcess([]string{"EXISTING=1"}, s)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFromProcess(t *testing.T) {
	s := FromProcess([]string{"KEY=val", "NOEQUALS", "X=y=z"})
	v, ok := s.Get("KEY")
	if !ok || v != "val" {
		t.Fatalf("expected KEY=val")
	}
	v2, ok2 := s.Get("X")
	if !ok2 || v2 != "y=z" {
		t.Fatalf("expected X=y=z")
	}
}
