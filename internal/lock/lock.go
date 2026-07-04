// Package lock reads and writes grimoire.lock files.
// grimoire.lock pins exact resolved versions for reproducible installs.
// It should be committed to version control.
package lock

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Entry is one resolved skill in the lockfile.
type Entry struct {
	Name     string `toml:"name"`
	Version  string `toml:"version"`
	Source   string `toml:"source"`   // dep key prefix, e.g. "acmecorp/practices" or "github.com/jeffreytse/grimoire-core"
	Resolved string `toml:"resolved"` // full git URL
	Commit   string `toml:"commit"`   // git commit SHA
	Checksum string `toml:"checksum"` // sha256:<hex>
}

// LockFile is the parsed grimoire.lock file.
type LockFile struct {
	Skills []Entry `toml:"skills"`
}

// ParseFile reads a grimoire.lock file.
// Returns an empty LockFile when the file is absent.
func ParseFile(path string) (LockFile, error) {
	var lf LockFile
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return lf, nil
	}
	if err != nil {
		return lf, err
	}
	if err := toml.Unmarshal(data, &lf); err != nil {
		return lf, err
	}
	return lf, nil
}

// WriteFile serializes lf to a grimoire.lock file, creating parent dirs as needed.
func WriteFile(path string, lf LockFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(lf)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Find returns the Entry for the given skill name, or nil if not present.
func (lf *LockFile) Find(name string) *Entry {
	for i := range lf.Skills {
		if lf.Skills[i].Name == name {
			return &lf.Skills[i]
		}
	}
	return nil
}

// Upsert adds or replaces the entry for e.Name.
func (lf *LockFile) Upsert(e *Entry) {
	for i := range lf.Skills {
		if lf.Skills[i].Name == e.Name {
			lf.Skills[i] = *e
			return
		}
	}
	lf.Skills = append(lf.Skills, *e)
}

// Remove deletes the entry for the given skill name. No-op if not present.
func (lf *LockFile) Remove(name string) {
	for i, e := range lf.Skills {
		if e.Name == name {
			lf.Skills = append(lf.Skills[:i], lf.Skills[i+1:]...)
			return
		}
	}
}
