package tempfile

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	once sync.Once
	mu   sync.Mutex
	r    *rand.Rand
	b    []byte
)

const (
	maxAttempts = 1000
	reseedAfter = 10
)

// New creates a new temporary file in the specified directory the same way
// ioutils.TempFile does. If dir is an empty string, it'll be created in the
// OS' temp directory. The caller is responsible for removing the file. Drop-in
// replacement for ioutils.TempFile
func New(dir, prefix string) (f *os.File, err error) {
	return NewSuffixAndMode(dir, prefix, "", 0600)
}

// NewSuffix does the same as New() but allows a suffix to be added to the name.
func NewSuffix(dir, prefix, suffix string) (f *os.File, err error) {
	return NewSuffixAndMode(dir, prefix, suffix, 0600)
}

// NewMode does the same as New() but allows the caller to provide a file mode
// to give the file permissions different from the default 0600.
func NewMode(dir, prefix string, perm os.FileMode) (f *os.File, err error) {
	return NewSuffixAndMode(dir, prefix, "", perm)
}

// NewSuffixAndMode does the same as New(), but allows the caller to provide
// a suffix as well as a file mode to be used for the temporary file.
func NewSuffixAndMode(dir, prefix, suffix string, perm os.FileMode) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	var conflicts int
	for i := 0; i < maxAttempts; i++ {
		name := filepath.Join(dir, prefix+nextSuffix()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, perm)
		if os.IsExist(err) {
			if conflicts++; conflicts > reseedAfter {
				reseed()
			}
			continue
		}
		return
	}
	return nil, fmt.Errorf("unable to create tempfile after %d attempts", maxAttempts)
}

func nextSuffix() string {
	once.Do(func() {
		b = make([]byte, 8)
		reseed()
	})
	mu.Lock()
	defer mu.Unlock()
	r.Read(b)
	return fmt.Sprintf(".%x", b)
}

func reseed() {
	mu.Lock()
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu.Unlock()
}
