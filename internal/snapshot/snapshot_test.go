package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSave_And_Latest(t *testing.T) {
	store, err := snapshot.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e := snapshot.Entry{
		Path:    "secret/app",
		Secrets: map[string]string{"KEY": "value"},
	}
	if err := store.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Latest("secret/app")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if got.Secrets["KEY"] != "value" {
		t.Errorf("expected value, got %q", got.Secrets["KEY"])
	}
}

func TestLatest_NoSnapshots_ReturnsNil(t *testing.T) {
	store, _ := snapshot.New(tempDir(t))
	got, err := store.Latest("secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestList_OrderedByTime(t *testing.T) {
	store, _ := snapshot.New(tempDir(t))
	t1 := time.Now().UTC().Add(-2 * time.Hour)
	t2 := time.Now().UTC().Add(-1 * time.Hour)
	_ = store.Save(snapshot.Entry{Path: "secret/app", Secrets: map[string]string{"A": "1"}, CreatedAt: t1})
	_ = store.Save(snapshot.Entry{Path: "secret/app", Secrets: map[string]string{"A": "2"}, CreatedAt: t2})
	list, err := store.List("secret/app")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	if list[0].Secrets["A"] != "1" || list[1].Secrets["A"] != "2" {
		t.Error("entries not in expected order")
	}
}

func TestSave_SetsTimestampIfZero(t *testing.T) {
	store, _ := snapshot.New(tempDir(t))
	e := snapshot.Entry{Path: "secret/ts", Secrets: map[string]string{}}
	_ = store.Save(e)
	got, _ := store.Latest("secret/ts")
	if got.CreatedAt.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
