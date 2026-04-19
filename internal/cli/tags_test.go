package cli

import (
	"strings"
	"testing"
)

func seedTags(t *testing.T) *storeT {
	t.Helper()
	st, pass := newTempStore(t)
	mustSet(t, st, pass, "alpha", "AWS_KEY=val1", "DB_URL=val2")
	mustSet(t, st, pass, "beta", "AWS_KEY=val3", "GH_TOKEN=secret")
	mustSet(t, st, pass, "gamma", "DB_URL=val4", "REDIS_URL=val5")
	return &storeT{st: st, pass: pass}
}

func mustSet(t *testing.T, st interface{ Save(string, interface{}) error }, pass, proj string, entries ...string) {
	t.Helper()
	es := newEnvSet(t, entries...)
	if err := mustEnvSet(t, st, pass, proj, es); err != nil {
		t.Fatalf("mustSet: %v", err)
	}
}

func TestCmdTagsAll(t *testing.T) {
	fix := seedTags(t)
	var buf strings.Builder
	if err := cmdTags(fix.st, fix.pass, "", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.Split(strings.TrimSpace(buf.String()), "\n")
	want := []string{"AWS_KEY", "DB_URL", "GH_TOKEN", "REDIS_URL"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("[%d] got %q want %q", i, got[i], w)
		}
	}
}

func TestCmdTagsPrefix(t *testing.T) {
	fix := seedTags(t)
	var buf strings.Builder
	if err := cmdTags(fix.st, fix.pass, "AWS", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "AWS_KEY" {
		t.Errorf("got %q want AWS_KEY", got)
	}
}

func TestCmdTagsWrongPassphrase(t *testing.T) {
	fix := seedTags(t)
	var buf strings.Builder
	err := cmdTags(fix.st, "wrong", "", &buf)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

type storeT struct {
	st   interface {
		List() ([]string, error)
		Load(string, string) (interface{ Keys() []string }, error)
	}
	pass string
}
