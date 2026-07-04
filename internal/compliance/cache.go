package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CheckCache stores per-file checksums and diagnostics for incremental re-checking.
// Stored at .grimoire/cache.json.
type CheckCache struct {
	Version          string                       `json:"version"`
	ConfigHash       string                       `json:"config_hash"`
	RubricHash       string                       `json:"rubric_hash,omitempty"`
	PracticeTotals   map[string]int               `json:"practice_totals"`
	PracticeCriteria map[string][]CriterionDetail `json:"practice_criteria,omitempty"`
	Files            map[string]FileCacheEntry    `json:"files"`
}

// FileCacheEntry is the cached state for one source file.
type FileCacheEntry struct {
	Hash        string       `json:"hash"`
	Mtime       int64        `json:"mtime"` // UnixNano — fast-path: skip SHA256 when mtime+size match
	Size        int64        `json:"size"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

const DefaultCachePath = ".grimoire/cache.json"

// LoadCache reads the cache from .grimoire/cache.json.
// Returns an initialized empty CheckCache if the file doesn't exist or is corrupt.
func LoadCache(projectDir string) (*CheckCache, error) {
	path := filepath.Join(projectDir, DefaultCachePath)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &CheckCache{Files: make(map[string]FileCacheEntry)}, nil
	}
	if err != nil {
		return &CheckCache{Files: make(map[string]FileCacheEntry)}, err
	}
	var c CheckCache
	if err := json.Unmarshal(data, &c); err != nil {
		return &CheckCache{Files: make(map[string]FileCacheEntry)}, nil
	}
	if c.Files == nil {
		c.Files = make(map[string]FileCacheEntry)
	}
	return &c, nil
}

// SaveCache writes the cache to .grimoire/cache.json.
func SaveCache(c *CheckCache, projectDir string) error {
	dir := filepath.Join(projectDir, ".grimoire")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating .grimoire dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling cache: %w", err)
	}
	return os.WriteFile(filepath.Join(projectDir, DefaultCachePath), data, 0o644)
}

// FileHash returns the SHA256 hex digest of the file at path.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ConfigHash returns the SHA256 hex digest of grimoire.toml.
// Returns a stable empty-file hash if grimoire.toml doesn't exist.
func ConfigHash(projectDir string) (string, error) {
	h := sha256.New()
	data, err := os.ReadFile(filepath.Join(projectDir, "grimoire.toml"))
	if err == nil {
		_, _ = h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
