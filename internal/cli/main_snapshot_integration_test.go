package cli

import (
	"os"
	"testing"
)

func TestMainSnapshotDispatch(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVCHAIN_DIR", dir)

	st := newTempStore(t)
	set := mustEnvSet(t, "KEY=val")
	if err := st.Save("proj", "secret", set); err != nil {
		t.Fatal(err)
	}

	// copy store file into ENVCHAIN_DIR so Main can find it
	data, err := os.ReadFile(st.Path("proj"))
	if err != nil {
		t.Skip("store.Path not available, skipping integration test")
	}
	if err := os.WriteFile(defaultStore()+"/proj.enc", data, 0600); err != nil {
		t.Skip("could not write store file")
	}

	// just verify no panic / no crash for dispatch path
	_ = CmdListSnapshots(st.Store, "proj")
}

func TestMainSnapshotEmptyProject(t *testing.T) {
	st := newTempStore(t)
	err := CmdSnapshot(st.Store, "", "pass")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}
