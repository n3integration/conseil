package conseil

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// StageTestDir sets up a testing directory structure for our tests
func StageTestDir(t *testing.T) (string, func()) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	wd := filepath.Join(cwd, "__staging")
	if err := os.MkdirAll(wd, 0755); err != nil {
		t.Fatalf("err: %s", err)
	}
	return wd, func() {
		CleanDir(t, wd)
	}
}

// CleanDir recursively removes all generated content in dir
func CleanDir(t *testing.T, dir string) {
	d, err := os.Open(dir)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

// FileExists asserts that path exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileCount gets a count of files matching pattern
func FileCount(base string, pattern string) int {
	count := 0
	re := regexp.MustCompile(pattern)
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if re.MatchString(path) {
			count++
		}
		return err
	})
	return count
}
