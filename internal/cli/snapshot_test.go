package cli

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/env"
)

func seedSnapshot(t *testing.T) (*testStore, string) {
	t.Helper()
	st := newTempStore(t)
	set := mustEnvSet(t, "FOO=bar", "BAZ=qux")
	if err := st.Save("myproject", "pass", set); err != nil {
		t.Fatal(err)
	}
	return st, "pass"
}

func TestCmdSnapshotCreates(t *testing.T) {
	st, pass := seedSnapshot(t)

	if err := CmdSnapshot(st.Store, "myproject", pass); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	projects, err := st.ListProjects()
	if err != nil {
		t.Fatal(err)
	}

	var snaps []string
	for _, p := range projects {
		if strings.HasPrefix(p, "myproject@") {
			snaps = append(snaps, p)
		}
	}
	if len(snaps) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snaps))
	}
}

func TestCmdSnapshotWrongPassphrase(t *testing.T) {
	st, _ := seedSnapshot(t)
	err := CmdSnapshot(st.Store, "myproject", "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestCmdRestoreSnapshot(t *testing.T) {
	st, pass := seedSnapshot(t)

	if err := CmdSnapshot(st.Store, "myproject", pass); err != nil {
		t.Fatal(err)
	}

	projects, _ := st.ListProjects()
	var snapName string
	for _, p := range projects {
		if strings.HasPrefix(p, "myproject@") {
			snapName = p
		}
	}

	// mutate original
	newSet := mustEnvSet(t, "FOO=changed")
	_ = st.Save("myproject", pass, newSet)

	if err := CmdRestoreSnapshot(st.Store, snapName, pass); err != nil {
		t.Fatalf("restore error: %v", err)
	}

	loaded, err := st.Load("myproject", pass)
	if err != nil {
		t.Fatal(err)
	}
	val, _ := loaded.Get("FOO")
	if val != "bar" {
		t.Fatalf("expected bar, got %s", val)
	}
}

func TestCmdRestoreSnapshotBadName(t *testing.T) {
	st, pass := seedSnapshot(t)
	_ = st.Save("notasnapshot", pass, mustEnvSet(t, "X=1"))
	err := CmdRestoreSnapshot(st.Store, "notasnapshot", pass)
	if err == nil {
		t.Fatal("expected error for non-snapshot name")
	}
}

func TestCmdListSnapshots(t *testing.T) {
	st, pass := seedSnapshot(t)
	_ = CmdSnapshot(st.Store, "myproject", pass)
	_ = CmdSnapshot(st.Store, "myproject", pass)

	if err := CmdListSnapshots(st.Store, "myproject"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var _ = env.NewSet
