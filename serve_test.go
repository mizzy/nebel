package nebel

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSnapshotDetectsChanges(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.markdown")
	if err := os.WriteFile(file, []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirs := []string{dir}

	base := snapshot(dirs)
	if snapshot(dirs) != base {
		t.Fatal("snapshot changed without file changes")
	}

	if err := os.Chtimes(file, time.Now(), time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	modified := snapshot(dirs)
	if modified == base {
		t.Fatal("modification not detected")
	}

	if err := os.WriteFile(filepath.Join(dir, "b.markdown"), []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}
	added := snapshot(dirs)
	if added == modified {
		t.Fatal("addition not detected")
	}

	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	if snapshot(dirs) == added {
		t.Fatal("removal not detected")
	}
}
