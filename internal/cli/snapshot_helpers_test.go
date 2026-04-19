package cli

import (
	"testing"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

type testStore struct {
	*store.Store
	dir string
}

func (ts *testStore) Path(project string) string {
	return ts.dir + "/" + project + ".enc"
}

func newTempStoreSnap(t *testing.T) *testStore {
	t.Helper()
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	return &testStore{Store: st, dir: dir}
}
