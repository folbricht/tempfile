package tempfile

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	// Basic form without params. Should create random file in OS TempDir
	// with 0600 perms
	f, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	dir := filepath.Dir(f.Name())
	if dir != os.TempDir() {
		t.Fatalf("expected temp file in '%s', got '%s'", os.TempDir(), dir)
	}
	stat, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Mode() != 0600 {
		t.Fatalf("expected perms 0600, got %s", stat.Mode())
	}

	// With filename prefix and in the current dir
	f, err = New(".", "prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	dir = filepath.Dir(f.Name())
	base := filepath.Base(f.Name())
	if dir != "." {
		t.Fatalf("expected temp file in current dir, got '%s'", dir)
	}
	if !strings.HasPrefix(base, "prefix") {
		t.Fatalf("expected filename with prefix, got %s", base)
	}

	// Filename with suffix
	f, err = NewSuffix("", "", "suffix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	base = filepath.Base(f.Name())
	if !strings.HasSuffix(base, "suffix") {
		t.Fatalf("expected filename with suffix, got %s", base)
	}

	// File with non-standard perms
	f, err = NewMode("", "", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	stat, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Mode() != 0755 {
		t.Fatalf("expected perms 0755, got %s", stat.Mode())
	}
}

func TestUnique(t *testing.T) {
	// Create several tempfiles and store them in a map. Fail if any is duplicated
	fmap := make(map[string]struct{})
	for i := 0; i < 20; i++ {
		f, err := New("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
		if _, ok := fmap[f.Name()]; ok {
			t.Fatalf("file %s existed already", f.Name())
		}
		fmap[f.Name()] = struct{}{}
	}
}

func TestCollision(t *testing.T) {
	// Make several files starting with a fixed seed
	r = rand.New(rand.NewSource(0))
	for i := 0; i < 20; i++ {
		f, err := New("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
	}

	// Reset the seed and make more files, those should conflict and trigger
	// a re-seed of the sequence after a few tries
	r = rand.New(rand.NewSource(0))
	for i := 0; i < 20; i++ {
		f, err := New("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
	}
}
